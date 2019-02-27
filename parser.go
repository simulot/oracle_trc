package main

import (
	"bufio"
	"bytes"
	"io"
	"path/filepath"
	"time"

	"github.com/simulot/oracle_trc/ts"

	"github.com/pkg/errors"
)

type packetType int

const (
	nsbasic_bsd packetType = iota
	nsbasic_brc
)

type Packet struct {
	pid     int
	ts      []byte
	t       time.Time
	line    int
	payload []byte
}

type packetAndError struct {
	*Packet
	err error
}

type parser struct {
	name    string
	line    int
	s       *bufio.Scanner
	buff    bytes.Buffer
	pkChan  chan packetAndError
	tParser ts.TimeParserFn
	clients map[int]string
}

func New(r io.Reader, name string, tParser ts.TimeParserFn) *parser {
	p := &parser{
		name:    name,
		s:       bufio.NewScanner(r),
		pkChan:  make(chan packetAndError),
		tParser: tParser,
		clients: make(map[int]string),
	}

	go func() {
		for fn := waitInterstingLines; fn != nil; {
			fn = fn(p)
		}
		close(p.pkChan)
	}()
	return p
}

func (p *parser) NextPacket() (*Packet, error) {
	pk := <-p.pkChan
	return pk.Packet, pk.err
}

func (p *parser) Scan() bool {
	b := p.s.Scan()
	p.line++
	return b
}

type stateFn func(p *parser) stateFn

func waitInterstingLines(p *parser) stateFn {
	for p.Scan() {
		if bytes.HasSuffix(p.s.Bytes(), []byte("nsbasic_bsd: packet dump")) {
			return inNSBasic
		}
		if bytes.Contains(p.s.Bytes(), []byte("nsc2addr:")) {
			return inNSC2Addr
		}
	}
	if p.s.Err() != nil && p.s.Err() != io.EOF {
		p.pkChan <- packetAndError{
			Packet: nil,
			err:    errors.Wrapf(p.s.Err(), "%s(%d)", p.name, p.line),
		}
	}
	return nil
}

func inNSBasic(p *parser) stateFn {
	p.buff.Reset()
	var b []byte
	pk := &Packet{
		line: p.line,
		ts:   nil,
	}
	for p.Scan() {
		b = p.s.Bytes()
		if bytes.HasSuffix(b, []byte("nsbasic_bsd: exit (0)")) {
			break
		}
		p.scanPacketLine(pk, p.s.Bytes())
	}
	if p.s.Err() != nil && p.s.Err() != io.EOF {
		p.pkChan <- packetAndError{
			Packet: nil,
			err:    errors.Wrapf(p.s.Err(), "%s(%d)", p.name, p.line),
		}
	}
	pk.payload = make([]byte, p.buff.Len())
	copy(pk.payload, p.buff.Bytes())
	p.pkChan <- packetAndError{
		Packet: pk,
		err:    nil,
	}
	return waitInterstingLines
}

func (p *parser) scanPID(b []byte) (int, []byte) {
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

func (p *parser) scanPacketLine(pk *Packet, b []byte) {
	if pk.pid == 0 {
		pk.pid, b = p.scanPID(b)
	}

	i := 0
	for i = 0; i < len(b); i++ {
		if b[i] == '[' {
			i++
			break
		}
	}
	b = b[i:]
	for i = 0; i < len(b) && b[i] != ']'; i++ {
	}
	if len(pk.ts) == 0 && b[i] == ']' {
		pk.ts = make([]byte, i)
		copy(pk.ts, b[0:i])
	}

	if i >= len(b) {
		return
	}
	// skip nsbasic_bsd:
	b = b[i+1:]
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
		p.buff.WriteByte(hex(b[i : i+2]))
	}

}

func hex(buf []byte) byte {
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

func inNSC2Addr(p *parser) stateFn {
	var b []byte

	for p.Scan() {
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
				p.clients[pid] = filepath.Base(string(b[:pos]))
			}
		}
	}
	if p.s.Err() != nil && p.s.Err() != io.EOF {
		p.pkChan <- packetAndError{
			Packet: nil,
			err:    errors.Wrapf(p.s.Err(), "%s(%d)", p.name, p.line),
		}
	}
	return waitInterstingLines
}
