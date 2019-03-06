package packet

import (
	"fmt"
	"strings"
)

/*
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

type TNSHeader struct {
	Length         int16
	CheckSum       int16
	PkType         PacketType
	Flags          int8
	HeaderCheckSum int16
}

func (th TNSHeader) String() string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("PktLen(%d),", th.Length))
	sb.WriteString(fmt.Sprintf("Chksum(%04x),", th.CheckSum))
	sb.WriteString(fmt.Sprintf("PkType(%d,%s),", th.PkType, packetTypeStr[th.PkType]))
	sb.WriteString(fmt.Sprintf("Flags(%04x),", th.Flags))
	sb.WriteString(fmt.Sprintf("HdrChkSum(%04x),", th.HeaderCheckSum))
	return sb.String()
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

type TNSSQLDataPkt struct {
	TNSHeader
	DataFlag      int16 // 0000 DATA, 0020 MORE (Oracle 12c), 0040 EOF
	TTCCode       int8
	Function      int8
	PacketCounter int8
	CID           int8 // cid
	OPESIZ        int8
}
