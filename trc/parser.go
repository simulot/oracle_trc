package trc

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"runtime"
	"strings"

	"github.com/pkg/errors"
)

// Packet hold packet content and context of packet
type Packet struct {
	Name    string // Trace file name
	Typ     string // Packet type as read in trc
	Line    int    // Line number in trace file
	Pid     int    // Client PID
	Client  string // Client name
	TS      []byte // Event time as written in trc file
	Payload []byte // Packet content
}

type PacketDirection int

// Parser is used to parse trc files and extract data packets
type Parser struct {
	name       string              // Source name, used for reporting errors
	line       int                 // current line in trc file
	s          *bufio.Scanner      // text scanner
	buff       bytes.Buffer        // buffer used for gathering packet fragments
	pkChan     chan packetAndError // gather extracted packets
	clients    map[int]string      // Hold client names per PID
	packetType string              // current packet type as seen in trc file
	pk         *Packet             // current packet
}

// String implement the basic representation of packet: Packet's context and its content in hexadecimal
func (pk Packet) String() string {
	sb := strings.Builder{}
	pk.context(&sb)
	writeEol(&sb)
	pk.hexdump(&sb)
	writeEol(&sb)
	return sb.String()
}

type packetAndError struct {
	pk  *Packet
	err error
}

// New create a trc parser
func New(r io.Reader, name string) *Parser {
	name = baseName(name)
	p := &Parser{
		s:          bufio.NewScanner(r),
		pkChan:     make(chan packetAndError),
		clients:    make(map[int]string),
		name:       name,
		packetType: "",
	}

	go func() {
		for fn := waitInterstingLines; fn != nil; {
			fn = fn(p)
		}
		close(p.pkChan)
	}()
	return p
}

// NextPacket deliver each packet until EOF or error
func (p *Parser) NextPacket() (*Packet, error) {
	pk := <-p.pkChan
	return pk.pk, pk.err
}

// scan wrap the standard scan for counting lines
func (p *Parser) scan() bool {
	b := p.s.Scan()
	p.line++
	return b
}

type stateFn func(p *Parser) stateFn

func waitInterstingLines(p *Parser) stateFn {
	for p.scan() {
		if bytes.HasSuffix(p.s.Bytes(), []byte("packet dump")) {
			return p.determinePacketType(p.s.Bytes())
		}
		if bytes.Contains(p.s.Bytes(), []byte("nsc2addr:")) {
			return inNSC2Addr
		}
	}
	if p.s.Err() != nil && p.s.Err() != io.EOF {
		p.pkChan <- packetAndError{
			pk:  nil,
			err: errors.Wrapf(p.s.Err(), "Error reading %s at line %d", p.name, p.line),
		}
	}
	return nil
}

func (p *Parser) determinePacketType(b []byte) stateFn {
	var i = 0
	for i = 0; i < len(b) && b[i] != ']'; i++ {
	}
	i += 2
	b = b[i:]
	i = bytes.Index(b, []byte{':'})
	if i > 0 {
		p.packetType = string(b[:i])
	}
	p.pk = &Packet{
		Name: p.name,
		Line: p.line,
		TS:   nil,
		Typ:  p.packetType,
	}

	return inDumpPacket
}

// inDumpPacket scan nsbasic lines then yield to waitInterstingLines
func inDumpPacket(p *Parser) stateFn {
	p.buff.Reset()
	var packetType = []byte(p.pk.Typ)
	var b []byte

	for p.scan() {
		b = p.s.Bytes()
		if !bytes.Contains(b, packetType) {
			break
		}
		p.scanPacketLine(p.s.Bytes())
	}
	if p.s.Err() != nil && p.s.Err() != io.EOF {
		p.pkChan <- packetAndError{
			pk:  nil,
			err: errors.Wrapf(p.s.Err(), "Error reading %s at line %d", p.name, p.line),
		}
	}
	p.pk.Payload = make([]byte, p.buff.Len())
	copy(p.pk.Payload, p.buff.Bytes())
	p.pkChan <- packetAndError{
		pk:  p.pk,
		err: nil,
	}
	return waitInterstingLines
}

