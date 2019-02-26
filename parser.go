package main

import (
	"bufio"
	"bytes"
	"io"

	"github.com/pkg/errors"
)

type packetType int

const (
	nsbasic_bsd packetType = iota
	nsbasic_brc
)

type Packet struct {
	pid     int
	ts      string
	line    int
	payload []byte
}

type packetAndError struct {
	*Packet
	err error
}

type parser struct {
	name   string
	line   int
	s      *bufio.Scanner
	buff   bytes.Buffer
	pkChan chan packetAndError
}

func New(r io.Reader, name string) *parser {
	p := &parser{
		name:   name,
		s:      bufio.NewScanner(r),
		pkChan: make(chan packetAndError),
	}

	go func() {
		for fn := waitNSBasicBSD; fn != nil; {
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

func waitNSBasicBSD(p *parser) stateFn {
	for p.Scan() {
		if bytes.HasSuffix(p.s.Bytes(), []byte("nsbasic_bsd: packet dump")) {
			p.buff.Reset()
			return inNSBasic
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
	// if p.line == 5422 {
	// 	runtime.Breakpoint()
	// }
	p.buff.Reset()
	var b []byte
	pk := &Packet{
		line: p.line,
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
	return waitNSBasicBSD
}

func (p *parser) scanPacketLine(pk *Packet, b []byte) {
	if len(b) == 0 {
		return
	}
	b = b[1:] // Skip '('
	i := 0
	if pk.pid == 0 {
		for i = 0; i < len(b) && b[i] != ')'; i++ {
			pk.pid = pk.pid*10 + int(b[i]-'0')
		}
	} else {
		for i = 0; i < len(b) && b[i] != ')'; i++ {
		}
	}
	if len(b) < i+3 {
		return
	}
	b = b[i+3:]
	for i = 0; i < len(b) && b[i] != ']'; i++ {
	}
	if len(pk.ts) == 0 && b[i] == ']' {
		pk.ts = string(b[0:i])
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
