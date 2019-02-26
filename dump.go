package main

import (
	"fmt"
	"strings"
)

func (p *parser) dumpPackets() {

	for pk := range p.pkChan {
		fmt.Println(pk.pid, pk.ts)
		fmt.Println(dump(pk.payload))
	}
}

var hexC = [16]byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'A', 'B', 'C', 'D', 'E', 'F'}

func dump(b []byte) string {
	sb := strings.Builder{}
	for r := 0; r < len(b); r += 16 {
		sb.WriteString(fmt.Sprintf("0x%04X: ", r))
		for p := 0; p < 16; p++ {
			if r+p < len(b) {
				c := b[r+p]
				sb.WriteByte(hexC[(c&0xf0)>>4])
				sb.WriteByte(hexC[(c & 0x0f)])
				sb.WriteByte(' ')
			} else {
				sb.WriteString("   ")
			}
			if p == 7 {
				sb.WriteString("- ")
			}
		}
		sb.WriteString(" | ")
		for p := 0; p < 16; p++ {
			if r+p < len(b) {
				c := b[r+p]
				if c > ' ' && c <= '~' {
					sb.WriteByte(c)
				} else {
					sb.WriteByte('.')
				}
			} else {
				sb.WriteByte(' ')
			}
		}
		sb.WriteString(" |\n")
	}
	return sb.String()
}
