package main

import (
	"reflect"
	"strings"
	"testing"
)

func Test_parser_NextPacket(t *testing.T) {
	tests := []struct {
		name string
		p    *parser
		want *Packet
	}{
		{
			name: "select",
			p: New(strings.NewReader(
				`
(6604) [12-FEB-2019 16:40:00:651] nsbasic_bsd: packet dump
(6604) [12-FEB-2019 16:40:00:651] nsbasic_bsd: 00 C3 00 00 06 00 00 00  |........|
(6604) [12-FEB-2019 16:40:00:651] nsbasic_bsd: 00 00 11 69 0D 01 01 01  |...i....|
(6604) [12-FEB-2019 16:40:00:651] nsbasic_bsd: 01 01 03 5E 0E 02 80 61  |...^...a|
(6604) [12-FEB-2019 16:40:00:651] nsbasic_bsd: 00 01 02 01 6B 01 01 0D  |....k...|
(6604) [12-FEB-2019 16:40:00:651] nsbasic_bsd: 01 01 00 01 64 00 00 00  |....d...|
(6604) [12-FEB-2019 16:40:00:651] nsbasic_bsd: 00 01 00 01 01 01 00 00  |........|
(6604) [12-FEB-2019 16:40:00:651] nsbasic_bsd: 01 01 00 00 00 00 00 FE  |........|
(6604) [12-FEB-2019 16:40:00:651] nsbasic_bsd: 40 73 65 6C 65 63 74 20  |@select.|
(6604) [12-FEB-2019 16:40:00:651] nsbasic_bsd: 43 4F 55 4E 54 28 2A 29  |COUNT(*)|
(6604) [12-FEB-2019 16:40:00:651] nsbasic_bsd: 20 66 72 6F 6D 20 75 73  |.from.us|
(6604) [12-FEB-2019 16:40:00:651] nsbasic_bsd: 65 72 5F 74 61 62 5F 63  |er_tab_c|
(6604) [12-FEB-2019 16:40:00:651] nsbasic_bsd: 6F 6C 75 6D 6E 73 20 77  |olumns.w|
(6604) [12-FEB-2019 16:40:00:651] nsbasic_bsd: 68 65 72 65 20 74 61 62  |here.tab|
(6604) [12-FEB-2019 16:40:00:651] nsbasic_bsd: 6C 65 5F 6E 61 6D 65 20  |le_name.|
(6604) [12-FEB-2019 16:40:00:651] nsbasic_bsd: 3D 20 4E 27 42 57 5F 55  |=.N'BW_U|
(6604) [12-FEB-2019 16:40:00:651] nsbasic_bsd: 53 39 45 52 5F 41 55 54  |S9ER_AUT|
(6604) [12-FEB-2019 16:40:00:651] nsbasic_bsd: 48 45 4E 54 49 43 41 54  |HENTICAT|
(6604) [12-FEB-2019 16:40:00:651] nsbasic_bsd: 49 4F 4E 27 20 61 6E 64  |ION'.and|
(6604) [12-FEB-2019 16:40:00:651] nsbasic_bsd: 20 63 6F 6C 75 6D 6E 5F  |.column_|
(6604) [12-FEB-2019 16:40:00:651] nsbasic_bsd: 6E 61 6D 65 20 3D 20 4E  |name.=.N|
(6604) [12-FEB-2019 16:40:00:651] nsbasic_bsd: 27 55 41 55 5F 46 4F 52  |'UAU_FOR|
(6604) [12-FEB-2019 16:40:00:651] nsbasic_bsd: 47 4F 54 5F 41 43 54 49  |GOT_ACTI|
(6604) [12-FEB-2019 16:40:00:651] nsbasic_bsd: 4F 4E 27 00 01 01 00 00  |ON'.....|
(6604) [12-FEB-2019 16:40:00:651] nsbasic_bsd: 00 00 00 00 01 01 00 00  |........|
(6604) [12-FEB-2019 16:40:00:651] nsbasic_bsd: 00 00 00                 |...     |
(6604) [12-FEB-2019 16:40:00:651] nsbasic_bsd: exit (0)
(6604) [12-FEB-2019 16:40:00:651] nsbasic_brc: entry: oln/tot=0
(6604) [12-FEB-2019 16:40:00:651] nttfprd: entry
(6604) [12-FEB-2019 16:40:00:651] nttfprd: socket 836 had bytes read=183
(6604) [12-FEB-2019 16:40:00:651] nttfprd: exit
`)),
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.NextPacket(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parser.NextPacket() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
