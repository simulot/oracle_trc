package packet

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"runtime"
	"strings"
)

/*
	Google findings on TNS packets
	https://blog.pythian.com/repost-oracle-protocol/
	https://flylib.com/books/en/2.680.1/the_oracle_network_architecture.html
	https://github.com/wireshark/wireshark/blob/master/epan/dissectors/packet-tns.c

*/

type PacketType uint8

const (
	Connect PacketType = iota + 1
	Accept
	Ack
	Refuse
	Redirect
	Data
	NULL
	Abort
	Resend
	Marker
	Attention
	Control
)

var packetTypeStr = map[PacketType]string{
	Connect:   "Connect",
	Accept:    "Accept",
	Ack:       "Ack",
	Refuse:    "Refuse",
	Redirect:  "Redirect",
	Data:      "Data",
	NULL:      "NULL",
	Abort:     "Abort",
	Resend:    "Resend",
	Marker:    "Marker",
	Attention: "Attention",
	Control:   "Control",
}

type StringBuilder interface {
	StringBuilder(*strings.Builder)
}

type Stringer interface {
	String() string
}

func PacketStringer(b []byte) Stringer {
	t := b[4]
	switch t {
	case 6:
		return ReadTNSData(b)
	default:
		return ReadTNSHeader(b)
	}

}

// func PacketStringBuilder(b []byte) StringBuilder {
// 	t := b[4]
// 	switch t {
// 	case 6:
// 		return ReadTNSData(b)
// 	default:
// 		return ReadTNSHeader(b)
// 	}

// }

type TNSHeader struct {
	Buffer         []byte
	Length         uint16     // ofs: 0
	CheckSum       uint16     // ofs: 2
	PkType         PacketType // ofs: 4
	Flags          uint8      // ofs: 5
	HeaderCheckSum uint16     // ofs: 6
}

func (th TNSHeader) String() string {
	sb := strings.Builder{}
	th.writeFields(&sb)
	writeEol(&sb)
	writePayload(&sb, th.Buffer)
	writeEol(&sb)
	return sb.String()
}

func (th TNSHeader) writeFields(sb *strings.Builder) {
	sb.WriteString("Packet header: ")
	sb.WriteString(fmt.Sprintf("PktLen(%d),", th.Length))
	sb.WriteString(fmt.Sprintf("Chksum(%04x),", th.CheckSum))
	sb.WriteString(fmt.Sprintf("PkType(%d=%s),", th.PkType, packetTypeStr[th.PkType]))
	sb.WriteString(fmt.Sprintf("Flags(%04x),", th.Flags))
	sb.WriteString(fmt.Sprintf("HdrChkSum(%04x)", th.HeaderCheckSum))
}

// func (th TNSHeader) StringBuilder(sb *strings.Builder) {
// 	th.writeFields(sb)
// 	writePayload(sb, th.Buffer[8:])
// }

func ReadTNSHeader(b []byte) TNSHeader {
	return TNSHeader{
		Buffer:         b,
		Length:         binary.BigEndian.Uint16(b[0:]),
		CheckSum:       binary.BigEndian.Uint16(b[2:]),
		PkType:         PacketType(uint8(b[4])),
		Flags:          uint8(b[5]),
		HeaderCheckSum: binary.BigEndian.Uint16(b[6:]),
	}
}

type TNSData struct {
	TNSHeader               // ofs: 0 Len 7
	DataFlag         uint16 // ofs: 8
	Function         uint8  // ofs: 10
	Sequence         uint8  // ofs: 11
	ExtendedFunction uint8  // ofs: 12
}

func ReadTNSData(b []byte) TNSData {
	return TNSData{
		TNSHeader:        ReadTNSHeader(b),
		DataFlag:         binary.BigEndian.Uint16(b[8:]),
		Function:         uint8(b[10]),
		Sequence:         uint8(b[11]),
		ExtendedFunction: uint8(b[12]),
	}
}
func (td TNSData) String() string {
	sb := strings.Builder{}
	td.TNSHeader.writeFields(&sb)
	writeEol(&sb)
	td.writeFields(&sb)
	writeEol(&sb)
	writePayload(&sb, td.TNSHeader.Buffer)
	return sb.String()
}

func (td TNSData) StringBuilder(sb *strings.Builder) {
	td.TNSHeader.writeFields(sb)
	writeEol(sb)
	if td.PkType == 6 {
		td.writeFields(sb)
		writeEol(sb)
	}
	writePayload(sb, td.TNSHeader.Buffer)
}

func (td TNSData) writeFields(sb *strings.Builder) {
	sb.WriteString("Data packet fields: ")
	sb.WriteString(fmt.Sprintf("DataFlag(%016b),", td.DataFlag))
	sb.WriteString(fmt.Sprintf("Function(%02x),", td.Function))
	sb.WriteString(fmt.Sprintf("Sequence(%d),", td.Sequence))
	sb.WriteString(fmt.Sprintf("ExtendedFunction(%02x),", td.ExtendedFunction))
}

/*

	Common Packet header
	00 5D: 		Packet length
	00 00: 		Checksum is always 0000
	06: 		Packet type: 6 DATA
	00:			Flags
	00 00:		Header checksum is always 0000

	Data Packet
	00 00:		Data flag (0x0040 for EOF)

	11:			11 for extended Two-Task Interface(TTI) function
	69:			Function ???
	4A:			Sequence number
	01 01 01: ???
	01 : 1 / No bind variable, 2 Bind variable
	02 03 5E 4B
	02 80 61 00
	01 01 45 01
	01 0D 01 01
	00 01 64 00
	00 00 00 01
	00 01 01 01
	00 00 01 01
	00 00 00 00
	00 : No IDEA 41 bytes
	17 : Length
	73 65 6C 65 63 74 20 64 6F 63 5F 69 64 20 66 72 6F 6D 20 44 4F 43 53 01: Query
	01 00 00 00 00 00 00 01 01 00 00 00 00 00 : Packet end.
*/

// writeEol add EOF marker to the string builder according running OS
func writeEol(sb *strings.Builder) {
	if runtime.GOOS == "WINDOWS" {
		sb.Write([]byte{0x0a, 0x0b})
		return
	}
	sb.WriteByte('\n')
}

func writePayload(sb *strings.Builder, b []byte) {
	// l := len(b)
	// if l > 0x100 {
	// 	sb.WriteString(hex.Dump(b[:0x100]))
	// 	sb.WriteString("... cut at 0x100 ...")
	// } else {
	sb.WriteString(hex.Dump(b))
	// }
}
