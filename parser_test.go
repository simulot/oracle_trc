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
`), "test"),
			want: nil,
		},
		{
			name: "select",
			p: New(strings.NewReader(
				`
				(9520) [26-FEB-2019 11:52:23:665] nsbasic_bsd: packet dump
				(9520) [26-FEB-2019 11:52:23:665] nsbasic_bsd: 01 4F 00 00 06 00 00 00  |.O......|
				(9520) [26-FEB-2019 11:52:23:665] nsbasic_bsd: 00 00 11 69 11 01 01 01  |...i....|
				(9520) [26-FEB-2019 11:52:23:665] nsbasic_bsd: 01 01 03 5E 12 02 81 29  |...^...)|
				(9520) [26-FEB-2019 11:52:23:665] nsbasic_bsd: 00 01 02 01 C5 01 01 0D  |........|
				(9520) [26-FEB-2019 11:52:23:665] nsbasic_bsd: 01 01 00 01 01 00 01 01  |........|
				(9520) [26-FEB-2019 11:52:23:665] nsbasic_bsd: 04 00 01 00 01 01 01 00  |........|
				(9520) [26-FEB-2019 11:52:23:665] nsbasic_bsd: 00 01 01 00 00 00 00 00  |........|
				(9520) [26-FEB-2019 11:52:23:665] nsbasic_bsd: FE 40 49 4E 53 45 52 54  |.@INSERT|
				(9520) [26-FEB-2019 11:52:23:665] nsbasic_bsd: 20 49 4E 54 4F 20 42 57  |.INTO.BW|
				(9520) [26-FEB-2019 11:52:23:665] nsbasic_bsd: 5F 55 53 45 52 5F 41 55  |_USER_AU|
				(9520) [26-FEB-2019 11:52:23:665] nsbasic_bsd: 54 48 5F 4C 4F 47 20 28  |TH_LOG.(|
				(9520) [26-FEB-2019 11:52:23:665] nsbasic_bsd: 55 41 4C 5F 49 44 2C 55  |UAL_ID,U|
				(9520) [26-FEB-2019 11:52:23:665] nsbasic_bsd: 41 4C 5F 55 53 45 52 5F  |AL_USER_|
				(9520) [26-FEB-2019 11:52:23:665] nsbasic_bsd: 4E 45 54 57 4F 52 4B 5F  |NETWORK_|
				(9520) [26-FEB-2019 11:52:23:665] nsbasic_bsd: 4E 41 4D 45 2C 55 41 4C  |NAME,UAL|
				(9520) [26-FEB-2019 11:52:23:665] nsbasic_bsd: 5F 4C 40 4F 47 5F 43 4F  |_L@OG_CO|
				(9520) [26-FEB-2019 11:52:23:665] nsbasic_bsd: 44 45 2C 55 41 4C 5F 53  |DE,UAL_S|
				(9520) [26-FEB-2019 11:52:23:665] nsbasic_bsd: 54 41 4D 50 5F 44 41 54  |TAMP_DAT|
				(9520) [26-FEB-2019 11:52:23:665] nsbasic_bsd: 45 2C 55 41 4C 5F 41 50  |E,UAL_AP|
				(9520) [26-FEB-2019 11:52:23:665] nsbasic_bsd: 50 4C 49 43 41 54 49 4F  |PLICATIO|
				(9520) [26-FEB-2019 11:52:23:665] nsbasic_bsd: 4E 29 20 56 41 4C 55 45  |N).VALUE|
				(9520) [26-FEB-2019 11:52:23:665] nsbasic_bsd: 53 20 28 3A 31 2C 3A 32  |S.(:1,:2|
				(9520) [26-FEB-2019 11:52:23:665] nsbasic_bsd: 2C 3A 33 2C 28 53 45 4C  |,:3,(SEL|
				(9520) [26-FEB-2019 11:52:23:665] nsbasic_bsd: 45 43 54 17 20 53 59 53  |ECT..SYS|
				(9520) [26-FEB-2019 11:52:23:665] nsbasic_bsd: 44 41 54 45 20 46 52 4F  |DATE.FRO|
				(9520) [26-FEB-2019 11:52:23:665] nsbasic_bsd: 4D 20 44 55 41 4C 29 2C  |M.DUAL),|
				(9520) [26-FEB-2019 11:52:23:665] nsbasic_bsd: 3A 34 29 00 01 01 01 01  |:4).....|
				(9520) [26-FEB-2019 11:52:23:665] nsbasic_bsd: 00 00 00 00 00 01 04 00  |........|
				(9520) [26-FEB-2019 11:52:23:665] nsbasic_bsd: 00 00 00 00 60 01 00 00  |....x...|
				(9520) [26-FEB-2019 11:52:23:665] nsbasic_bsd: 01 78 00 01 10 00 00 02  |.x......|
				(9520) [26-FEB-2019 11:52:23:666] nsbasic_bsd: 03 69 01 00 60 01 00 00  |.i..s...|
				(9520) [26-FEB-2019 11:52:23:666] nsbasic_bsd: 01 36 00 01 10 00 00 02  |.6......|
				(9520) [26-FEB-2019 11:52:23:666] nsbasic_bsd: 03 69 01 00 02 01 00 00  |.i......|
				(9520) [26-FEB-2019 11:52:23:666] nsbasic_bsd: 01 16 00 00 00 00 00 00  |........|
				(9520) [26-FEB-2019 11:52:23:666] nsbasic_bsd: 00 60 01 00 00 01 3C 00  |.s....<.|
				(9520) [26-FEB-2019 11:52:23:666] nsbasic_bsd: 01 10 00 00 02 03 69 01  |......i.|
				(9520) [26-FEB-2019 11:52:23:666] nsbasic_bsd: 00 07 14 44 38 33 42 36  |...D83B6|
				(9520) [26-FEB-2019 11:52:23:666] nsbasic_bsd: 30 30 36 38 31 39 42 34  |006819B4|
				(9520) [26-FEB-2019 11:52:23:666] nsbasic_bsd: 46 30 33 42 31 45 43 09  |F03B1EC.|
				(9520) [26-FEB-2019 11:52:23:666] nsbasic_bsd: 4A 46 2E 43 41 53 53 41  |JF.CASSA|
				(9520) [26-FEB-2019 11:52:23:666] nsbasic_bsd: 4E 02 C1 02 0A 4D 61 73  |N....Mas|
				(9520) [26-FEB-2019 11:52:23:666] nsbasic_bsd: 74 65 72 20 35 2E 31     |ter.5.1 |
				(9520) [26-FEB-2019 11:52:23:666] nsbasic_bsd: exit (0)
				(9520) [26-FEB-2019 11:52:23:666] nsbasic_brc: entry: oln/tot=0
				(9520) [26-FEB-2019 11:52:23:666] nttfprd: entry
`), "test"),
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := tt.p.NextPacket(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parser.NextPacket() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
