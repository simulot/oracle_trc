package queries

import (
	"bytes"
	"fmt"
	"io"
	"runtime"
	"strconv"
	"strings"

	"github.com/simulot/oracle_trc/trc"
)

// Packet hold packet content and context of packet
type Query struct {
	Packet      *trc.Packet // Query's packet
	Query       string      // Query text
	ExeOp       uint32
	CursorId    uint32
	Len         uint32
	RowToFetch  uint32
	ParamLen    uint32
	NbofDefCols uint32
	Params      []*ParameterInfo
}

// String implement the basic representation of packet: Packet's context and its content in hexadecimal
func (q Query) String() string {
	sb := strings.Builder{}
	q.Packet.WriteContext(&sb)
	writeEol(&sb)
	sb.WriteString(q.Query)
	writeEol(&sb)
	for i, p := range q.Params {
		sb.WriteString("  :")
		sb.WriteString(strconv.Itoa(i + 1))
		sb.WriteString(" = ")
		switch p.DataType {
		case CHAR:
			sb.WriteByte('\'')
			sb.WriteString(strings.Replace(p.String(), "'", "''", -1))
			sb.WriteByte('\'')
		default:
			sb.WriteString(p.String())
		}
		writeEol(&sb)
	}
	return sb.String()
}

// Parser is used to parse trc files and extract queries
type Parser struct {
	p          *trc.Parser // Trace file parser
	q          *Query      // current query
	rowsToRead int         // rows to be read
	qChan      chan queryAndError
}

type queryAndError struct {
	q   *Query
	err error
}

// New create a trc parser
func New(r io.Reader, name string, rowsToRead int) *Parser {
	p := &Parser{
		p:          trc.New(r, name),
		qChan:      make(chan queryAndError),
		rowsToRead: rowsToRead,
	}

	go func() {
		for fn := waitQuery; fn != nil; {
			fn = fn(p)
		}
		close(p.qChan)
	}()
	return p
}

// NextPacket deliver each packet until EOF or error
func (p *Parser) Next() (*Query, error) {
	r := <-p.qChan
	return r.q, r.err
}

// stateFn is the function that handle a state
type stateFn func(p *Parser) stateFn

// waitQuery wait a packet sent by the client with a query
func waitQuery(p *Parser) stateFn {
	for {
		pk, err := p.p.NextPacket()
		if pk == nil && err == nil {
			break
		}
		if pk.Typ == "nsbasic_bsd" {
			return p.parseQuery(pk)
		}
		if err != nil {
			fmt.Println(err)
		}
	}
	return nil
}

// waitQuery wait a packet sent by the client with a query
func waitResponse(p *Parser) stateFn {
	for {
		pk, err := p.p.NextPacket()
		if pk == nil && err == nil {
			break
		}
		if pk.Typ == "nsbasic_brc" {
			return p.parseResponse(pk)
		}
		if err != nil {
			fmt.Println(err)
		}
	}
	return nil
}

