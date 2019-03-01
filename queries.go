package main

import (
	"bytes"
	"fmt"
	"strings"
	"time"
)

func (p *parser) dumpQueries(after time.Time) (err error) {

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

		disp := after.IsZero()
		if !disp {
			t, err := p.tParser(pk.ts)
			if err == nil {
				disp = t.After(after)
			}
		}
		if disp {
			fmt.Printf("%s(%d) %s(%d) %s ", p.name, pk.line, p.clients[pk.pid], pk.pid, pk.ts)
			fmt.Println(query(pl[pos:]))
		}
	}
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
