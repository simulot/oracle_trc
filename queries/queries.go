package queries

import (
	"bytes"
	"fmt"
	"io"
	"runtime"
	"strings"

	"github.com/simulot/oracle_trc/trc"
)

// Packet hold packet content and context of packet
type Query struct {
	Query  string      // Query text
	Packet *trc.Packet // Query's packet
}

// String implement the basic representation of packet: Packet's context and its content in hexadecimal
func (q Query) String() string {
	sb := strings.Builder{}
	q.Packet.WriteContext(&sb)
	writeEol(&sb)
	sb.WriteString(q.Query)
	writeEol(&sb)
	return sb.String()
}

// Parser is used to parse trc files and extract queries
type Parser struct {
	p     *trc.Parser // Trace file parser
	q     *Query      // current query
	qChan chan queryAndError
}

type queryAndError struct {
	q   *Query
	err error
}

// New create a trc parser
func New(r io.Reader, name string) *Parser {
	p := &Parser{
		p:     trc.New(r, name),
		qChan: make(chan queryAndError),
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

// parseQuery and returns the next stateFn
func (p *Parser) parseQuery(pk *trc.Packet) stateFn {
	payload := pk.Payload
	pos := detectQuery(payload)
	if pos < 0 {
		return waitQuery
	}

	q := &Query{
		Packet: pk,
		Query:  query(payload[pos:]),
	}
	p.qChan <- queryAndError{
		q:   q,
		err: nil,
	}
	return waitQuery
}

var queryKeyWords = [][]byte{
	[]byte("SELECT"),
	[]byte("INSERT"),
	[]byte("UPDATE"),
	[]byte("MERGE"),
	[]byte("DELETE"),
	[]byte("ALTER"),
	[]byte("WIDTH"),
}

func detectQuery(pl []byte) int {
	l := len(pl)
	if l > 0x100 {
		l = 0x100
	}
	b := toUpperAscii(pl[:l])
	pos := -1

	for k := 0; k < len(queryKeyWords); k++ {
		p := bytes.Index(b, queryKeyWords[k])
		if pos == -1 || (p > 0 && p < pos) {
			pos = p
		}
	}

	if pos == -1 {
		return pos
	}
	// Check white chars prepending the query
	for pos > 1 {
		switch b[pos-1] {
		case ' ', '\t', '\r', '\n', '(':
			pos--
		default:
			return pos - 1
		}
	}
	return -1
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

func query(b []byte) string {
	sb := strings.Builder{}
	for len(b) > 0 {
		l := int(b[0])
		if l == 0 || l == 1 || l > len(b) {
			break
		}
		b = b[1:]
		if l > 0 {
			if l > len(b) {
				l = len(b)
			}
			sb.Write(b[:l])
			b = b[l:]
		}
	}
	return sb.String()
}

// writeEol add EOF marker to the string builder according running OS
func writeEol(sb *strings.Builder) {
	if runtime.GOOS == "WINDOWS" {
		sb.Write([]byte{0x0a, 0x0b})
		return
	}
	sb.WriteByte('\n')
}