// parseQuery and returns the next stateFn
func (p *Parser) parseQuery(pk *trc.Packet) stateFn {

	var err error
	var b byte

	q := &Query{
		Packet: pk,
	}

	if len(pk.Payload) < 0x20 {
		return waitQuery
	}

	pl := pk.Payload
	buff := bytes.NewBuffer(pl[0x12:])
	b, err = buff.ReadByte()

	// Field 1 Marker
	if err == nil && b == 0x03 {
		b, err = buff.ReadByte()
		if err != nil || b != 0x5E {
			return waitQuery
		}
	} else {
		return waitQuery
	}
	if err != nil {
		return waitQuery
	}

	// Field 2 Flag
	b, err = buff.ReadByte() // discard Flag
	if err != nil {
		return waitQuery
	}

	// Field 3 ExeOp
	q.ExeOp, err = GetUInt(buff, 4, true, true) // Read ExeOp
	if err != nil {
		return waitQuery
	}

	// Field 4 Cursor ID
	q.CursorId, err = GetUInt(buff, 2, true, true) // Read Cursor ID
	if err != nil {
		return waitQuery
	}

	// Field 5
	b, err = buff.ReadByte() // discard byte after cursor id
	if err != nil {
		return waitQuery
	}

	// Field 6, statment length ?
	q.Len, err = GetUInt(buff, 4, true, true) // Read Statment length ??
	if err != nil {
		return waitQuery
	}

	// Field 7
	b, err = buff.ReadByte() // discard byte after statment length
	if err != nil {
		return waitQuery
	}

	var discardedInt int32
	// Field 8 always 13???
	discardedInt, err = GetInt(buff, 2, true, true) // Is always 13, purpose?.
	if err != nil {
		return waitQuery
	}

	// Field 9
	b, err = buff.ReadByte()
	if err != nil {
		return waitQuery
	}

	// Field 10
	b, err = buff.ReadByte()
	if err != nil {
		return waitQuery
	}

	// Field 11
	discardedInt, err = GetInt(buff, 4, true, true)
	if err != nil {
		return waitQuery
	}

	// Field 12
	q.RowToFetch, err = GetUInt(buff, 4, true, true) // row to fetch
	if err != nil {
		return waitQuery
	}

	// Field 13
	discardedInt, err = GetInt(buff, 4, true, true) // Should be 0, unknown
	if err != nil {
		return waitQuery
	}

	// Field 14 Has Parameters == 1
	b, err = buff.ReadByte() // Paramter flag
	if err != nil {
		return waitQuery
	}

	if b > 0 {
		// Field 15
		q.ParamLen, err = GetUInt(buff, 2, true, true)
		if err != nil {
			return waitQuery
		}
	}

	// Detect query start... sort of.

	for {
		b, err = buff.ReadByte()
		if err != nil {
			return waitQuery
		}
		if b > 5 {
			// One step back
			err = buff.UnreadByte()
			if err != nil {
				return waitQuery
			}
			break
		}
	}

	var stmt []byte
	stmt, err = readBytes(buff)
	if err != nil {
		return waitQuery
	}
	q.Query = string(stmt)

	// Skip 13 int for structure AL8I4
	for i := 0; i < 13; i++ {
		discardedInt, err = GetInt(buff, 2, true, true)
		if err != nil {
			return waitQuery
		}
	}

	if q.ParamLen > 0 {
		q.Params = []*ParameterInfo{}
		for i := 0; i < int(q.ParamLen); i++ {
			var p *ParameterInfo
			p, err = GetParamInfo(buff)
			if err != nil {
				return waitQuery
			}
			q.Params = append(q.Params, p)
		}

		// Skip byte 7
		b, err = buff.ReadByte()
		if err != nil {
			return waitQuery
		}

		for _, p := range q.Params {
			var v []byte
			v, err = readBytes(buff)
			if err != nil {
				return waitQuery
			}
			p.Value = v
		}
	}

	p.qChan <- queryAndError{
		q:   q,
		err: err,
	}

	_ = discardedInt
	if p.rowsToRead == 0 {
		return waitQuery
	}
	return waitResponse
}

type EndianNess int

var queryKeyWords = [][]byte{
	[]byte("SELECT"),
	[]byte("INSERT"),
	[]byte("UPDATE"),
	[]byte("MERGE"),
	[]byte("DELETE"),
	[]byte("ALTER"),
	[]byte("WIDTH"),
	[]byte("DECLARE"),
	[]byte("BEGIN"),
}

func toUpperAscii(b []byte) []byte {
	o := make([]byte, len(b))
	copy(o, b)
	for i := 0; i < len(o); i++ {
		if o[i] >= 'a' && o[i] <= 'z' {
			o[i] -= 'a' - 'A'
		}
	}
	return o
}

// writeEol add EOF marker to the string builder according running OS
func writeEol(sb *strings.Builder) {
	if runtime.GOOS == "WINDOWS" {
		sb.Write([]byte{0x0a, 0x0b})
		return
	}
	sb.WriteByte('\n')
}

func (p *Parser) waitResponse(pk *trc.Packet) stateFn {
	for {
		pk, err := p.p.NextPacket()
		if pk == nil && err == nil {
			break
		}
		if pk.Typ == "nsbasic_brc" {
			return p.parseResponse(pk)
		}
		if err != nil {
			fmt.Println(err)
		}
	}
	return nil
}

// parseQuery and returns the next stateFn
func (p *Parser) parseResponse(pk *trc.Packet) stateFn {

	var err error
	var b byte

	q := &Response{
		Packet: pk,
	}
}

type Response struct {
	Packet *trc.Packet // Query's packet
}