// scanPID get PID from scanned line
func (p *Parser) scanPID(b []byte) (int, []byte) {
	pid := 0
	if len(b) == 0 {
		return 0, b
	}

	i := 0
	for i = 0; i < len(b); i++ {
		if b[i] == '(' {
			i++
			break
		}
	}
	b = b[i:] // Skip '('

	for i = 0; i < len(b) && b[i] != ')'; i++ {
		pid = pid*10 + int(b[i]-'0')
	}
	b = b[i:]
	return pid, b
}

// scanPacketLine scan one of packets lines
func (p *Parser) scanPacketLine(b []byte) {
	if p.pk.Pid == 0 {
		// Get PID and client
		p.pk.Pid, b = p.scanPID(b)
		p.pk.Client = p.clients[p.pk.Pid]
	}

	// Go to time stamp begin
	i := 0
	for i = 0; i < len(b); i++ {
		if b[i] == '[' {
			i++
			break
		}
	}

	// Determine time stop
	b = b[i:]
	for i = 0; i < len(b) && b[i] != ']'; i++ {
	}
	if len(p.pk.TS) == 0 && b[i] == ']' {
		p.pk.TS = make([]byte, i)
		copy(p.pk.TS, b[0:i])
	}

	if i >= len(b) {
		return
	}

	// skip packet type
	b = b[i+2:]
	for i = 0; i < len(b) && b[i] != ':'; i++ {
	}

	if i >= len(b) {
		return
	}
	b = b[i+2:] // At beging of HEX
	for i = 0; i < len(b); i += 3 {
		if b[i] == ' ' {
			break
		}
		p.buff.WriteByte(hexToByte(b[i : i+2])) // Accumulate decoded bytes
	}

}

// convert 2 hex chars into a byte
func hexToByte(buf []byte) byte {
	b := byte(0)
	for i := 0; i < 2; i++ {
		c := buf[i]
		switch {
		case c >= 'A' && c <= 'F':
			b = b<<4 + c + 10 - 'A'
		case c >= '0' && c <= '9':
			b = b<<4 + c - '0'
		}
	}
	return b
}

// inNSC2Addr extract client program
func inNSC2Addr(p *Parser) stateFn {
	var b []byte

	for p.scan() {
		b = p.s.Bytes()
		if bytes.HasSuffix(b, []byte("nsc2addr: normal exit")) {
			break
		}
		pos := bytes.Index(b, []byte("PROGRAM"))
		if pos >= 0 {
			pid, _ := p.scanPID(b)

			b = b[pos+len("PROGRAM="):]
			pos = bytes.Index(b, []byte(")"))
			if pos >= 0 {
				p.clients[pid] = baseName(string(b[:pos]))
			}
		}
	}
	if p.s.Err() != nil && p.s.Err() != io.EOF {
		p.pkChan <- packetAndError{
			pk:  nil,
			err: errors.Wrapf(p.s.Err(), "Error reading %s at line %d", p.name, p.line),
		}
	}
	return waitInterstingLines
}

// context write packet context to the string builder
func (pk Packet) context(sb *strings.Builder) {
	sb.WriteString(pk.Name)
	sb.WriteString(fmt.Sprintf("(%d),", pk.Line))
	sb.Write(pk.TS)
	sb.WriteByte(',')
	sb.WriteString(pk.Client)
	sb.WriteString(fmt.Sprintf("(%d),", pk.Pid))
	sb.WriteString(pk.Typ)
	sb.WriteByte(':')
}

// hexDump
func (pk Packet) hexdump(sb *strings.Builder) {
	sb.WriteString(hex.Dump(pk.Payload))
}

// writeEol add EOF marker to the string builder according running OS
func writeEol(sb *strings.Builder) {
	if runtime.GOOS == "WINDOWS" {
		sb.Write([]byte{0x0a, 0x0b})
		return
	}
	sb.WriteByte('\n')
}

// baseName is filepath.base using whatever seprarator '/' or '\\'
func baseName(s string) string {
	for i := len(s) - 1; i >= 0; {
		c := s[i]
		if c != '\\' && c != '/' {
			i--
			continue
		}
		return s[i+1:]
	}
	return s
}
