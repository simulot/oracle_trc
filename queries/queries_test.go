package queries

import (
	"errors"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

func Test_toUpperAscii(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "All caps",
			args: args{
				b: []byte("ALL CAPS"),
			},
			want: []byte("ALL CAPS"),
		},
		{
			name: "All small",
			args: args{
				b: []byte("all small"),
			},
			want: []byte("ALL SMALL"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toUpperAscii(tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toUpperAscii() = %v, want %v", got, tt.want)
			}
		})
	}
}

var withChars = regexp.MustCompile(`\p{Zs}|\r|\n|\t`)

// significantCharsEqual strips all separators to ease query comparison
func significantCharsEqual(s1, s2 string) bool {
	s1 = withChars.ReplaceAllLiteralString(s1, "")
	s2 = withChars.ReplaceAllLiteralString(s2, "")
	return s1 == s2
}

// getQueryFromTraceSnippet is an helper to get sql query from trace by using the parser type
func getQueryFromTraceSnippet(trc string) (*Query, error) {
	p := New(strings.NewReader(trc), "test")

	for {
		q, err := p.Next()
		if q == nil && err == nil {
			return nil, nil
		}
		return q, err
	}
	return nil, errors.New("Should not happen")
}
func Test_getQueryFromTraceSnippet(t *testing.T) {
	type args struct {
		trc string
	}
	type want struct {
		query  string
		params []string
	}
	tests := []struct {
		name    string
		args    args
		want    want
		wantErr bool
	}{
		{
			name: "simple short select",
			args: args{
				trc: `(5236) [22-OCT-2020 12:44:14:750] nsbasic_bsd: entry
(5236) [22-OCT-2020 12:44:14:750] nsbasic_bsd: tot=0, plen=164.
(5236) [22-OCT-2020 12:44:14:750] nttfpwr: entry
(5236) [22-OCT-2020 12:44:14:750] nttfpwr: socket 1288 had bytes written=164
(5236) [22-OCT-2020 12:44:14:750] nttfpwr: exit
(5236) [22-OCT-2020 12:44:14:750] nsbasic_bsd: packet dump
(5236) [22-OCT-2020 12:44:14:750] nsbasic_bsd: 00 A4 00 00 06 00 00 00  |........|
(5236) [22-OCT-2020 12:44:14:750] nsbasic_bsd: 00 00 11 69 15 01 01 01  |...i....|
(5236) [22-OCT-2020 12:44:14:750] nsbasic_bsd: 01 02 03 5E 16 02 80 69  |...^...i|
(5236) [22-OCT-2020 12:44:14:750] nsbasic_bsd: 00 01 01 AE 01 01 0D 01  |........|
(5236) [22-OCT-2020 12:44:14:750] nsbasic_bsd: 01 00 01 64 00 01 01 01  |...d....|
(5236) [22-OCT-2020 12:44:14:750] nsbasic_bsd: 00 01 00 01 01 01 00 00  |........|
(5236) [22-OCT-2020 12:44:14:750] nsbasic_bsd: 01 01 00 00 00 00 00 3A  |.......:|
(5236) [22-OCT-2020 12:44:14:750] nsbasic_bsd: 53 45 4C 45 43 54 20 70  |SELECT.p|
(5236) [22-OCT-2020 12:44:14:750] nsbasic_bsd: 61 72 61 6D 5F 76 61 6C  |aram_val|
(5236) [22-OCT-2020 12:44:14:750] nsbasic_bsd: 75 65 20 46 52 4F 4D 20  |ue.FROM.|
(5236) [22-OCT-2020 12:44:14:750] nsbasic_bsd: 65 66 6C 6F 77 5F 70 61  |eflow_pa|
(5236) [22-OCT-2020 12:44:14:750] nsbasic_bsd: 72 61 6D 73 20 57 48 45  |rams.WHE|
(5236) [22-OCT-2020 12:44:14:750] nsbasic_bsd: 52 45 20 70 61 72 61 6D  |RE.param|
(5236) [22-OCT-2020 12:44:14:750] nsbasic_bsd: 5F 6E 61 6D 65 20 3D 20  |_name.=.|
(5236) [22-OCT-2020 12:44:14:750] nsbasic_bsd: 3A 31 01 01 00 00 00 00  |:1......|
(5236) [22-OCT-2020 12:44:14:750] nsbasic_bsd: 00 00 01 01 00 00 00 00  |........|
(5236) [22-OCT-2020 12:44:14:750] nsbasic_bsd: 00 60 00 00 00 01 66 00  |.` + "`" + `....f.|
(5236) [22-OCT-2020 12:44:14:750] nsbasic_bsd: 01 10 00 00 02 03 69 01  |......i.|
(5236) [22-OCT-2020 12:44:14:750] nsbasic_bsd: 00 07 11 4C 49 43 45 58  |...LICEX|
(5236) [22-OCT-2020 12:44:14:750] nsbasic_bsd: 50 49 52 41 54 49 4F 4E  |PIRATION|
(5236) [22-OCT-2020 12:44:14:750] nsbasic_bsd: 44 41 54 45              |DATE    |
(5236) [22-OCT-2020 12:44:14:750] nsbasic_bsd: exit (0)
`,
			},
			want: want{
				"SELECT param_value FROM eflow_params WHERE param_name = :1",
				[]string{
					"LICEXPIRATIONDATE",
				},
			},
		},
		{
			name: "simple short select",
			args: args{
				trc: `(5236) [22-OCT-2020 12:44:09:974] nsbasic_bsd: entry
(5236) [22-OCT-2020 12:44:09:974] nsbasic_bsd: tot=0, plen=133.
(5236) [22-OCT-2020 12:44:09:974] nttfpwr: entry
(5236) [22-OCT-2020 12:44:09:974] nttfpwr: socket 1288 had bytes written=133
(5236) [22-OCT-2020 12:44:09:974] nttfpwr: exit
(5236) [22-OCT-2020 12:44:09:974] nsbasic_bsd: packet dump
(5236) [22-OCT-2020 12:44:09:974] nsbasic_bsd: 00 85 00 00 06 00 00 00  |........|
(5236) [22-OCT-2020 12:44:09:974] nsbasic_bsd: 00 00 11 69 09 01 01 01  |...i....|
(5236) [22-OCT-2020 12:44:09:974] nsbasic_bsd: 01 01 03 5E 0A 02 80 61  |...^...a|
(5236) [22-OCT-2020 12:44:09:974] nsbasic_bsd: 00 01 01 BD 01 01 0D 01  |........|
(5236) [22-OCT-2020 12:44:09:974] nsbasic_bsd: 01 00 01 64 00 00 00 00  |...d....|
(5236) [22-OCT-2020 12:44:09:974] nsbasic_bsd: 01 00 01 01 01 00 00 01  |........|
(5236) [22-OCT-2020 12:44:09:974] nsbasic_bsd: 01 00 00 00 00 00 3F 53  |......?S|
(5236) [22-OCT-2020 12:44:09:974] nsbasic_bsd: 45 4C 45 43 54 20 44 49  |ELECT.DI|
(5236) [22-OCT-2020 12:44:09:974] nsbasic_bsd: 53 54 49 4E 43 54 20 27  |STINCT.'|
(5236) [22-OCT-2020 12:44:09:974] nsbasic_bsd: 65 66 6C 6F 77 5F 70 61  |eflow_pa|
(5236) [22-OCT-2020 12:44:09:974] nsbasic_bsd: 72 61 6D 73 27 2C 65 66  |rams',ef|
(5236) [22-OCT-2020 12:44:09:974] nsbasic_bsd: 6C 6F 77 5F 70 61 72 61  |low_para|
(5236) [22-OCT-2020 12:44:09:974] nsbasic_bsd: 6D 73 2E 2A 20 46 52 4F  |ms.*.FRO|
(5236) [22-OCT-2020 12:44:09:974] nsbasic_bsd: 4D 20 65 66 6C 6F 77 5F  |M.eflow_|
(5236) [22-OCT-2020 12:44:09:974] nsbasic_bsd: 70 61 72 61 6D 73 01 01  |params..|
(5236) [22-OCT-2020 12:44:09:974] nsbasic_bsd: 00 00 00 00 00 00 01 01  |........|
(5236) [22-OCT-2020 12:44:09:974] nsbasic_bsd: 00 00 00 00 00           |.....   |
(5236) [22-OCT-2020 12:44:09:974] nsbasic_bsd: exit (0)
(5236) [22-OCT-2020 12:44:09:974] nsbasic_brc: entry: oln/tot=0`,
			},
			want: want{
				"SELECT DISTINCT 'eflow_params',eflow_params.* FROM eflow_params",
				nil,
			},
		},
		{
			name: "simple short select 2",
			args: args{
				trc: `(5236) [22-OCT-2020 12:44:14:753] nsbasic_bsd: entry
(5236) [22-OCT-2020 12:44:14:753] nsbasic_bsd: tot=0, plen=94.
(5236) [22-OCT-2020 12:44:14:753] nttfpwr: entry
(5236) [22-OCT-2020 12:44:14:753] nttfpwr: socket 1288 had bytes written=94
(5236) [22-OCT-2020 12:44:14:753] nttfpwr: exit
(5236) [22-OCT-2020 12:44:14:753] nsbasic_bsd: packet dump
(5236) [22-OCT-2020 12:44:14:753] nsbasic_bsd: 00 5E 00 00 06 00 00 00  |.^......|
(5236) [22-OCT-2020 12:44:14:753] nsbasic_bsd: 00 00 11 69 17 01 01 01  |...i....|
(5236) [22-OCT-2020 12:44:14:753] nsbasic_bsd: 01 01 03 5E 18 02 80 61  |...^...a|
(5236) [22-OCT-2020 12:44:14:753] nsbasic_bsd: 00 01 01 48 01 01 0D 01  |...H....|
(5236) [22-OCT-2020 12:44:14:753] nsbasic_bsd: 01 00 01 64 00 00 00 00  |...d....|
(5236) [22-OCT-2020 12:44:14:753] nsbasic_bsd: 01 00 01 01 01 00 00 01  |........|
(5236) [22-OCT-2020 12:44:14:753] nsbasic_bsd: 01 00 00 00 00 00 18 53  |.......S|
(5236) [22-OCT-2020 12:44:14:753] nsbasic_bsd: 45 4C 45 43 54 20 53 59  |ELECT.SY|
(5236) [22-OCT-2020 12:44:14:753] nsbasic_bsd: 53 44 41 54 45 20 46 52  |SDATE.FR|
(5236) [22-OCT-2020 12:44:14:753] nsbasic_bsd: 4F 4D 20 44 55 41 4C 01  |OM.DUAL.|
(5236) [22-OCT-2020 12:44:14:753] nsbasic_bsd: 01 00 00 00 00 00 00 01  |........|
(5236) [22-OCT-2020 12:44:14:753] nsbasic_bsd: 01 00 00 00 00 00        |......  |
(5236) [22-OCT-2020 12:44:14:753] nsbasic_bsd: exit (0)
		`,
			},
			want: want{
				"SELECT SYSDATE FROM DUAL",
				nil,
			},
		},
		{
			name: "long query with 3 parameters",
			args: args{
				trc: `(5236) [22-OCT-2020 12:44:16:136] nsbasic_bsd: entry
(5236) [22-OCT-2020 12:44:16:136] nsbasic_bsd: tot=0, plen=434.
(5236) [22-OCT-2020 12:44:16:136] nttfpwr: entry
(5236) [22-OCT-2020 12:44:16:136] nttfpwr: socket 1288 had bytes written=434
(5236) [22-OCT-2020 12:44:16:136] nttfpwr: exit
(5236) [22-OCT-2020 12:44:16:136] nsbasic_bsd: packet dump
(5236) [22-OCT-2020 12:44:16:136] nsbasic_bsd: 01 B2 00 00 06 00 00 00  |........|
(5236) [22-OCT-2020 12:44:16:136] nsbasic_bsd: 00 00 11 69 59 01 01 01  |...iY...|
(5236) [22-OCT-2020 12:44:16:136] nsbasic_bsd: 01 02 03 5E 5A 02 80 69  |...^Z..i|
(5236) [22-OCT-2020 12:44:16:136] nsbasic_bsd: 00 01 02 03 21 01 01 0D  |....!...|
(5236) [22-OCT-2020 12:44:16:136] nsbasic_bsd: 01 01 00 01 64 00 01 01  |....d...|
(5236) [22-OCT-2020 12:44:16:136] nsbasic_bsd: 03 00 01 00 01 01 01 00  |........|
(5236) [22-OCT-2020 12:44:16:136] nsbasic_bsd: 00 01 01 00 00 00 00 00  |........|
(5236) [22-OCT-2020 12:44:16:136] nsbasic_bsd: FE 40 53 45 4C 45 43 54  |.@SELECT|
(5236) [22-OCT-2020 12:44:16:136] nsbasic_bsd: 20 74 31 2E 70 73 65 5F  |.t1.pse_|
(5236) [22-OCT-2020 12:44:16:136] nsbasic_bsd: 64 61 74 61 2C 20 74 32  |data,.t2|
(5236) [22-OCT-2020 12:44:16:136] nsbasic_bsd: 2E 70 73 64 5F 61 70 70  |.psd_app|
(5236) [22-OCT-2020 12:44:16:136] nsbasic_bsd: 6C 69 63 61 74 69 6F 6E  |lication|
(5236) [22-OCT-2020 12:44:16:136] nsbasic_bsd: 2C 20 74 32 2E 70 73 64  |,.t2.psd|
(5236) [22-OCT-2020 12:44:16:136] nsbasic_bsd: 5F 61 64 6D 69 6E 5F 76  |_admin_v|
(5236) [22-OCT-2020 12:44:16:136] nsbasic_bsd: 69 73 69 62 6C 65 2C 20  |isible,.|
(5236) [22-OCT-2020 12:44:16:136] nsbasic_bsd: 74 32 40 2E 70 73 64 5F  |t2@.psd_|
(5236) [22-OCT-2020 12:44:16:136] nsbasic_bsd: 63 6F 6D 6D 65 6E 74 2C  |comment,|
(5236) [22-OCT-2020 12:44:16:136] nsbasic_bsd: 20 74 32 2E 70 73 64 5F  |.t2.psd_|
(5236) [22-OCT-2020 12:44:16:136] nsbasic_bsd: 63 6F 6D 6D 65 6E 74 20  |comment.|
(5236) [22-OCT-2020 12:44:16:136] nsbasic_bsd: 46 52 4F 4D 20 69 70 5F  |FROM.ip_|
(5236) [22-OCT-2020 12:44:16:136] nsbasic_bsd: 70 65 72 73 6F 6E 61 6C  |personal|
(5236) [22-OCT-2020 12:44:16:136] nsbasic_bsd: 5F 73 65 74 74 69 6E 67  |_setting|
(5236) [22-OCT-2020 12:44:16:136] nsbasic_bsd: 73 20 74 31 2C 20 69 70  |s.t1,.ip|
(5236) [22-OCT-2020 12:44:16:136] nsbasic_bsd: 5F 70 65 40 72 73 6F 6E  |_pe@rson|
(5236) [22-OCT-2020 12:44:16:136] nsbasic_bsd: 61 6C 5F 73 65 74 74 69  |al_setti|
(5236) [22-OCT-2020 12:44:16:136] nsbasic_bsd: 6E 67 73 5F 64 65 66 20  |ngs_def.|
(5236) [22-OCT-2020 12:44:16:136] nsbasic_bsd: 74 32 20 57 48 45 52 45  |t2.WHERE|
(5236) [22-OCT-2020 12:44:16:137] nsbasic_bsd: 20 74 31 2E 70 73 64 5F  |.t1.psd_|
(5236) [22-OCT-2020 12:44:16:137] nsbasic_bsd: 64 61 74 61 5F 63 6F 64  |data_cod|
(5236) [22-OCT-2020 12:44:16:137] nsbasic_bsd: 65 20 3D 20 3A 31 20 41  |e.=.:1.A|
(5236) [22-OCT-2020 12:44:16:137] nsbasic_bsd: 4E 44 20 74 31 2E 70 73  |ND.t1.ps|
(5236) [22-OCT-2020 12:44:16:137] nsbasic_bsd: 64 5F 64 61 40 74 61 5F  |d_da@ta_|
(5236) [22-OCT-2020 12:44:16:137] nsbasic_bsd: 63 6F 64 65 20 3D 20 74  |code.=.t|
(5236) [22-OCT-2020 12:44:16:137] nsbasic_bsd: 32 2E 70 73 64 5F 64 61  |2.psd_da|
(5236) [22-OCT-2020 12:44:16:137] nsbasic_bsd: 74 61 5F 63 6F 64 65 20  |ta_code.|
(5236) [22-OCT-2020 12:44:16:137] nsbasic_bsd: 41 4E 44 20 74 31 2E 75  |AND.t1.u|
(5236) [22-OCT-2020 12:44:16:137] nsbasic_bsd: 73 65 72 5F 6E 65 74 77  |ser_netw|
(5236) [22-OCT-2020 12:44:16:137] nsbasic_bsd: 6F 72 6B 5F 6E 61 6D 65  |ork_name|
(5236) [22-OCT-2020 12:44:16:137] nsbasic_bsd: 20 3D 20 3A 32 20 41 4E  |.=.:2.AN|
(5236) [22-OCT-2020 12:44:16:137] nsbasic_bsd: 44 20 74 31 2E 0B 64 6F  |D.t1..do|
(5236) [22-OCT-2020 12:44:16:137] nsbasic_bsd: 6D 61 69 6E 20 3D 20 3A  |main.=.:|
(5236) [22-OCT-2020 12:44:16:137] nsbasic_bsd: 33 00 01 01 00 00 00 00  |3.......|
(5236) [22-OCT-2020 12:44:16:137] nsbasic_bsd: 00 00 01 01 00 00 00 00  |........|
(5236) [22-OCT-2020 12:44:16:137] nsbasic_bsd: 00 60 00 00 00 01 A2 00  |.` + "`" + `......|
(5236) [22-OCT-2020 12:44:16:137] nsbasic_bsd: 01 10 00 00 02 03 69 01  |......i.|
(5236) [22-OCT-2020 12:44:16:137] nsbasic_bsd: 00 60 00 00 00 01 36 00  |.` + "`" + `....6.|
(5236) [22-OCT-2020 12:44:16:137] nsbasic_bsd: 01 10 00 00 02 03 69 01  |......i.|
(5236) [22-OCT-2020 12:44:16:137] nsbasic_bsd: 00 60 00 00 00 01 06 00  |.` + "`" + `......|
(5236) [22-OCT-2020 12:44:16:137] nsbasic_bsd: 01 10 00 00 02 03 69 01  |......i.|
(5236) [22-OCT-2020 12:44:16:137] nsbasic_bsd: 00 07 1B 49 50 5F 4D 41  |...IP_MA|
(5236) [22-OCT-2020 12:44:16:137] nsbasic_bsd: 53 54 45 52 5F 53 45 41  |STER_SEA|
(5236) [22-OCT-2020 12:44:16:137] nsbasic_bsd: 52 43 48 5F 43 52 49 54  |RCH_CRIT|
(5236) [22-OCT-2020 12:44:16:137] nsbasic_bsd: 45 52 49 41 5F 32 09 4A  |ERIA_2.J|
(5236) [22-OCT-2020 12:44:16:137] nsbasic_bsd: 46 2E 43 41 53 53 41 4E  |F.CASSAN|
(5236) [22-OCT-2020 12:44:16:137] nsbasic_bsd: 01 20                    |..      |
(5236) [22-OCT-2020 12:44:16:137] nsbasic_bsd: exit (0)
(5236) [22-OCT-2020 12:44:16:137] nsbasic_brc: entry: oln/tot=0
(5236) [22-OCT-2020 12:44:16:137] nttfprd: entry
`,
			},
			want: want{
				"SELECT t1.pse_data, t2.psd_application, t2.psd_admin_visible, t2.psd_comment, t2.psd_comment FROM ip_personal_settings t1, ip_personal_settings_def t2 WHERE t1.psd_data_code = :1 AND t1.psd_data_code = t2.psd_data_code AND t1.user_network_name = :2 AND t1.domain = :3",
				[]string{
					"IP_MASTER_SEARCH_CRITERIA_2",
					"JF.CASSAN",
					" ",
				},
			},
		},
		{
			name: "simple long select",
			args: args{
				trc: `(5236) [22-OCT-2020 12:44:14:828] nsbasic_bsd: entry
(5236) [22-OCT-2020 12:44:14:828] nsbasic_bsd: tot=0, plen=1768.
(5236) [22-OCT-2020 12:44:14:828] nttfpwr: entry
(5236) [22-OCT-2020 12:44:14:828] nttfpwr: socket 1288 had bytes written=1768
(5236) [22-OCT-2020 12:44:14:828] nttfpwr: exit
(5236) [22-OCT-2020 12:44:14:828] nsbasic_bsd: packet dump
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 06 E8 00 00 06 00 00 00  |........|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 00 00 11 69 31 01 01 01  |...i1...|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 01 02 03 5E 32 02 80 61  |...^2..a|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 00 01 02 13 8F 01 01 0D  |........|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 01 01 00 01 64 00 00 00  |....d...|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 00 01 00 01 01 01 00 00  |........|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 01 01 00 00 00 00 00 FE  |........|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 40 53 45 4C 45 43 54 20  |@SELECT.|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 44 49 53 54 49 4E 43 54  |DISTINCT|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 20 74 32 2E 67 72 6F 75  |.t2.grou|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 70 5F 6E 61 6D 65 2C 20  |p_name,.|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 74 32 2E 67 72 6F 75 70  |t2.group|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 5F 74 79 70 65 2C 20 74  |_type,.t|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 32 2E 67 72 6F 75 70 5F  |2.group_|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 72 6F 6C 65 2C 20 74 32  |role,.t2|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 2E 40 67 72 6F 75 70 5F  |.@group_|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 73 65 63 75 72 69 74 79  |security|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 5F 6C 65 76 65 6C 2C 20  |_level,.|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 74 34 2E 63 6F 6D 70 5F  |t4.comp_|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 6E 6F 2C 20 74 34 2E 63  |no,.t4.c|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 6F 6D 70 5F 6E 61 6D 65  |omp_name|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 2C 20 74 32 2E 64 65 73  |,.t2.des|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 63 72 69 70 74 69 6F 6E  |cription|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 2C 20 40 74 33 2E 72 69  |,.@t3.ri|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 67 5F 73 63 61 6E 2C 20  |g_scan,.|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 74 33 2E 72 69 67 5F 61  |t3.rig_a|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 63 63 72 75 61 6C 72 65  |ccrualre|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 70 6F 72 74 2C 20 74 33  |port,.t3|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 2E 72 69 67 5F 63 6F 6E  |.rig_con|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 74 72 61 63 74 61 70 70  |tractapp|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 72 6F 76 61 6C 2C 20 74  |roval,.t|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 33 2E 72 40 69 67 5F 63  |3.r@ig_c|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 6F 6E 74 72 61 63 74 61  |ontracta|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 70 70 72 6F 76 61 6C 5F  |pproval_|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 6C 69 6D 69 74 2C 20 74  |limit,.t|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 33 2E 72 69 67 5F 63 6C  |3.rig_cl|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 69 65 6E 74 2C 20 74 33  |ient,.t3|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 2E 72 69 67 5F 69 6E 76  |.rig_inv|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 61 64 6D 69 6E 2C 20 74  |admin,.t|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 33 2E 72 69 40 67 5F 72  |3.ri@g_r|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 65 70 6F 72 74 2C 20 74  |eport,.t|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 33 2E 72 69 67 5F 61 70  |3.rig_ap|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 70 72 6F 76 65 72 2C 20  |prover,.|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 74 33 2E 72 69 67 5F 61  |t3.rig_a|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 70 70 72 6F 76 65 5F 6C  |pprove_l|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 69 6D 69 74 2C 20 74 33  |imit,.t3|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 2E 72 69 67 5F 6E 75 6C  |.rig_nul|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 6C 69 66 69 65 40 72 2C  |lifie@r,|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 20 74 33 2E 72 69 67 5F  |.t3.rig_|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 61 64 6D 69 6E 61 70 70  |adminapp|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 72 6F 76 65 2C 20 74 33  |rove,.t3|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 2E 72 69 67 5F 6D 61 74  |.rig_mat|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 63 68 69 6E 67 2C 20 74  |ching,.t|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 33 2E 72 69 67 5F 63 6D  |3.rig_cm|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 76 69 65 77 2C 20 74 33  |view,.t3|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 2E 72 69 67 5F 63 40 6D  |.rig_c@m|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 63 72 65 61 74 65 2C 20  |create,.|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 74 33 2E 72 69 67 5F 63  |t3.rig_c|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 6D 64 65 6C 65 74 65 2C  |mdelete,|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 20 74 33 2E 72 69 67 5F  |.t3.rig_|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 63 6D 61 70 70 72 6F 76  |cmapprov|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 65 2C 20 74 33 2E 72 69  |e,.t3.ri|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 67 5F 63 6D 61 70 70 72  |g_cmappr|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 6F 76 65 6C 69 6D 69 40  |ovelimi@|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 74 2C 20 74 33 2E 72 69  |t,.t3.ri|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 67 5F 63 6D 77 6F 72 6B  |g_cmwork|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 66 6C 6F 77 6D 61 6E 61  |flowmana|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 67 65 72 2C 20 74 33 2E  |ger,.t3.|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 72 69 67 5F 63 6D 72 65  |rig_cmre|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 70 6F 72 74 69 6E 67 2C  |porting,|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 20 74 33 2E 72 69 67 5F  |.t3.rig_|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 69 6E 76 6F 69 63 65 5F  |invoice_|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 40 61 64 6D 69 6E 5F 61  |@admin_a|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 70 70 72 6F 76 65 72 2C  |pprover,|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 20 74 33 2E 72 69 67 5F  |.t3.rig_|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 69 6E 76 6F 69 63 65 5F  |invoice_|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 61 64 6D 69 6E 5F 6C 69  |admin_li|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 6D 69 74 2C 20 74 33 2E  |mit,.t3.|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 72 69 67 5F 63 6D 63 6F  |rig_cmco|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 6E 74 72 61 63 74 5F 61  |ntract_a|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 70 40 70 72 6F 76 65 5F  |p@prove_|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 6D 69 6C 65 73 74 6F 2C  |milesto,|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 20 74 33 2E 72 69 67 5F  |.t3.rig_|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 63 6D 63 6F 6E 74 72 61  |cmcontra|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 63 74 5F 64 65 6C 65 74  |ct_delet|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 65 5F 6D 69 6C 65 73 74  |e_milest|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 6F 2C 20 74 33 2E 72 69  |o,.t3.ri|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 67 5F 63 6D 63 6F 6E 74  |g_cmcont|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 72 61 40 63 74 5F 64 65  |ra@ct_de|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 6C 5F 6D 69 6C 65 73 74  |l_milest|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 6F 5F 64 6F 63 2C 20 74  |o_doc,.t|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 33 2E 72 69 67 5F 63 6D  |3.rig_cm|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 63 6F 6E 74 72 61 63 74  |contract|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 5F 64 65 6C 65 74 65 5F  |_delete_|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 64 6F 63 75 6D 65 6E 74  |document|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 2C 20 74 33 2E 72 69 67  |,.t3.rig|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 5F 63 6D 40 63 6F 6E 74  |_cm@cont|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 72 61 63 74 5F 65 64 69  |ract_edi|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 74 5F 64 6F 63 75 6D 65  |t_docume|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 6E 74 2C 20 74 33 2E 72  |nt,.t3.r|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 69 67 5F 63 6D 63 6F 6E  |ig_cmcon|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 74 72 61 63 74 5F 72 61  |tract_ra|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 74 69 66 79 2C 20 74 33  |tify,.t3|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 2E 72 69 67 5F 63 6D 66  |.rig_cmf|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 69 6C 65 5F 40 76 69 65  |ile_@vie|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 77 2C 20 74 33 2E 72 69  |w,.t3.ri|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 67 5F 63 6D 66 69 6C 65  |g_cmfile|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 5F 63 72 65 61 74 65 2C  |_create,|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 20 74 33 2E 72 69 67 5F  |.t3.rig_|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 63 6D 66 69 6C 65 5F 64  |cmfile_d|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 65 6C 65 74 65 2C 20 74  |elete,.t|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 33 2E 72 69 67 5F 63 6D  |3.rig_cm|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 66 69 6C 65 5F 40 72 65  |file_@re|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 70 6F 72 74 2C 20 74 33  |port,.t3|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 2E 72 69 67 5F 63 6D 66  |.rig_cmf|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 69 6C 65 5F 64 65 6C 65  |ile_dele|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 74 65 5F 64 6F 63 75 6D  |te_docum|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 65 6E 74 2C 20 74 33 2E  |ent,.t3.|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 72 69 67 5F 63 6D 66 69  |rig_cmfi|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 6C 65 5F 65 64 69 74 5F  |le_edit_|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 64 6F 63 75 6D 65 40 6E  |docume@n|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 74 2C 20 74 33 2E 72 69  |t,.t3.ri|
(5236) [22-OCT-2020 12:44:14:829] nsbasic_bsd: 67 5F 63 6D 66 69 6C 65  |g_cmfile|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 5F 62 6C 61 6E 6B 5F 31  |_blank_1|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 2C 20 74 33 2E 72 69 67  |,.t3.rig|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 5F 63 6D 66 69 6C 65 5F  |_cmfile_|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 62 6C 61 6E 6B 5F 32 2C  |blank_2,|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 20 74 33 2E 72 69 67 5F  |.t3.rig_|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 63 6D 66 69 6C 65 5F 40  |cmfile_@|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 62 6C 61 6E 6B 5F 33 2C  |blank_3,|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 20 74 33 2E 72 69 67 5F  |.t3.rig_|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 63 6D 66 69 6C 65 5F 62  |cmfile_b|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 6C 61 6E 6B 5F 34 2C 20  |lank_4,.|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 74 33 2E 72 69 67 5F 63  |t3.rig_c|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 6D 6F 74 68 65 72 5F 73  |mother_s|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 75 70 70 6C 69 65 72 5F  |upplier_|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 6D 61 6E 61 67 65 6D 65  |manageme|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 40 2C 20 74 33 2E 72 69  |@,.t3.ri|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 67 5F 63 6D 6F 74 68 65  |g_cmothe|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 72 5F 73 79 73 74 65 6D  |r_system|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 5F 6D 61 6E 61 67 65 6D  |_managem|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 65 6E 74 2C 20 74 33 2E  |ent,.t3.|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 72 69 67 5F 63 6D 6F 74  |rig_cmot|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 68 65 72 5F 63 72 65 61  |her_crea|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 74 65 5F 69 6E 76 6F 69  |te_invoi|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 63 40 65 2C 20 74 33 2E  |c@e,.t3.|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 72 69 67 5F 63 6D 6F 74  |rig_cmot|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 68 65 72 5F 63 72 65 61  |her_crea|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 74 65 5F 70 6F 2C 20 74  |te_po,.t|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 33 2E 72 69 67 5F 63 6D  |3.rig_cm|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 6F 74 68 65 72 5F 63 72  |other_cr|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 65 61 74 65 5F 67 72 2C  |eate_gr,|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 20 74 33 2E 72 69 67 5F  |.t3.rig_|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 63 6D 40 6F 74 68 65 72  |cm@other|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 5F 62 6C 61 6E 6B 5F 31  |_blank_1|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 2C 20 74 33 2E 72 69 67  |,.t3.rig|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 5F 63 6D 6F 74 68 65 72  |_cmother|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 5F 62 6C 61 6E 6B 5F 32  |_blank_2|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 20 46 52 4F 4D 20 69 70  |.FROM.ip|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 5F 63 6F 6D 70 61 6E 79  |_company|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 5F 75 73 65 72 5F 67 72  |_user_gr|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 6F 75 70 40 20 74 31 2C  |oup@.t1,|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 20 65 66 6C 6F 77 5F 67  |.eflow_g|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 72 6F 75 70 73 20 74 32  |roups.t2|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 2C 20 69 70 5F 72 69 67  |,.ip_rig|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 68 74 73 20 74 33 2C 20  |hts.t3,.|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 63 6F 6D 70 61 6E 69 65  |companie|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 73 20 74 34 20 57 48 45  |s.t4.WHE|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 52 45 20 74 32 2E 67 72  |RE.t2.gr|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 6F 75 70 5F 40 6E 61 6D  |oup_@nam|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 65 20 3D 20 74 31 2E 67  |e.=.t1.g|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 72 6F 75 70 5F 6E 61 6D  |roup_nam|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 65 20 41 4E 44 20 74 33  |e.AND.t3|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 2E 67 72 6F 75 70 5F 6E  |.group_n|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 61 6D 65 20 3D 20 74 32  |ame.=.t2|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 2E 67 72 6F 75 70 5F 6E  |.group_n|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 61 6D 65 20 41 4E 44 20  |ame.AND.|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 28 74 31 2E 63 40 6F 6D  |(t1.c@om|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 70 5F 6E 6F 20 49 4E 28  |p_no.IN(|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 53 45 4C 45 43 54 20 44  |SELECT.D|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 49 53 54 49 4E 43 54 20  |ISTINCT.|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 74 31 2E 63 6F 6D 70 5F  |t1.comp_|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 6E 6F 20 46 52 4F 4D 20  |no.FROM.|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 69 70 5F 63 6F 6D 70 61  |ip_compa|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 6E 79 5F 75 73 65 72 5F  |ny_user_|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 67 72 6F 75 70 20 40 74  |group.@t|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 31 2C 20 69 70 5F 67 72  |1,.ip_gr|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 6F 75 70 5F 75 73 65 72  |oup_user|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 20 74 32 20 2C 69 70 5F  |.t2.,ip_|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 72 69 67 68 74 73 20 74  |rights.t|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 34 20 57 48 45 52 45 20  |4.WHERE.|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 74 32 2E 75 73 65 72 5F  |t2.user_|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 6E 65 74 77 6F 72 6B 5F  |network_|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 6E 61 6D 65 20 3D 20 40  |name.=.@|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 27 4A 46 2E 43 41 53 53  |'JF.CASS|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 41 4E 27 20 41 4E 44 20  |AN'.AND.|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 74 32 2E 64 6F 6D 61 69  |t2.domai|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 6E 3D 20 27 20 27 20 41  |n=.'.'.A|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 4E 44 20 20 74 31 2E 67  |ND..t1.g|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 72 6F 75 70 5F 6E 61 6D  |roup_nam|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 65 20 3D 20 74 32 2E 67  |e.=.t2.g|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 72 6F 75 70 5F 6E 61 6D  |roup_nam|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 40 65 20 20 20 41 4E 44  |@e...AND|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 20 74 31 2E 67 72 6F 75  |.t1.grou|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 70 5F 6E 61 6D 65 20 3D  |p_name.=|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 20 74 34 2E 67 72 6F 75  |.t4.grou|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 70 5F 6E 61 6D 65 20 20  |p_name..|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 20 41 4E 44 20 28 74 34  |.AND.(t4|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 2E 72 69 67 5F 69 6E 76  |.rig_inv|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 61 64 6D 69 6E 20 3D 20  |admin.=.|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 31 40 20 29 20 20 29 29  |1@.)..))|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 20 41 4E 44 20 74 34 2E  |.AND.t4.|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 63 6F 6D 70 5F 6E 6F 20  |comp_no.|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 3D 20 74 31 2E 63 6F 6D  |=.t1.com|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 70 5F 6E 6F 20 4F 52 44  |p_no.ORD|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 45 52 20 42 59 20 74 32  |ER.BY.t2|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 2E 67 72 6F 75 70 5F 6E  |.group_n|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 61 6D 65 2C 20 74 34 2E  |ame,.t4.|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 63 6F 05 6D 70 5F 6E 6F  |co.mp_no|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 00 01 01 00 00 00 00 00  |........|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: 00 01 01 00 00 00 00 00  |........|
(5236) [22-OCT-2020 12:44:14:830] nsbasic_bsd: exit (0)
(5236) [22-OCT-2020 12:44:14:830] nsbasic_brc: entry: oln/tot=0
(5236) [22-OCT-2020 12:44:14:830] nttfprd: entry
(5236) [22-OCT-2020 12:44:14:836] nttfprd: socket 1288 had bytes read=8155
(5236) [22-OCT-2020 12:44:14:836] nttfprd: exit
`,
			},
			want: want{
				"SELECT DISTINCT t2.group_name, t2.group_type, t2.group_role, t2.group_security_level, t4.comp_no, t4.comp_name, t2.description, t3.rig_scan, t3.rig_accrualreport, t3.rig_contractapproval, t3.rig_contractapproval_limit, t3.rig_client, t3.rig_invadmin, t3.rig_report, t3.rig_approver, t3.rig_approve_limit, t3.rig_nullifier, t3.rig_adminapprove, t3.rig_matching, t3.rig_cmview, t3.rig_cmcreate, t3.rig_cmdelete, t3.rig_cmapprove, t3.rig_cmapprovelimit, t3.rig_cmworkflowmanager, t3.rig_cmreporting, t3.rig_invoice_admin_approver, t3.rig_invoice_admin_limit, t3.rig_cmcontract_approve_milesto, t3.rig_cmcontract_delete_milesto, t3.rig_cmcontract_del_milesto_doc, t3.rig_cmcontract_delete_document, t3.rig_cmcontract_edit_document, t3.rig_cmcontract_ratify, t3.rig_cmfile_view, t3.rig_cmfile_create, t3.rig_cmfile_delete, t3.rig_cmfile_report, t3.rig_cmfile_delete_document, t3.rig_cmfile_edit_document, t3.rig_cmfile_blank_1, t3.rig_cmfile_blank_2, t3.rig_cmfile_blank_3, t3.rig_cmfile_blank_4, t3.rig_cmother_supplier_manageme, t3.rig_cmother_system_management, t3.rig_cmother_create_invoice, t3.rig_cmother_create_po, t3.rig_cmother_create_gr, t3.rig_cmother_blank_1, t3.rig_cmother_blank_2 FROM ip_company_user_group t1, eflow_groups t2, ip_rights t3, companies t4 WHERE t2.group_name = t1.group_name AND t3.group_name = t2.group_name AND (t1.comp_no IN(SELECT DISTINCT t1.comp_no FROM ip_company_user_group t1, ip_group_user t2 ,ip_rights t4 WHERE t2.user_network_name = 'JF.CASSAN' AND t2.domain= ' ' AND  t1.group_name = t2.group_name   AND t1.group_name = t4.group_name   AND (t4.rig_invadmin = 1 )  )) AND t4.comp_no = t1.comp_no ORDER BY t2.group_name, t4.comp_no",
				nil,
			},
		},
		{
			name: "Long query, with CRLF and parameters",
			args: args{
				trc: `(5236) [22-OCT-2020 12:44:52:419] nioqrc: entry
(5236) [22-OCT-2020 12:44:52:419] nsbasic_bsd: entry
(5236) [22-OCT-2020 12:44:52:419] nsbasic_bsd: tot=0, plen=964.
(5236) [22-OCT-2020 12:44:52:419] nttfpwr: entry
(5236) [22-OCT-2020 12:44:52:419] nttfpwr: socket 1288 had bytes written=964
(5236) [22-OCT-2020 12:44:52:419] nttfpwr: exit
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: packet dump
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 03 C4 00 00 06 00 00 00  |........|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 00 00 11 69 3D 01 01 01  |...i=...|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 01 01 03 5E 3E 02 80 69  |...^>..i|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 00 01 02 08 A6 01 01 0D  |........|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 01 01 00 01 64 00 01 01  |....d...|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 05 00 01 00 01 01 01 00  |........|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 00 01 01 00 00 00 00 00  |........|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: FE 40 73 65 6C 65 63 74  |.@select|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 20 27 4E 4F 20 44 55 50  |.'NO.DUP|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 4C 49 43 41 54 45 53 27  |LICATES'|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 20 66 72 6F 6D 20 44 55  |.from.DU|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 41 4C 20 77 68 65 72 65  |AL.where|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 20 30 20 3D 20 6E 76 6C  |.0.=.nvl|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 28 28 20 0A 73 65 6C 65  |((..sele|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 63 74 20 0A 20 20 73 75  |ct....su|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 6D 28 40 20 63 61 73 65  |m(@.case|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 20 77 68 65 6E 20 46 4C  |.when.FL|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 2E 41 43 54 49 4F 4E 5F  |.ACTION_|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 49 4E 44 45 58 20 3D 20  |INDEX.=.|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 33 20 74 68 65 6E 20 30  |3.then.0|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 20 65 6C 73 65 20 31 20  |.else.1.|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 65 6E 64 29 20 61 73 20  |end).as.|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 43 20 0A 66 72 6F 6D 20  |C..from.|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 0A 20 20 40 28 73 65 6C  |...@(sel|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 65 63 74 20 3A 31 20 44  |ect.:1.D|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 4F 43 5F 49 44 2C 20 3A  |OC_ID,.:|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 32 20 53 55 50 50 4C 49  |2.SUPPLI|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 45 52 5F 4E 55 4D 2C 20  |ER_NUM,.|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 3A 33 20 49 4E 56 4F 49  |:3.INVOI|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 43 45 5F 4E 55 4D 2C 20  |CE_NUM,.|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 3A 34 20 49 4E 56 4F 49  |:4.INVOI|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 43 45 5F 44 40 41 54 45  |CE_D@ATE|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 2C 20 3A 35 20 43 4F 4D  |,.:5.COM|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 50 5F 4E 4F 20 66 72 6F  |P_NO.fro|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 6D 20 44 55 41 4C 20 29  |m.DUAL.)|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 52 45 46 20 6A 6F 69 6E  |REF.join|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 20 20 0A 20 20 44 4F 43  |.....DOC|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 53 20 44 20 6F 6E 20 52  |S.D.on.R|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 45 46 2E 44 4F 43 5F 49  |EF.DOC_I|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 44 20 3C 3E 20 40 44 2E  |D.<>.@D.|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 44 4F 43 5F 49 44 20 61  |DOC_ID.a|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 6E 64 20 52 45 46 2E 49  |nd.REF.I|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 4E 56 4F 49 43 45 5F 4E  |NVOICE_N|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 55 4D 20 3D 20 44 2E 49  |UM.=.D.I|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 4E 56 4F 49 43 45 5F 4E  |NVOICE_N|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 55 4D 20 61 6E 64 20 52  |UM.and.R|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 45 46 2E 43 4F 4D 50 5F  |EF.COMP_|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 4E 4F 20 3D 20 44 40 2E  |NO.=.D@.|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 43 4F 4D 50 5F 4E 4F 20  |COMP_NO.|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 61 6E 64 20 52 45 46 2E  |and.REF.|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 53 55 50 50 4C 49 45 52  |SUPPLIER|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 5F 4E 55 4D 20 3D 20 44  |_NUM.=.D|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 2E 53 55 50 50 4C 49 45  |.SUPPLIE|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 52 5F 4E 55 4D 20 61 6E  |R_NUM.an|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 64 20 44 2E 53 54 41 54  |d.D.STAT|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 55 53 5F 49 4E 44 45 40  |US_INDE@|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 58 20 3C 3E 20 34 20 0A  |X.<>.4..|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 20 20 20 20 61 6E 64 20  |....and.|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 44 2E 49 4E 56 4F 49 43  |D.INVOIC|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 45 5F 44 41 54 45 20 62  |E_DATE.b|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 65 74 77 65 65 6E 20 52  |etween.R|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 45 46 2E 49 4E 56 4F 49  |EF.INVOI|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 43 45 5F 44 41 54 45 2D  |CE_DATE-|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 31 38 30 20 61 6E 64 20  |180.and.|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 40 52 45 46 2E 49 4E 56  |@REF.INV|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 4F 49 43 45 5F 44 41 54  |OICE_DAT|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 45 2B 31 38 30 20 20 0A  |E+180...|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 20 20 6C 65 66 74 20 6F  |..left.o|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 75 74 65 72 20 6A 6F 69  |uter.joi|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 6E 20 28 20 20 0A 20 20  |n.(.....|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 20 20 73 65 6C 65 63 74  |..select|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 20 46 4C 2E 44 4F 43 5F  |.FL.DOC_|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 49 40 44 2C 20 46 4C 2E  |I@D,.FL.|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 41 43 54 49 4F 4E 5F 49  |ACTION_I|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 4E 44 45 58 20 20 66 72  |NDEX..fr|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 6F 6D 20 46 4C 4F 57 5F  |om.FLOW_|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 4C 4F 47 20 46 4C 20 77  |LOG.FL.w|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 68 65 72 65 20 46 4C 2E  |here.FL.|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 41 43 54 49 4F 4E 5F 49  |ACTION_I|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 4E 44 45 58 20 3D 20 33  |NDEX.=.3|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 20 20 40 0A 20 20 20 20  |..@.....|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 20 20 61 6E 64 20 46 4C  |..and.FL|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 2E 53 45 4E 44 45 44 5F  |.SENDED_|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 54 4F 5F 54 49 4D 45 53  |TO_TIMES|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 54 41 4D 50 20 3D 20 28  |TAMP.=.(|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 73 65 6C 65 63 74 20 6D  |select.m|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 61 78 28 46 4C 32 2E 53  |ax(FL2.S|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 45 4E 44 45 44 5F 54 4F  |ENDED_TO|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 5F 54 49 40 4D 45 53 54  |_TI@MEST|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 41 4D 50 29 20 66 72 6F  |AMP).fro|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 6D 20 46 4C 4F 57 5F 4C  |m.FLOW_L|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 4F 47 20 46 4C 32 20 77  |OG.FL2.w|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 68 65 72 65 20 46 4C 2E  |here.FL.|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 44 4F 43 5F 49 44 20 3D  |DOC_ID.=|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 20 46 4C 32 2E 44 4F 43  |.FL2.DOC|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 5F 49 44 20 29 20 0A 20  |_ID.)...|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 20 20 20 29 22 46 4C 20  |...)"FL.|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 6F 6E 20 46 4C 2E 44 4F  |on.FL.DO|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 43 5F 49 44 20 3D 20 44  |C_ID.=.D|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 2E 44 4F 43 5F 49 44 20  |.DOC_ID.|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 0A 29 2C 30 29 20 0A 00  |.),0)...|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 01 01 00 00 00 00 00 00  |........|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 01 01 00 00 00 00 00 60  |.......` + "`" + `|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 00 00 00 01 C0 00 01 10  |........|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 00 00 02 03 69 01 00 60  |....i..` + "`" + `|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 00 00 00 01 30 00 01 10  |....0...|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 00 00 02 03 69 01 00 60  |....i..` + "`" + `|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 00 00 00 01 24 00 01 10  |....$...|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 00 00 02 03 69 01 00 0C  |....i...|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 00 00 00 01 01 00 01 10  |........|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 00 00 01 01 00 00 60 00  |......` + "`" + `.|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 00 00 01 12 00 01 10 00  |........|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 00 02 03 69 01 00 07 20  |...i....|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 46 30 38 46 35 43 41 43  |F08F5CAC|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 31 35 39 30 34 44 38 38  |15904D88|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 39 41 31 33 37 37 44 30  |9A1377D0|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 46 46 33 39 35 43 34 44  |FF395C4D|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 08 36 33 35 34 39 37 35  |.6354975|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 34 06 31 37 31 30 32 31  |4.171021|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 07 78 75 0A 1F 01 01 01  |.xu.....|
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: 03 33 35 35              |.355    |
(5236) [22-OCT-2020 12:44:52:420] nsbasic_bsd: exit (0)
(5236) [22-OCT-2020 12:44:52:420] nsbasic_brc: entry: oln/tot=0
(5236) [22-OCT-2020 12:44:52:420] nttfprd: entry
`,
			},
			want: want{
				`select 'NO DUPLICATES' from DUAL where 0 = nvl(( 
select 
	sum( case when FL.ACTION_INDEX = 3 then 0 else 1 end) as C 
from 
	(select :1 DOC_ID, :2 SUPPLIER_NUM, :3 INVOICE_NUM, :4 INVOICE_DATE, :5 COMP_NO from DUAL )REF join  
	DOCS D on REF.DOC_ID <> D.DOC_ID and REF.INVOICE_NUM = D.INVOICE_NUM and REF.COMP_NO = D.COMP_NO and REF.SUPPLIER_NUM = D.SUPPLIER_NUM and D.STATUS_INDEX <> 4 
	and D.INVOICE_DATE between REF.INVOICE_DATE-180 and REF.INVOICE_DATE+180  
	left outer join (  
	select FL.DOC_ID, FL.ACTION_INDEX  from FLOW_LOG FL where FL.ACTION_INDEX = 3  
		and FL.SENDED_TO_TIMESTAMP = (select max(FL2.SENDED_TO_TIMESTAMP) from FLOW_LOG FL2 where FL.DOC_ID = FL2.DOC_ID ) 
	)FL on FL.DOC_ID = D.DOC_ID 
),0) `,
				[]string{
					"F08F5CAC15904D889A1377D0FF395C4D",
					"63549754",
					"171021",
					"2017-10-31T00:00:00Z",
					"355",
				},
			},
		},
		{
			name: "Long insert with 16 parameters, some are decimal numbers",
			args: args{
				trc: `(2548) [05-NOV-2020 06:54:51:739] nsbasic_bsd: entry
				(2548) [05-NOV-2020 06:54:51:739] nsbasic_bsd: tot=0, plen=1177.
				(2548) [05-NOV-2020 06:54:51:739] nttfpwr: entry
				(2548) [05-NOV-2020 06:54:51:739] nttfpwr: socket 1284 had bytes written=1177
				(2548) [05-NOV-2020 06:54:51:739] nttfpwr: exit
				(2548) [05-NOV-2020 06:54:51:739] nsbasic_bsd: packet dump
				(2548) [05-NOV-2020 06:54:51:739] nsbasic_bsd: 04 99 00 00 06 00 00 00  |........|
				(2548) [05-NOV-2020 06:54:51:739] nsbasic_bsd: 00 00 11 69 C2 01 01 01  |...i....|
				(2548) [05-NOV-2020 06:54:51:739] nsbasic_bsd: 01 01 03 5E C3 02 80 29  |...^...)|
				(2548) [05-NOV-2020 06:54:51:739] nsbasic_bsd: 00 01 02 08 34 01 01 0D  |....4...|
				(2548) [05-NOV-2020 06:54:51:739] nsbasic_bsd: 01 01 00 01 01 00 01 01  |........|
				(2548) [05-NOV-2020 06:54:51:739] nsbasic_bsd: 10 00 01 00 01 01 01 00  |........|
				(2548) [05-NOV-2020 06:54:51:739] nsbasic_bsd: 00 01 01 00 00 00 00 00  |........|
				(2548) [05-NOV-2020 06:54:51:739] nsbasic_bsd: FE 40 55 50 44 41 54 45  |.@UPDATE|
				(2548) [05-NOV-2020 06:54:51:739] nsbasic_bsd: 20 64 6F 63 73 20 53 45  |.docs.SE|
				(2548) [05-NOV-2020 06:54:51:739] nsbasic_bsd: 54 20 73 74 61 6D 70 5F  |T.stamp_|
				(2548) [05-NOV-2020 06:54:51:739] nsbasic_bsd: 64 61 74 65 3D 20 53 59  |date=.SY|
				(2548) [05-NOV-2020 06:54:51:739] nsbasic_bsd: 53 44 41 54 45 2C 20 73  |SDATE,.s|
				(2548) [05-NOV-2020 06:54:51:739] nsbasic_bsd: 74 61 6D 70 5F 75 69 64  |tamp_uid|
				(2548) [05-NOV-2020 06:54:51:739] nsbasic_bsd: 3D 20 3A 31 2C 20 63 6F  |=.:1,.co|
				(2548) [05-NOV-2020 06:54:51:739] nsbasic_bsd: 6D 70 5F 6E 6F 20 3D 20  |mp_no.=.|
				(2548) [05-NOV-2020 06:54:51:739] nsbasic_bsd: 3A 32 40 2C 20 73 75 70  |:2@,.sup|
				(2548) [05-NOV-2020 06:54:51:739] nsbasic_bsd: 70 6C 69 65 72 5F 6E 75  |plier_nu|
				(2548) [05-NOV-2020 06:54:51:739] nsbasic_bsd: 6D 20 3D 20 3A 33 2C 20  |m.=.:3,.|
				(2548) [05-NOV-2020 06:54:51:739] nsbasic_bsd: 73 75 70 70 6C 69 65 72  |supplier|
				(2548) [05-NOV-2020 06:54:51:739] nsbasic_bsd: 5F 6E 61 6D 65 20 3D 20  |_name.=.|
				(2548) [05-NOV-2020 06:54:51:739] nsbasic_bsd: 3A 34 2C 20 69 6E 76 6F  |:4,.invo|
				(2548) [05-NOV-2020 06:54:51:739] nsbasic_bsd: 69 63 65 5F 6E 75 6D 20  |ice_num.|
				(2548) [05-NOV-2020 06:54:51:739] nsbasic_bsd: 3D 20 3A 35 2C 20 69 6E  |=.:5,.in|
				(2548) [05-NOV-2020 06:54:51:739] nsbasic_bsd: 76 6F 69 40 63 65 5F 74  |voi@ce_t|
				(2548) [05-NOV-2020 06:54:51:739] nsbasic_bsd: 79 70 65 20 3D 20 3A 36  |ype.=.:6|
				(2548) [05-NOV-2020 06:54:51:739] nsbasic_bsd: 2C 20 69 6E 76 6F 69 63  |,.invoic|
				(2548) [05-NOV-2020 06:54:51:739] nsbasic_bsd: 65 5F 64 61 74 65 20 3D  |e_date.=|
				(2548) [05-NOV-2020 06:54:51:739] nsbasic_bsd: 20 54 4F 5F 44 41 54 45  |.TO_DATE|
				(2548) [05-NOV-2020 06:54:51:739] nsbasic_bsd: 28 27 30 34 2D 31 31 2D  |('04-11-|
				(2548) [05-NOV-2020 06:54:51:739] nsbasic_bsd: 32 30 32 30 20 30 30 3A  |2020.00:|
				(2548) [05-NOV-2020 06:54:51:739] nsbasic_bsd: 30 30 3A 30 30 27 2C 27  |00:00','|
				(2548) [05-NOV-2020 06:54:51:739] nsbasic_bsd: 44 44 2D 4D 40 4D 2D 59  |DD-M@M-Y|
				(2548) [05-NOV-2020 06:54:51:739] nsbasic_bsd: 59 59 59 20 48 48 32 34  |YYY.HH24|
				(2548) [05-NOV-2020 06:54:51:739] nsbasic_bsd: 3A 4D 49 3A 53 53 27 29  |:MI:SS')|
				(2548) [05-NOV-2020 06:54:51:739] nsbasic_bsd: 2C 20 63 6F 6E 74 72 61  |,.contra|
				(2548) [05-NOV-2020 06:54:51:739] nsbasic_bsd: 63 74 5F 6E 75 6D 20 3D  |ct_num.=|
				(2548) [05-NOV-2020 06:54:51:739] nsbasic_bsd: 20 4E 55 4C 4C 2C 20 6F  |.NULL,.o|
				(2548) [05-NOV-2020 06:54:51:739] nsbasic_bsd: 72 64 65 72 5F 6E 75 6D  |rder_num|
				(2548) [05-NOV-2020 06:54:51:739] nsbasic_bsd: 20 3D 20 4E 55 4C 4C 2C  |.=.NULL,|
				(2548) [05-NOV-2020 06:54:51:739] nsbasic_bsd: 20 61 74 74 72 40 69 62  |.attr@ib|
				(2548) [05-NOV-2020 06:54:51:739] nsbasic_bsd: 5F 74 35 20 3D 20 3A 37  |_t5.=.:7|
				(2548) [05-NOV-2020 06:54:51:739] nsbasic_bsd: 2C 20 61 74 74 72 69 62  |,.attrib|
				(2548) [05-NOV-2020 06:54:51:739] nsbasic_bsd: 5F 74 33 20 3D 20 4E 55  |_t3.=.NU|
				(2548) [05-NOV-2020 06:54:51:739] nsbasic_bsd: 4C 4C 2C 20 61 74 74 72  |LL,.attr|
				(2548) [05-NOV-2020 06:54:51:739] nsbasic_bsd: 69 62 5F 74 32 20 3D 20  |ib_t2.=.|
				(2548) [05-NOV-2020 06:54:51:739] nsbasic_bsd: 3A 38 2C 20 61 74 74 72  |:8,.attr|
				(2548) [05-NOV-2020 06:54:51:739] nsbasic_bsd: 69 62 5F 74 36 20 3D 20  |ib_t6.=.|
				(2548) [05-NOV-2020 06:54:51:739] nsbasic_bsd: 4E 55 4C 4C 2C 20 40 69  |NULL,.@i|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 6E 76 6F 69 63 65 5F 63  |nvoice_c|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 75 72 72 65 6E 63 79 20  |urrency.|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 3D 20 3A 39 2C 20 65 78  |=.:9,.ex|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 63 68 61 6E 67 65 5F 72  |change_r|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 61 74 65 20 3D 20 3A 31  |ate.=.:1|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 30 2C 20 69 6E 76 6F 69  |0,.invoi|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 63 65 5F 73 75 6D 20 3D  |ce_sum.=|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 20 3A 31 31 2C 20 69 40  |.:11,.i@|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 6E 76 6F 69 63 65 5F 73  |nvoice_s|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 75 6D 5F 63 61 6C 63 20  |um_calc.|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 3D 20 3A 31 32 2C 20 76  |=.:12,.v|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 61 74 5F 73 75 6D 20 3D  |at_sum.=|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 20 3A 31 33 2C 20 6E 65  |.:13,.ne|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 74 5F 73 75 6D 20 3D 20  |t_sum.=.|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 3A 31 34 2C 20 61 74 74  |:14,.att|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 72 69 62 5F 64 32 20 3D  |rib_d2.=|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 40 20 4E 55 4C 4C 2C 20  |@.NULL,.|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 69 6E 76 6F 69 63 65 5F  |invoice_|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 6C 61 73 74 5F 64 61 74  |last_dat|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 65 20 3D 20 54 4F 5F 44  |e.=.TO_D|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 41 54 45 28 27 32 34 2D  |ATE('24-|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 31 31 2D 32 30 32 30 20  |11-2020.|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 30 30 3A 30 30 3A 30 30  |00:00:00|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 27 2C 27 44 44 2D 4D 4D  |','DD-MM|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 2D 40 59 59 59 59 20 48  |-@YYYY.H|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 48 32 34 3A 4D 49 3A 53  |H24:MI:S|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 53 27 29 2C 20 63 61 73  |S'),.cas|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 68 5F 64 61 74 65 20 3D  |h_date.=|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 20 4E 55 4C 4C 2C 20 61  |.NULL,.a|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 74 74 72 69 62 5F 74 34  |ttrib_t4|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 20 3D 20 4E 55 4C 4C 2C  |.=.NULL,|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 20 61 74 74 72 69 62 5F  |.attrib_|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 74 37 40 20 3D 20 4E 55  |t7@.=.NU|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 4C 4C 2C 20 65 6E 74 72  |LL,.entr|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 79 5F 64 61 74 65 20 3D  |y_date.=|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 20 4E 55 4C 4C 2C 20 76  |.NULL,.v|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 6F 75 63 68 65 72 5F 6E  |oucher_n|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 75 6D 20 3D 20 4E 55 4C  |um.=.NUL|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 4C 2C 20 70 61 79 6D 65  |L,.payme|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 6E 74 5F 64 61 74 65 20  |nt_date.|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 3D 20 4E 3C 55 4C 4C 2C  |=.N<ULL,|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 20 61 74 74 72 69 62 5F  |.attrib_|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 74 31 20 3D 20 4E 55 4C  |t1.=.NUL|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 4C 2C 20 6E 65 74 5F 73  |L,.net_s|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 75 6D 5F 63 61 6C 63 20  |um_calc.|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 3D 20 3A 31 35 20 57 48  |=.:15.WH|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 45 52 45 20 64 6F 63 5F  |ERE.doc_|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 69 64 20 3D 20 3A 31 36  |id.=.:16|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 00 01 01 01 01 00 00 00  |........|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 00 00 01 02 00 00 00 00  |........|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 00 60 01 00 00 01 BA 00  |.` + "`" + `......|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 01 10 00 00 02 03 69 01  |......i.|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 00 60 01 00 00 01 12 00  |.` + "`" + `......|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 01 10 00 00 02 03 69 01  |......i.|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 00 60 01 00 00 01 30 00  |.` + "`" + `....0.|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 01 10 00 00 02 03 69 01  |......i.|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 00 60 01 00 00 01 24 00  |.` + "`" + `....$.|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 01 10 00 00 02 03 69 01  |......i.|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 00 60 01 00 00 01 30 00  |.` + "`" + `....0.|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 01 10 00 00 02 03 69 01  |......i.|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 00 60 01 00 00 01 0C 00  |.` + "`" + `......|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 01 10 00 00 02 03 69 01  |......i.|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 00 60 01 00 00 01 48 00  |.` + "`" + `....H.|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 01 10 00 00 02 03 69 01  |......i.|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 00 60 01 00 00 01 96 00  |.` + "`" + `......|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 01 10 00 00 02 03 69 01  |......i.|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 00 60 01 00 00 01 12 00  |.` + "`" + `......|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 01 10 00 00 02 03 69 01  |......i.|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 00 02 01 00 00 01 15 00  |........|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 00 00 00 00 00 00 02 01  |........|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 00 00 01 15 00 00 00 00  |........|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 00 00 00 02 01 00 00 01  |........|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 15 00 00 00 00 00 00 00  |........|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 02 01 00 00 01 15 00 00  |........|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 00 00 00 00 00 02 01 00  |........|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 00 01 15 00 00 00 00 00  |........|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 00 00 02 01 00 00 01 15  |........|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 00 00 00 00 00 00 00 60  |.......` + "`" + `|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 01 00 00 01 78 00 01 10  |....x...|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 00 00 02 03 69 01 00 07  |....i...|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 1F 43 41 53 53 41 4E 20  |.CASSAN.|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 4A 45 41 4E 2D 46 52 41  |JEAN-FRA|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 4E 43 4F 49 53 20 28 45  |NCOIS.(E|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 58 54 45 52 4E 41 4C 29  |XTERNAL)|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 03 39 39 32 08 39 39 32  |.992.992|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 5F 34 31 38 39 06 46 4C  |_4189.FL|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 55 58 59 4D 08 54 45 53  |UXYM.TES|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 54 31 32 33 34 02 49 4E  |T1234.IN|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 0C 53 53 43 5F 53 55 50  |.SSC_SUP|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 50 4C 49 45 52 19 43 4F  |PLIER.CO|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 4E 54 52 41 43 54 20 5F  |NTRACT._|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 20 43 4F 4E 54 52 41 43  |.CONTRAC|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 54 20 49 53 53 55 45 03  |T.ISSUE.|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 45 55 52 02 C1 02 04 C2  |EUR.....|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 02 36 0D 04 C2 02 36 0D  |.6....6.|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 01 80 04 C2 02 36 0D 04  |.....6..|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: C2 02 36 0D 14 35 30 46  |..6..50F|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 45 45 42 38 42 33 33 38  |EEB8B338|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 34 34 45 36 39 42 46 38  |44E69BF8|
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: 43                       |C       |
				(2548) [05-NOV-2020 06:54:51:740] nsbasic_bsd: exit (0)
`,
			},
			want: want{
				"UPDATE docs SET stamp_date= SYSDATE, stamp_uid= :1, comp_no = :2, supplier_num = :3, supplier_name = :4, invoice_num = :5, invoice_type = :6, invoice_date = TO_DATE('04-11-2020 00:00:00','DD-MM-YYYY HH24:MI:SS'), contract_num = NULL, order_num = NULL, attrib_t5 = :7, attrib_t3 = NULL, attrib_t2 = :8, attrib_t6 = NULL, invoice_currency = :9, exchange_rate = :10, invoice_sum = :11, invoice_sum_calc = :12, vat_sum = :13, net_sum = :14, attrib_d2 = NULL, invoice_last_date = TO_DATE('24-11-2020 00:00:00','DD-MM-YYYY HH24:MI:SS'), cash_date = NULL, attrib_t4 = NULL, attrib_t7 = NULL, entry_date = NULL, voucher_num = NULL, payment_date = NULL, attrib_t1 = NULL, net_sum_calc = :15 WHERE doc_id = :16",
				[]string{
					"CASSAN JEAN-FRANCOIS (EXTERNAL)",
					"992",
					"992_4189",
					"FLUXYM",
					"TEST1234",
					"IN",
					"SSC_SUPPLIER",
					"CONTRACT _ CONTRACT ISSUE",
					"EUR",
					"1",
					"153.12",
					"153.12",
					"0",
					"153.12",
					"153.12",
					"50FEEB8B33844E69BF8C",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getQueryFromTraceSnippet(tt.args.trc)
			if err != nil {
				t.Errorf("Error returned error = %v", err)
				return
			}
			if got == nil {
				t.Errorf("Query not found")
				return
			}
			if !significantCharsEqual(got.Query, tt.want.query) {
				t.Errorf("Query obtained \n%v,\n want \n%v", got, tt.want)
			}
			if len(got.Params) != len(tt.want.params) {
				t.Errorf("Parameters number = \n%v,\n want \n%v", len(got.Params), len(tt.want.params))
			} else {
				for i := 0; i < len(tt.want.params); i++ {
					if got.Params[i].String() != tt.want.params[i] {
						t.Errorf("Parameters #%d = \n%v,\n want \n%v", i+1, got.Params[i].String(), tt.want.params[i])
					}
				}
			}
		})
	}
}
