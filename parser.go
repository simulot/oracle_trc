package main

import (
	"bufio"
	"bytes"
	"io"
	"log"
)

type packetType int

const (
	nsbasic_bsd packetType = iota
	nsbasic_brc
)

type Packet struct {
	pid     int
	t       packetType
	ts      string
	len     int
	payload []byte
}

type parser struct {
	s      *bufio.Scanner
	buff   bytes.Buffer
	pkChan chan *Packet
}

func New(r io.Reader) *parser {
	p := &parser{
		s:      bufio.NewScanner(r),
		pkChan: make(chan *Packet),
	}

	go func() {
		for fn := waitNSBasicBSD; fn != nil; {
			fn = fn(p)
		}
		close(p.pkChan)
	}()
	return p
}

func (p *parser) NextPacket() *Packet {
	return <-p.pkChan
}

type stateFn func(p *parser) stateFn

func waitNSBasicBSD(p *parser) stateFn {
	for p.s.Scan() {
		if bytes.HasSuffix(p.s.Bytes(), []byte("nsbasic_bsd: packet dump")) {
			p.buff.Reset()
			return inNSBasic
		}
	}
	if p.s.Err() != io.EOF {
		log.Println(p.s.Err())
	}
	return nil
}

func inNSBasic(p *parser) stateFn {
	var b []byte
	pk := &Packet{}
	for p.s.Scan() {
		b = p.s.Bytes()
		if bytes.HasSuffix(b, []byte("nsbasic_bsd: exit (0)")) {
			break
		}
		p.scanPacketLine(pk, p.s.Bytes())
	}
	if p.s.Err() != nil && p.s.Err() != io.EOF {
		log.Println(p.s.Err())
	}
	pk.payload = make([]byte, p.buff.Len())
	copy(pk.payload, p.buff.Bytes())
	p.pkChan <- pk
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
