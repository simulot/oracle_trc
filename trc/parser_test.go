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
			name: "simple",
			args: args{
				trc: "(5304) [12-FEB-2019 17:25:10:804] nsprecv: 8 bytes from transport\r\n" +
					"(5304) [12-FEB-2019 17:25:10:804] nsprecv: tlen=8, plen=8, type=11\r\n" +
					"(5304) [12-FEB-2019 17:25:10:804] nsprecv: packet dump\r\n" +
					"(5304) [12-FEB-2019 17:25:10:804] nsprecv: 00 08 00 00 0B 00 00 00  |........|\r\n" +
					"(5304) [12-FEB-2019 17:25:10:804] nsprecv: normal exit\r\n" +
					"(5304) [12-FEB-2019 17:25:10:804] nscon: got NSPTRS packet\r\n" +
					"",
			},
			want: &Packet{
				Line:    3,
				Pid:     5304,
				TS:      []byte("12-FEB-2019 17:25:10:804"),
				Typ:     "nsprecv",
				Payload: []byte{0x00, 0x08, 0x00, 0x00, 0x0B, 0x00, 0x00, 0x00},
			},
		},

		{
			name: "simple with connection info",
			args: args{
				trc: "(5304) [12-FEB-2019 17:25:10:647] nsc2addr: entry\r\n" +
					"(5304) [12-FEB-2019 17:25:10:647] nsc2addr: (DESCRIPTION=(CONNECT_DATA=(SID=BSWD01)(CID=(PROGRAM=C:\\App\\Service.exe)(HOST=HOST)(USER=SYSTEM)))(ADDRESS=(PROTOCOL=TCP)(HOST=10.30.194.77)(PORT=1525)))\r\n" +
					"(5304) [12-FEB-2019 17:25:10:647] nttbnd2addr: entry\r\n" +
					"(5304) [12-FEB-2019 17:25:10:647] nsc2addr: normal exit\r\n" +
					"(5304) [12-FEB-2019 17:25:10:804] nsprecv: 8 bytes from transport\r\n" +
					"(5304) [12-FEB-2019 17:25:10:804] nsprecv: tlen=8, plen=8, type=11\r\n" +
					"(5304) [12-FEB-2019 17:25:10:804] nsprecv: packet dump\r\n" +
					"(5304) [12-FEB-2019 17:25:10:804] nsprecv: 00 08 00 00 0B 00 00 00  |........|\r\n" +
					"(5304) [12-FEB-2019 17:25:10:804] nsprecv: normal exit\r\n" +
					"(5304) [12-FEB-2019 17:25:10:804] nscon: got NSPTRS packet\r\n" +
					"",
			},
			want: &Packet{
				Line:    7,
				Pid:     5304,
				TS:      []byte("12-FEB-2019 17:25:10:804"),
				Typ:     "nsprecv",
				Client:  "Service.exe",
				Payload: []byte{0x00, 0x08, 0x00, 0x00, 0x0B, 0x00, 0x00, 0x00},
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
