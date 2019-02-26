package main

import (
	"bytes"
	"fmt"
	"strings"
)

func (p *parser) dumpQueries() (err error) {

	for {
		pk, err := p.NextPacket()
		if err != nil {
			return err
		}
		if pk == nil {
			return nil
		}
		pl := pk.payload
		pos := detectQuery(pl)
		if pos < 0 {
			continue
		}

		fmt.Printf("%s(%d) (%d) %s ", p.name, pk.line, pk.pid, pk.ts)
		fmt.Println(query(pl[pos:]))
	}
}

var queryKeyWords = [][]byte{
	[]byte("SELECT"),
	[]byte("INSERT"),
	[]byte("UPDATE"),
	[]byte("MERGE"),
	[]byte("DELETE"),
	[]byte("ALTER"),
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
		if pos == -1 || p > 0 {
			pos = p
		}
	}

	if pos == -1 {
		return pos
	}
	// Check white chars prepending the query
	for pos > 1 {
		switch b[pos-1] {
		case ' ', '\t', '\r', '\n':
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
	for {
		l := int(b[0])
		if l == 0 || l == 1 || l > len(b) {
			break
		}
		if l > 0 && l <= len(b) {
			sb.Write(b[1 : l+1])
			b = b[l+1:]
		}
	}
	return sb.String()
}
