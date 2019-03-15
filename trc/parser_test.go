package trc

import (
	"errors"
	"reflect"
	"strings"
	"testing"
)

func getPacketFromTraceSnippet(trc string) (*Packet, error) {
	p := New(strings.NewReader(trc), "")

	for {
		pk, err := p.NextPacket()
		if pk == nil && err == nil {
			return nil, nil
		}
		return pk, err
	}
	return nil, errors.New("Should not happen")
}

func Test_getPacketFromTraceSnippet(t *testing.T) {
	type args struct {
		trc string
	}
	tests := []struct {
		name    string
		args    args
		want    *Packet
		wantErr bool
	}{

		{
			name: "simple with connection info and socket",
			args: args{
				trc: "(5304) [12-FEB-2019 17:25:10:647] nsc2addr: entry\r\n" +
					"(5304) [12-FEB-2019 17:25:10:647] nsc2addr: (DESCRIPTION=(CONNECT_DATA=(SID=BSWD01)(CID=(PROGRAM=C:\\App\\Service.exe)(HOST=APPSERVR)(USER=SYSTEM)))(ADDRESS=(PROTOCOL=TCP)(HOST=10.30.194.77)(PORT=1525)))\r\n" +
					"(5304) [12-FEB-2019 17:25:10:647] nttbnd2addr: entry\r\n" +
					"(5304) [12-FEB-2019 17:25:10:647] nsc2addr: normal exit\r\n" +
					"(5304) [12-FEB-2019 17:25:10:647] nsprecv: entry\r\n" +
					"(5304) [12-FEB-2019 17:25:10:647] nsprecv: reading frm transport...\r\n" +
					"(5304) [12-FEB-2019 17:25:10:647] nttrd: entry\r\n" +
					"(5304) [12-FEB-2019 17:25:10:804] nttrd: socket 844 had bytes read=8\r\n" +
					"(5304) [12-FEB-2019 17:25:10:804] nttrd: exit\r\n" +
					"(5304) [12-FEB-2019 17:25:10:804] nsprecv: 8 bytes from transport\r\n" +
					"(5304) [12-FEB-2019 17:25:10:804] nsprecv: tlen=8, plen=8, type=11\r\n" +
					"(5304) [12-FEB-2019 17:25:10:804] nsprecv: packet dump\r\n" +
					"(5304) [12-FEB-2019 17:25:10:804] nsprecv: 00 08 00 00 0B 00 00 00  |........|\r\n" +
					"(5304) [12-FEB-2019 17:25:10:804] nsprecv: normal exit\r\n" +
					"(5304) [12-FEB-2019 17:25:10:804] nscon: got NSPTRS packet\r\n" +
					"",
			},
			want: &Packet{
				Line:    13,
				Pid:     5304,
				TS:      []byte("12-FEB-2019 17:25:10:804"),
				Typ:     "nsprecv",
				Client:  "Service.exe",
				Payload: []byte{0x00, 0x08, 0x00, 0x00, 0x0B, 0x00, 0x00, 0x00},
				Socket:  844,
			},
		},

		{
			name: "nsbasic_bsd",
			args: args{
				trc: "(2100) [11-MAR-2019 07:04:24:529] nsbasic_bsd: entry\r\n" +
					"(2100) [11-MAR-2019 07:04:24:529] nsbasic_bsd: tot=0, plen=37.\r\n" +
					"(2100) [11-MAR-2019 07:04:24:529] nttfpwr: entry\r\n" +
					"(2100) [11-MAR-2019 07:04:24:529] nttfpwr: socket 1644 had bytes written=37\r\n" +
					"(2100) [11-MAR-2019 07:04:24:529] nttfpwr: exit\r\n" +
					"(2100) [11-MAR-2019 07:04:24:529] nsbasic_bsd: packet dump\r\n" +
					"(2100) [11-MAR-2019 07:04:24:529] nsbasic_bsd: 00 25 00 00 06 00 00 00  |.%......|\r\n" +
					"(2100) [11-MAR-2019 07:04:24:529] nsbasic_bsd: 00 00 01 06 05 04 03 02  |........|\r\n" +
					"(2100) [11-MAR-2019 07:04:24:529] nsbasic_bsd: 01 00 49 42 4D 50 43 2F  |..IBMPC/|\r\n" +
					"(2100) [11-MAR-2019 07:04:24:529] nsbasic_bsd: 57 49 4E 5F 4E 54 2D 38  |WIN_NT-8|\r\n" +
					"(2100) [11-MAR-2019 07:04:24:529] nsbasic_bsd: 2E 31 2E 30 00           |.1.0.   |\r\n" +
					"(2100) [11-MAR-2019 07:04:24:529] nsbasic_bsd: exit (0)\r\n" +
					"(2100) [11-MAR-2019 07:04:24:529] nsbasic_brc: entry: oln/tot=0\r\n" +

					"",
			},
			want: &Packet{
				Line:    7,
				Pid:     2100,
				TS:      []byte("11-MAR-2019 07:04:24:529"),
				Typ:     "nsbasic_bsd",
				Client:  "",
				Payload: []uint8{0x0, 0x25, 0x0, 0x0, 0x6, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x6, 0x5, 0x4, 0x3, 0x2, 0x1, 0x0, 0x49, 0x42, 0x4d, 0x50, 0x43, 0x2f, 0x57, 0x49, 0x4e, 0x5f, 0x4e, 0x54, 0x2d, 0x38, 0x2e, 0x31, 0x2e, 0x30, 0x0},
				Socket:  1644,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getPacketFromTraceSnippet(tt.args.trc)
			if (err != nil) != tt.wantErr {
				t.Errorf("getPacketFromTraceSnippet() error = %#v, wantErr %#v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getPacketFromTraceSnippet() = \n%#v,\n want \n%#v", got, tt.want)
			}
		})
	}
}
