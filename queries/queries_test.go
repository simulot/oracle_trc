package queries

import (
	"errors"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/simulot/oracle_trc/ts"
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
func getQueryFromTraceSnippet(trc string) (string, error) {
	p := New(strings.NewReader(trc), "test", ts.OracleTS_DD_MON_YYYY_HH_MI_SS_FF9)

	for {
		pk, err := p.NextPacket()
		if err != nil {
			return "", err
		}
		if pk == nil {
			return "", nil
		}
		pl := pk.payload
		pos := detectQuery(pl)
		if pos < 0 {
			continue
		}
		return query(pl[pos:]), nil
	}
	return "", errors.New("Should not happen")
}
func Test_getQueryFromTraceSnippet(t *testing.T) {
	type args struct {
		trc string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "simple short select",
			args: args{
				trc: "(3224) [01-MAR-2019 17:27:32:320] nsbasic_bsd: packet dump\r\n" +
					"(3224) [01-MAR-2019 17:27:32:320] nsbasic_bsd: 00 5D 00 00 06 00 00 00  |.]......|\r\n" +
					"(3224) [01-MAR-2019 17:27:32:320] nsbasic_bsd: 00 00 11 69 4A 01 01 01  |...iJ...|\r\n" +
					"(3224) [01-MAR-2019 17:27:32:320] nsbasic_bsd: 01 02 03 5E 4B 02 80 61  |...^K..a|\r\n" +
					"(3224) [01-MAR-2019 17:27:32:320] nsbasic_bsd: 00 01 01 45 01 01 0D 01  |...E....|\r\n" +
					"(3224) [01-MAR-2019 17:27:32:321] nsbasic_bsd: 01 00 01 64 00 00 00 00  |...d....|\r\n" +
					"(3224) [01-MAR-2019 17:27:32:321] nsbasic_bsd: 01 00 01 01 01 00 00 01  |........|\r\n" +
					"(3224) [01-MAR-2019 17:27:32:321] nsbasic_bsd: 01 00 00 00 00 00 17 73  |.......s|\r\n" +
					"(3224) [01-MAR-2019 17:27:32:321] nsbasic_bsd: 65 6C 65 63 74 20 64 6F  |elect.do|\r\n" +
					"(3224) [01-MAR-2019 17:27:32:321] nsbasic_bsd: 63 5F 69 64 20 66 72 6F  |c_id.fro|\r\n" +
					"(3224) [01-MAR-2019 17:27:32:321] nsbasic_bsd: 6D 20 44 4F 43 53 01 01  |m.DOCS..|\r\n" +
					"(3224) [01-MAR-2019 17:27:32:321] nsbasic_bsd: 00 00 00 00 00 00 01 01  |........|\r\n" +
					"(3224) [01-MAR-2019 17:27:32:321] nsbasic_bsd: 00 00 00 00 00           |.....   |\r\n" +
					"(3224) [01-MAR-2019 17:27:32:321] nsbasic_bsd: exit (0)",
			},
			want:    "select doc_id from DOCS",
			wantErr: false,
		},
		{
			name: "simple insert",
			args: args{
				trc: "(7980) [28-FEB-2019 20:29:23:254] nsbasic_bsd: packet dump\r\n" +
					"(7980) [28-FEB-2019 20:29:23:254] nsbasic_bsd: 01 4B 00 00 06 00 00 00  |.K......|\r\n" +
					"(7980) [28-FEB-2019 20:29:23:254] nsbasic_bsd: 00 00 11 69 39 01 01 01  |...i9...|\r\n" +
					"(7980) [28-FEB-2019 20:29:23:254] nsbasic_bsd: 01 01 03 5E 3A 02 80 29  |...^:..)|\r\n" +
					"(7980) [28-FEB-2019 20:29:23:254] nsbasic_bsd: 00 01 02 01 59 01 01 0D  |....Y...|\r\n" +
					"(7980) [28-FEB-2019 20:29:23:254] nsbasic_bsd: 01 01 00 01 01 00 01 01  |........|\r\n" +
					"(7980) [28-FEB-2019 20:29:23:254] nsbasic_bsd: 04 00 01 00 01 01 01 00  |........|\r\n" +
					"(7980) [28-FEB-2019 20:29:23:254] nsbasic_bsd: 00 01 01 00 00 00 00 00  |........|\r\n" +
					"(7980) [28-FEB-2019 20:29:23:254] nsbasic_bsd: FE 40 49 4E 53 45 52 54  |.@INSERT|\r\n" +
					"(7980) [28-FEB-2019 20:29:23:254] nsbasic_bsd: 20 49 4E 54 4F 20 66 6C  |.INTO.fl|\r\n" +
					"(7980) [28-FEB-2019 20:29:23:254] nsbasic_bsd: 6F 77 5F 6C 6F 67 20 28  |ow_log.(|\r\n" +
					"(7980) [28-FEB-2019 20:29:23:254] nsbasic_bsd: 64 6F 63 5F 69 64 2C 73  |doc_id,s|\r\n" +
					"(7980) [28-FEB-2019 20:29:23:254] nsbasic_bsd: 65 6E 64 65 72 5F 6E 61  |ender_na|\r\n" +
					"(7980) [28-FEB-2019 20:29:23:254] nsbasic_bsd: 6D 65 2C 61 63 74 69 6F  |me,actio|\r\n" +
					"(7980) [28-FEB-2019 20:29:23:254] nsbasic_bsd: 6E 5F 6E 61 6D 65 2C 73  |n_name,s|\r\n" +
					"(7980) [28-FEB-2019 20:29:23:254] nsbasic_bsd: 65 6E 64 65 64 5F 74 6F  |ended_to|\r\n" +
					"(7980) [28-FEB-2019 20:29:23:254] nsbasic_bsd: 5F 74 33 69 6D 65 73 74  |_t3imest|\r\n" +
					"(7980) [28-FEB-2019 20:29:23:254] nsbasic_bsd: 61 6D 70 2C 61 63 74 69  |amp,acti|\r\n" +
					"(7980) [28-FEB-2019 20:29:23:254] nsbasic_bsd: 6F 6E 5F 69 6E 64 65 78  |on_index|\r\n" +
					"(7980) [28-FEB-2019 20:29:23:254] nsbasic_bsd: 29 20 56 41 4C 55 45 53  |).VALUES|\r\n" +
					"(7980) [28-FEB-2019 20:29:23:254] nsbasic_bsd: 20 28 3A 31 2C 3A 32 2C  |.(:1,:2,|\r\n" +
					"(7980) [28-FEB-2019 20:29:23:254] nsbasic_bsd: 3A 33 2C 53 59 53 44 41  |:3,SYSDA|\r\n" +
					"(7980) [28-FEB-2019 20:29:23:254] nsbasic_bsd: 54 45 2C 3A 34 29 00 01  |TE,:4)..|\r\n" +
					"(7980) [28-FEB-2019 20:29:23:254] nsbasic_bsd: 01 01 01 00 00 00 00 00  |........|\r\n" +
					"(7980) [28-FEB-2019 20:29:23:254] nsbasic_bsd: 01 04 00 00 00 00 00 60  |.......`|\r\n" +
					"(7980) [28-FEB-2019 20:29:23:254] nsbasic_bsd: 01 00 00 01 C0 00 01 10  |........|\r\n" +
					"(7980) [28-FEB-2019 20:29:23:254] nsbasic_bsd: 00 00 02 03 69 01 00 60  |....i..`|\r\n" +
					"(7980) [28-FEB-2019 20:29:23:254] nsbasic_bsd: 01 00 00 01 BA 00 01 10  |........|\r\n" +
					"(7980) [28-FEB-2019 20:29:23:254] nsbasic_bsd: 00 00 02 03 69 01 00 60  |....i..`|\r\n" +
					"(7980) [28-FEB-2019 20:29:23:254] nsbasic_bsd: 01 00 00 01 3C 00 01 10  |....<...|\r\n" +
					"(7980) [28-FEB-2019 20:29:23:254] nsbasic_bsd: 00 00 02 03 69 01 00 02  |....i...|\r\n" +
					"(7980) [28-FEB-2019 20:29:23:254] nsbasic_bsd: 01 00 00 01 15 00 00 00  |........|\r\n" +
					"(7980) [28-FEB-2019 20:29:23:254] nsbasic_bsd: 00 00 00 00 07 20 34 35  |......45|\r\n" +
					"(7980) [28-FEB-2019 20:29:23:254] nsbasic_bsd: 30 37 41 33 34 34 39 42  |07A3449B|\r\n" +
					"(7980) [28-FEB-2019 20:29:23:254] nsbasic_bsd: 38 44 34 32 41 44 41 38  |8D42ADA8|\r\n" +
					"(7980) [28-FEB-2019 20:29:23:254] nsbasic_bsd: 32 44 34 45 38 34 39 32  |2D4E8492|\r\n" +
					"(7980) [28-FEB-2019 20:29:23:254] nsbasic_bsd: 46 35 38 46 41 32 1F 43  |F58FA2.C|\r\n" +
					"(7980) [28-FEB-2019 20:29:23:254] nsbasic_bsd: 41 53 53 41 4E 20 4A 45  |ASSAN.JE|\r\n" +
					"(7980) [28-FEB-2019 20:29:23:254] nsbasic_bsd: 41 4E 2D 46 52 41 4E 43  |AN-FRANC|\r\n" +
					"(7980) [28-FEB-2019 20:29:23:254] nsbasic_bsd: 4F 49 53 20 28 45 58 54  |OIS.(EXT|\r\n" +
					"(7980) [28-FEB-2019 20:29:23:254] nsbasic_bsd: 45 52 4E 41 4C 29 0A 54  |ERNAL).T|\r\n" +
					"(7980) [28-FEB-2019 20:29:23:254] nsbasic_bsd: 61 6C 6C 65 74 65 74 74  |alletett|\r\n" +
					"(7980) [28-FEB-2019 20:29:23:254] nsbasic_bsd: 75 01 80                 |u..     |\r\n" +
					"(7980) [28-FEB-2019 20:29:23:254] nsbasic_bsd: exit (0)				`,",
			},
			want:    "INSERT INTO flow_log (doc_id,sender_name,action_name,sended_to_timestamp,action_index) VALUES (:1,:2,:3,SYSDATE,:4)",
			wantErr: false,
		},
		{
			name: "simple select lowcase",
			args: args{
				trc: "(6604) [12-FEB-2019 16:40:00:651] nsbasic_bsd: packet dump\r\n" +
					"(6604) [12-FEB-2019 16:40:00:651] nsbasic_bsd: 00 C3 00 00 06 00 00 00  |........|\r\n" +
					"(6604) [12-FEB-2019 16:40:00:651] nsbasic_bsd: 00 00 11 69 0D 01 01 01  |...i....|\r\n" +
					"(6604) [12-FEB-2019 16:40:00:651] nsbasic_bsd: 01 01 03 5E 0E 02 80 61  |...^...a|\r\n" +
					"(6604) [12-FEB-2019 16:40:00:651] nsbasic_bsd: 00 01 02 01 6B 01 01 0D  |....k...|\r\n" +
					"(6604) [12-FEB-2019 16:40:00:651] nsbasic_bsd: 01 01 00 01 64 00 00 00  |....d...|\r\n" +
					"(6604) [12-FEB-2019 16:40:00:651] nsbasic_bsd: 00 01 00 01 01 01 00 00  |........|\r\n" +
					"(6604) [12-FEB-2019 16:40:00:651] nsbasic_bsd: 01 01 00 00 00 00 00 FE  |........|\r\n" +
					"(6604) [12-FEB-2019 16:40:00:651] nsbasic_bsd: 40 73 65 6C 65 63 74 20  |@select.|\r\n" +
					"(6604) [12-FEB-2019 16:40:00:651] nsbasic_bsd: 43 4F 55 4E 54 28 2A 29  |COUNT(*)|\r\n" +
					"(6604) [12-FEB-2019 16:40:00:651] nsbasic_bsd: 20 66 72 6F 6D 20 75 73  |.from.us|\r\n" +
					"(6604) [12-FEB-2019 16:40:00:651] nsbasic_bsd: 65 72 5F 74 61 62 5F 63  |er_tab_c|\r\n" +
					"(6604) [12-FEB-2019 16:40:00:651] nsbasic_bsd: 6F 6C 75 6D 6E 73 20 77  |olumns.w|\r\n" +
					"(6604) [12-FEB-2019 16:40:00:651] nsbasic_bsd: 68 65 72 65 20 74 61 62  |here.tab|\r\n" +
					"(6604) [12-FEB-2019 16:40:00:651] nsbasic_bsd: 6C 65 5F 6E 61 6D 65 20  |le_name.|\r\n" +
					"(6604) [12-FEB-2019 16:40:00:651] nsbasic_bsd: 3D 20 4E 27 42 57 5F 55  |=.N'BW_U|\r\n" +
					"(6604) [12-FEB-2019 16:40:00:651] nsbasic_bsd: 53 39 45 52 5F 41 55 54  |S9ER_AUT|\r\n" +
					"(6604) [12-FEB-2019 16:40:00:651] nsbasic_bsd: 48 45 4E 54 49 43 41 54  |HENTICAT|\r\n" +
					"(6604) [12-FEB-2019 16:40:00:651] nsbasic_bsd: 49 4F 4E 27 20 61 6E 64  |ION'.and|\r\n" +
					"(6604) [12-FEB-2019 16:40:00:651] nsbasic_bsd: 20 63 6F 6C 75 6D 6E 5F  |.column_|\r\n" +
					"(6604) [12-FEB-2019 16:40:00:651] nsbasic_bsd: 6E 61 6D 65 20 3D 20 4E  |name.=.N|\r\n" +
					"(6604) [12-FEB-2019 16:40:00:651] nsbasic_bsd: 27 55 41 55 5F 46 4F 52  |'UAU_FOR|\r\n" +
					"(6604) [12-FEB-2019 16:40:00:651] nsbasic_bsd: 47 4F 54 5F 41 43 54 49  |GOT_ACTI|\r\n" +
					"(6604) [12-FEB-2019 16:40:00:651] nsbasic_bsd: 4F 4E 27 00 01 01 00 00  |ON'.....|\r\n" +
					"(6604) [12-FEB-2019 16:40:00:651] nsbasic_bsd: 00 00 00 00 01 01 00 00  |........|\r\n" +
					"(6604) [12-FEB-2019 16:40:00:651] nsbasic_bsd: 00 00 00                 |...     |\r\n" +
					"(6604) [12-FEB-2019 16:40:00:651] nsbasic_bsd: exit (0)",
			},
			want:    "select COUNT(*) from user_tab_columns where table_name = N'BW_USER_AUTHENTICATION' and column_name = N'UAU_FORGOT_ACTION'",
			wantErr: false,
		},
		{
			name: "compound insert with select uppercase",
			args: args{
				trc: "(9520) [26-FEB-2019 11:52:23:665] nsbasic_bsd: packet dump\r\n" +
					"(9520) [26-FEB-2019 11:52:23:665] nsbasic_bsd: 01 4F 00 00 06 00 00 00  |.O......|\r\n" +
					"(9520) [26-FEB-2019 11:52:23:665] nsbasic_bsd: 00 00 11 69 11 01 01 01  |...i....|\r\n" +
					"(9520) [26-FEB-2019 11:52:23:665] nsbasic_bsd: 01 01 03 5E 12 02 81 29  |...^...)|\r\n" +
					"(9520) [26-FEB-2019 11:52:23:665] nsbasic_bsd: 00 01 02 01 C5 01 01 0D  |........|\r\n" +
					"(9520) [26-FEB-2019 11:52:23:665] nsbasic_bsd: 01 01 00 01 01 00 01 01  |........|\r\n" +
					"(9520) [26-FEB-2019 11:52:23:665] nsbasic_bsd: 04 00 01 00 01 01 01 00  |........|\r\n" +
					"(9520) [26-FEB-2019 11:52:23:665] nsbasic_bsd: 00 01 01 00 00 00 00 00  |........|\r\n" +
					"(9520) [26-FEB-2019 11:52:23:665] nsbasic_bsd: FE 40 49 4E 53 45 52 54  |.@INSERT|\r\n" +
					"(9520) [26-FEB-2019 11:52:23:665] nsbasic_bsd: 20 49 4E 54 4F 20 42 57  |.INTO.BW|\r\n" +
					"(9520) [26-FEB-2019 11:52:23:665] nsbasic_bsd: 5F 55 53 45 52 5F 41 55  |_USER_AU|\r\n" +
					"(9520) [26-FEB-2019 11:52:23:665] nsbasic_bsd: 54 48 5F 4C 4F 47 20 28  |TH_LOG.(|\r\n" +
					"(9520) [26-FEB-2019 11:52:23:665] nsbasic_bsd: 55 41 4C 5F 49 44 2C 55  |UAL_ID,U|\r\n" +
					"(9520) [26-FEB-2019 11:52:23:665] nsbasic_bsd: 41 4C 5F 55 53 45 52 5F  |AL_USER_|\r\n" +
					"(9520) [26-FEB-2019 11:52:23:665] nsbasic_bsd: 4E 45 54 57 4F 52 4B 5F  |NETWORK_|\r\n" +
					"(9520) [26-FEB-2019 11:52:23:665] nsbasic_bsd: 4E 41 4D 45 2C 55 41 4C  |NAME,UAL|\r\n" +
					"(9520) [26-FEB-2019 11:52:23:665] nsbasic_bsd: 5F 4C 40 4F 47 5F 43 4F  |_L@OG_CO|\r\n" +
					"(9520) [26-FEB-2019 11:52:23:665] nsbasic_bsd: 44 45 2C 55 41 4C 5F 53  |DE,UAL_S|\r\n" +
					"(9520) [26-FEB-2019 11:52:23:665] nsbasic_bsd: 54 41 4D 50 5F 44 41 54  |TAMP_DAT|\r\n" +
					"(9520) [26-FEB-2019 11:52:23:665] nsbasic_bsd: 45 2C 55 41 4C 5F 41 50  |E,UAL_AP|\r\n" +
					"(9520) [26-FEB-2019 11:52:23:665] nsbasic_bsd: 50 4C 49 43 41 54 49 4F  |PLICATIO|\r\n" +
					"(9520) [26-FEB-2019 11:52:23:665] nsbasic_bsd: 4E 29 20 56 41 4C 55 45  |N).VALUE|\r\n" +
					"(9520) [26-FEB-2019 11:52:23:665] nsbasic_bsd: 53 20 28 3A 31 2C 3A 32  |S.(:1,:2|\r\n" +
					"(9520) [26-FEB-2019 11:52:23:665] nsbasic_bsd: 2C 3A 33 2C 28 53 45 4C  |,:3,(SEL|\r\n" +
					"(9520) [26-FEB-2019 11:52:23:665] nsbasic_bsd: 45 43 54 17 20 53 59 53  |ECT..SYS|\r\n" +
					"(9520) [26-FEB-2019 11:52:23:665] nsbasic_bsd: 44 41 54 45 20 46 52 4F  |DATE.FRO|\r\n" +
					"(9520) [26-FEB-2019 11:52:23:665] nsbasic_bsd: 4D 20 44 55 41 4C 29 2C  |M.DUAL),|\r\n" +
					"(9520) [26-FEB-2019 11:52:23:665] nsbasic_bsd: 3A 34 29 00 01 01 01 01  |:4).....|\r\n" +
					"(9520) [26-FEB-2019 11:52:23:665] nsbasic_bsd: 00 00 00 00 00 01 04 00  |........|\r\n" +
					"(9520) [26-FEB-2019 11:52:23:665] nsbasic_bsd: 00 00 00 00 60 01 00 00  |....x...|\r\n" +
					"(9520) [26-FEB-2019 11:52:23:665] nsbasic_bsd: 01 78 00 01 10 00 00 02  |.x......|\r\n" +
					"(9520) [26-FEB-2019 11:52:23:666] nsbasic_bsd: 03 69 01 00 60 01 00 00  |.i..s...|\r\n" +
					"(9520) [26-FEB-2019 11:52:23:666] nsbasic_bsd: 01 36 00 01 10 00 00 02  |.6......|\r\n" +
					"(9520) [26-FEB-2019 11:52:23:666] nsbasic_bsd: 03 69 01 00 02 01 00 00  |.i......|\r\n" +
					"(9520) [26-FEB-2019 11:52:23:666] nsbasic_bsd: 01 16 00 00 00 00 00 00  |........|\r\n" +
					"(9520) [26-FEB-2019 11:52:23:666] nsbasic_bsd: 00 60 01 00 00 01 3C 00  |.s....<.|\r\n" +
					"(9520) [26-FEB-2019 11:52:23:666] nsbasic_bsd: 01 10 00 00 02 03 69 01  |......i.|\r\n" +
					"(9520) [26-FEB-2019 11:52:23:666] nsbasic_bsd: 00 07 14 44 38 33 42 36  |...D83B6|\r\n" +
					"(9520) [26-FEB-2019 11:52:23:666] nsbasic_bsd: 30 30 36 38 31 39 42 34  |006819B4|\r\n" +
					"(9520) [26-FEB-2019 11:52:23:666] nsbasic_bsd: 46 30 33 42 31 45 43 09  |F03B1EC.|\r\n" +
					"(9520) [26-FEB-2019 11:52:23:666] nsbasic_bsd: 4A 46 2E 43 41 53 53 41  |JF.CASSA|\r\n" +
					"(9520) [26-FEB-2019 11:52:23:666] nsbasic_bsd: 4E 02 C1 02 0A 4D 61 73  |N....Mas|\r\n" +
					"(9520) [26-FEB-2019 11:52:23:666] nsbasic_bsd: 74 65 72 20 35 2E 31     |ter.5.1 |\r\n" +
					"(9520) [26-FEB-2019 11:52:23:666] nsbasic_bsd: exit (0)",
			},
			want:    "INSERT INTO BW_USER_AUTH_LOG (UAL_ID,UAL_USER_NETWORK_NAME,UAL_LOG_CODE,UAL_STAMP_DATE,UAL_APPLICATION) VALUES (:1,:2,:3,(SELECT SYSDATE FROM DUAL),:4)",
			wantErr: false,
		},
		{
			name: "compound selects with blank chars before",
			args: args{
				trc: "(6632) [26-FEB-2019 11:49:09:541] nsbasic_bsd: packet dump\r\n" +
					"(6632) [26-FEB-2019 11:49:09:541] nsbasic_bsd: 03 70 00 00 06 00 00 00  |.p......|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:541] nsbasic_bsd: 00 00 11 69 08 01 01 01  |...i....|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:541] nsbasic_bsd: 01 02 03 5E 09 02 80 69  |...^...i|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:541] nsbasic_bsd: 00 01 02 08 B8 01 01 0D  |........|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:541] nsbasic_bsd: 01 01 00 01 19 00 01 01  |........|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:541] nsbasic_bsd: 01 00 01 00 01 01 01 00  |........|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:541] nsbasic_bsd: 00 01 01 00 00 00 00 00  |........|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:541] nsbasic_bsd: FE 40 0D 0A 09 09 09 09  |.@......|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:541] nsbasic_bsd: 09 09 73 65 6C 65 63 74  |..select|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:541] nsbasic_bsd: 20 63 6F 75 6E 74 28 2A  |.count(*|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:541] nsbasic_bsd: 29 20 66 72 6F 6D 20 28  |).from.(|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:541] nsbasic_bsd: 0D 0A 09 09 09 09 09 09  |........|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:541] nsbasic_bsd: 77 69 74 68 20 44 20 61  |with.D.a|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:541] nsbasic_bsd: 73 20 28 0D 0A 09 09 09  |s.(.....|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:541] nsbasic_bsd: 09 09 09 09 73 65 6C 65  |....sele|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:541] nsbasic_bsd: 63 74 40 20 44 4F 43 5F  |ct@.DOC_|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:541] nsbasic_bsd: 49 44 2C 53 54 41 54 55  |ID,STATU|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:541] nsbasic_bsd: 53 5F 49 4E 44 45 58 0D  |S_INDEX.|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:541] nsbasic_bsd: 0A 09 09 09 09 09 09 09  |........|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:541] nsbasic_bsd: 66 72 6F 6D 20 44 4F 43  |from.DOC|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:541] nsbasic_bsd: 53 20 77 68 65 72 65 20  |S.where.|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:541] nsbasic_bsd: 53 54 41 54 55 53 5F 49  |STATUS_I|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:541] nsbasic_bsd: 4E 44 45 58 20 69 6E 20  |NDEX.in.|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:541] nsbasic_bsd: 28 32 2C 40 39 39 29 0D  |(2,@99).|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:541] nsbasic_bsd: 0A 09 09 09 09 09 09 29  |.......)|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:541] nsbasic_bsd: 2C 6C 61 73 74 5F 41 4C  |,last_AL|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:541] nsbasic_bsd: 20 61 73 20 28 0D 0A 09  |.as.(...|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:541] nsbasic_bsd: 09 09 09 09 09 09 73 65  |......se|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:541] nsbasic_bsd: 6C 65 63 74 20 61 6C 2E  |lect.al.|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:541] nsbasic_bsd: 64 6F 63 5F 69 64 2C 20  |doc_id,.|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:541] nsbasic_bsd: 6D 61 78 28 61 6C 2E 73  |max(al.s|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:541] nsbasic_bsd: 74 61 6D 70 40 5F 64 61  |tamp@_da|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:541] nsbasic_bsd: 74 65 29 20 73 74 61 6D  |te).stam|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:541] nsbasic_bsd: 70 5F 64 61 74 65 0D 0A  |p_date..|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:541] nsbasic_bsd: 09 09 09 09 09 09 09 66  |.......f|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:541] nsbasic_bsd: 72 6F 6D 20 61 63 74 69  |rom.acti|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:542] nsbasic_bsd: 6F 6E 5F 6C 6F 67 20 61  |on_log.a|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:542] nsbasic_bsd: 6C 20 6A 6F 69 6E 20 44  |l.join.D|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:542] nsbasic_bsd: 20 6F 6E 20 44 2E 44 4F  |.on.D.DO|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:542] nsbasic_bsd: 43 5F 49 44 20 40 3D 20  |C_ID.@=.|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:542] nsbasic_bsd: 61 6C 2E 44 4F 43 5F 49  |al.DOC_I|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:542] nsbasic_bsd: 44 0D 0A 09 09 09 09 09  |D.......|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:542] nsbasic_bsd: 09 09 67 72 6F 75 70 20  |..group.|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:542] nsbasic_bsd: 62 79 20 61 6C 2E 64 6F  |by.al.do|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:542] nsbasic_bsd: 63 5F 69 64 0D 0A 09 09  |c_id....|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:542] nsbasic_bsd: 09 09 09 09 29 2C 20 0D  |....),..|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:542] nsbasic_bsd: 0A 09 09 09 09 09 09 6C  |.......l|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:542] nsbasic_bsd: 61 73 74 41 4C 5F 40 55  |astAL_@U|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:542] nsbasic_bsd: 49 44 20 61 73 20 28 0D  |ID.as.(.|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:542] nsbasic_bsd: 0A 09 09 09 09 09 09 09  |........|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:542] nsbasic_bsd: 73 65 6C 65 63 74 20 64  |select.d|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:542] nsbasic_bsd: 69 73 74 69 6E 63 74 20  |istinct.|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:542] nsbasic_bsd: 41 4C 2E 64 6F 63 5F 69  |AL.doc_i|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:542] nsbasic_bsd: 64 2C 41 4C 2E 73 74 61  |d,AL.sta|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:542] nsbasic_bsd: 6D 70 5F 75 69 64 20 0D  |mp_uid..|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:542] nsbasic_bsd: 0A 09 09 09 09 09 09 40  |.......@|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:542] nsbasic_bsd: 09 66 72 6F 6D 20 44 20  |.from.D.|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:542] nsbasic_bsd: 6A 6F 69 6E 20 6C 61 73  |join.las|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:542] nsbasic_bsd: 74 5F 41 4C 20 4C 41 4C  |t_AL.LAL|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:542] nsbasic_bsd: 20 6F 6E 20 44 2E 44 4F  |.on.D.DO|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:542] nsbasic_bsd: 43 5F 49 44 20 3D 20 4C  |C_ID.=.L|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:542] nsbasic_bsd: 41 4C 2E 44 4F 43 5F 49  |AL.DOC_I|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:542] nsbasic_bsd: 44 0D 0A 09 09 09 09 09  |D.......|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:542] nsbasic_bsd: 09 09 6A 6F 69 6E 20 61  |..join.a|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:542] nsbasic_bsd: 40 63 74 69 6F 6E 5F 6C  |@ction_l|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:542] nsbasic_bsd: 6F 67 20 41 4C 20 6F 6E  |og.AL.on|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:542] nsbasic_bsd: 20 4C 41 4C 2E 44 4F 43  |.LAL.DOC|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:542] nsbasic_bsd: 5F 49 44 20 3D 20 41 4C  |_ID.=.AL|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:542] nsbasic_bsd: 2E 44 4F 43 5F 49 44 20  |.DOC_ID.|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:542] nsbasic_bsd: 61 6E 64 20 4C 41 4C 2E  |and.LAL.|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:542] nsbasic_bsd: 73 74 61 6D 70 5F 64 61  |stamp_da|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:542] nsbasic_bsd: 74 65 20 3D 20 41 4C 2E  |te.=.AL.|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:542] nsbasic_bsd: 73 40 74 61 6D 70 5F 64  |s@tamp_d|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:542] nsbasic_bsd: 61 74 65 0D 0A 09 09 09  |ate.....|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:542] nsbasic_bsd: 09 09 09 29 0D 0A 09 09  |...)....|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:542] nsbasic_bsd: 09 09 09 09 73 65 6C 65  |....sele|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:542] nsbasic_bsd: 63 74 20 44 2E 44 4F 43  |ct.D.DOC|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:542] nsbasic_bsd: 5F 49 44 2C 20 4C 41 4C  |_ID,.LAL|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:542] nsbasic_bsd: 2E 53 54 41 4D 50 5F 55  |.STAMP_U|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:542] nsbasic_bsd: 49 44 20 66 72 6F 6D 20  |ID.from.|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:542] nsbasic_bsd: 6C 61 40 73 74 41 4C 5F  |la@stAL_|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:542] nsbasic_bsd: 55 49 44 20 4C 41 4C 20  |UID.LAL.|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:542] nsbasic_bsd: 6A 6F 69 6E 20 44 20 6F  |join.D.o|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:542] nsbasic_bsd: 6E 20 44 2E 44 4F 43 5F  |n.D.DOC_|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:542] nsbasic_bsd: 49 44 20 3D 20 4C 41 4C  |ID.=.LAL|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:542] nsbasic_bsd: 2E 44 4F 43 5F 49 44 0D  |.DOC_ID.|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:542] nsbasic_bsd: 0A 09 09 09 09 09 09 77  |.......w|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:542] nsbasic_bsd: 68 65 72 65 20 0D 0A 09  |here....|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:542] nsbasic_bsd: 09 09 09 40 09 09 09 4C  |...@...L|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:542] nsbasic_bsd: 41 4C 2E 73 74 61 6D 70  |AL.stamp|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:542] nsbasic_bsd: 5F 75 69 64 20 6E 6F 74  |_uid.not|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:542] nsbasic_bsd: 20 69 6E 20 28 27 49 4E  |.in.('IN|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:542] nsbasic_bsd: 56 4F 49 43 45 20 46 45  |VOICE.FE|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:542] nsbasic_bsd: 45 44 42 41 43 4B 27 2C  |EDBACK',|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:542] nsbasic_bsd: 20 27 41 55 54 4F 54 52  |.'AUTOTR|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:542] nsbasic_bsd: 41 4E 53 46 45 52 27 29  |ANSFER')|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:542] nsbasic_bsd: 0D 0A 09 09 28 09 09 09  |....(...|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:542] nsbasic_bsd: 09 09 61 6E 64 20 44 2E  |..and.D.|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:542] nsbasic_bsd: 44 4F 43 5F 49 44 20 3D  |DOC_ID.=|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:542] nsbasic_bsd: 20 3A 30 0D 0A 09 09 09  |.:0.....|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:542] nsbasic_bsd: 09 09 09 29 20 52 0D 0A  |...).R..|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:542] nsbasic_bsd: 09 09 09 09 09 00 01 01  |........|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:542] nsbasic_bsd: 00 00 00 00 00 00 01 01  |........|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:542] nsbasic_bsd: 00 00 00 00 00 01 03 00  |........|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:542] nsbasic_bsd: 00 01 80 00 01 10 00 00  |........|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:542] nsbasic_bsd: 02 03 69 01 01 20 07 20  |..i.....|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:542] nsbasic_bsd: 43 31 35 30 33 36 33 30  |C1503630|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:542] nsbasic_bsd: 34 42 33 31 34 44 33 38  |4B314D38|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:542] nsbasic_bsd: 39 45 32 32 45 46 41 39  |9E22EFA9|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:542] nsbasic_bsd: 37 34 34 39 44 41 32 36  |7449DA26|\r\n" +
					"(6632) [26-FEB-2019 11:49:09:542] nsbasic_bsd: exit (0)",
			},
			want: `
						select count(*) from (
						with D as (
							select DOC_ID,STATUS_INDEX
							from DOCS where STATUS_INDEX in (2,99)
						),last_AL as (
							select al.doc_id, max(al.stamp_date) stamp_date
							from action_log al join D on D.DOC_ID = al.DOC_ID
							group by al.doc_id
						), 
						lastAL_UID as (
							select distinct AL.doc_id,AL.stamp_uid 
							from D join last_AL LAL on D.DOC_ID = LAL.DOC_ID
							join action_log AL on LAL.DOC_ID = AL.DOC_ID and LAL.stamp_date = AL.stamp_date
						)
						select D.DOC_ID, LAL.STAMP_UID from lastAL_UID LAL join D on D.DOC_ID = LAL.DOC_ID
						where 
							LAL.stamp_uid not in ('INVOICE FEEDBACK', 'AUTOTRANSFER')
							and D.DOC_ID = :0
						) R
		`,
			wantErr: false,
		},
		{
			name: "select surrounded by ()",
			args: args{
				trc: "(7980) [01-MAR-2019 11:38:21:223] nsbasic_bsd: packet dump\r\n" +
					"(7980) [01-MAR-2019 11:38:21:223] nsbasic_bsd: 03 12 00 00 06 00 00 00  |........|\r\n" +
					"(7980) [01-MAR-2019 11:38:21:223] nsbasic_bsd: 00 00 11 69 05 01 01 01  |...i....|\r\n" +
					"(7980) [01-MAR-2019 11:38:21:223] nsbasic_bsd: 01 01 03 5E 06 02 80 69  |...^...i|\r\n" +
					"(7980) [01-MAR-2019 11:38:21:223] nsbasic_bsd: 00 01 02 06 A8 01 01 0D  |........|\r\n" +
					"(7980) [01-MAR-2019 11:38:21:223] nsbasic_bsd: 01 01 00 01 64 00 01 01  |....d...|\r\n" +
					"(7980) [01-MAR-2019 11:38:21:223] nsbasic_bsd: 06 00 01 00 01 01 01 00  |........|\r\n" +
					"(7980) [01-MAR-2019 11:38:21:223] nsbasic_bsd: 00 01 01 00 00 00 00 00  |........|\r\n" +
					"(7980) [01-MAR-2019 11:38:21:223] nsbasic_bsd: FE 40 28 53 45 4C 45 43  |.@(SELEC|\r\n" +
					"(7980) [01-MAR-2019 11:38:21:223] nsbasic_bsd: 54 20 44 49 53 54 49 4E  |T.DISTIN|\r\n" +
					"(7980) [01-MAR-2019 11:38:21:223] nsbasic_bsd: 43 54 20 74 31 30 2E 64  |CT.t10.d|\r\n" +
					"(7980) [01-MAR-2019 11:38:21:223] nsbasic_bsd: 6F 63 5F 69 64 20 46 52  |oc_id.FR|\r\n" +
					"(7980) [01-MAR-2019 11:38:21:223] nsbasic_bsd: 4F 4D 20 66 6C 6F 77 5F  |OM.flow_|\r\n" +
					"(7980) [01-MAR-2019 11:38:21:223] nsbasic_bsd: 6C 6F 67 20 74 31 30 2C  |log.t10,|\r\n" +
					"(7980) [01-MAR-2019 11:38:21:223] nsbasic_bsd: 20 66 6C 6F 77 5F 63 75  |.flow_cu|\r\n" +
					"(7980) [01-MAR-2019 11:38:21:223] nsbasic_bsd: 72 72 65 6E 74 20 74 31  |rrent.t1|\r\n" +
					"(7980) [01-MAR-2019 11:38:21:223] nsbasic_bsd: 31 20 40 20 57 48 45 52  |1.@.WHER|\r\n" +
					"(7980) [01-MAR-2019 11:38:21:223] nsbasic_bsd: 45 20 74 31 30 2E 61 63  |E.t10.ac|\r\n" +
					"(7980) [01-MAR-2019 11:38:21:223] nsbasic_bsd: 74 69 6F 6E 5F 69 6E 64  |tion_ind|\r\n" +
					"(7980) [01-MAR-2019 11:38:21:223] nsbasic_bsd: 65 78 20 3D 20 3A 31 20  |ex.=.:1.|\r\n" +
					"(7980) [01-MAR-2019 11:38:21:223] nsbasic_bsd: 41 4E 44 20 74 31 31 2E  |AND.t11.|\r\n" +
					"(7980) [01-MAR-2019 11:38:21:223] nsbasic_bsd: 72 65 63 69 70 69 65 6E  |recipien|\r\n" +
					"(7980) [01-MAR-2019 11:38:21:223] nsbasic_bsd: 74 5F 6E 61 6D 65 20 3D  |t_name.=|\r\n" +
					"(7980) [01-MAR-2019 11:38:21:223] nsbasic_bsd: 20 3A 32 20 41 4E 44 20  |.:2.AND.|\r\n" +
					"(7980) [01-MAR-2019 11:38:21:223] nsbasic_bsd: 74 31 31 40 2E 64 6F 63  |t11@.doc|\r\n" +
					"(7980) [01-MAR-2019 11:38:21:223] nsbasic_bsd: 5F 69 64 20 3D 20 74 31  |_id.=.t1|\r\n" +
					"(7980) [01-MAR-2019 11:38:21:223] nsbasic_bsd: 30 2E 64 6F 63 5F 69 64  |0.doc_id|\r\n" +
					"(7980) [01-MAR-2019 11:38:21:223] nsbasic_bsd: 20 41 4E 44 20 20 28 28  |.AND..((|\r\n" +
					"(7980) [01-MAR-2019 11:38:21:223] nsbasic_bsd: 74 31 30 2E 73 65 6E 64  |t10.send|\r\n" +
					"(7980) [01-MAR-2019 11:38:21:223] nsbasic_bsd: 65 64 5F 74 6F 5F 74 69  |ed_to_ti|\r\n" +
					"(7980) [01-MAR-2019 11:38:21:223] nsbasic_bsd: 6D 65 73 74 61 6D 70 20  |mestamp.|\r\n" +
					"(7980) [01-MAR-2019 11:38:21:223] nsbasic_bsd: 3E 20 28 53 45 4C 45 43  |>.(SELEC|\r\n" +
					"(7980) [01-MAR-2019 11:38:21:223] nsbasic_bsd: 54 20 4D 41 40 58 28 74  |T.MA@X(t|\r\n" +
					"(7980) [01-MAR-2019 11:38:21:223] nsbasic_bsd: 32 30 2E 73 65 6E 64 65  |20.sende|\r\n" +
					"(7980) [01-MAR-2019 11:38:21:223] nsbasic_bsd: 64 5F 74 6F 5F 74 69 6D  |d_to_tim|\r\n" +
					"(7980) [01-MAR-2019 11:38:21:223] nsbasic_bsd: 65 73 74 61 6D 70 29 20  |estamp).|\r\n" +
					"(7980) [01-MAR-2019 11:38:21:223] nsbasic_bsd: 46 52 4F 4D 20 66 6C 6F  |FROM.flo|\r\n" +
					"(7980) [01-MAR-2019 11:38:21:223] nsbasic_bsd: 77 5F 6C 6F 67 20 74 32  |w_log.t2|\r\n" +
					"(7980) [01-MAR-2019 11:38:21:223] nsbasic_bsd: 30 2C 20 66 6C 6F 77 5F  |0,.flow_|\r\n" +
					"(7980) [01-MAR-2019 11:38:21:223] nsbasic_bsd: 63 75 72 72 65 6E 74 20  |current.|\r\n" +
					"(7980) [01-MAR-2019 11:38:21:223] nsbasic_bsd: 74 32 31 20 20 40 57 48  |t21..@WH|\r\n" +
					"(7980) [01-MAR-2019 11:38:21:223] nsbasic_bsd: 45 52 45 20 74 32 30 2E  |ERE.t20.|\r\n" +
					"(7980) [01-MAR-2019 11:38:21:223] nsbasic_bsd: 61 63 74 69 6F 6E 5F 69  |action_i|\r\n" +
					"(7980) [01-MAR-2019 11:38:21:223] nsbasic_bsd: 6E 64 65 78 20 3D 20 3A  |ndex.=.:|\r\n" +
					"(7980) [01-MAR-2019 11:38:21:223] nsbasic_bsd: 33 20 41 4E 44 20 74 32  |3.AND.t2|\r\n" +
					"(7980) [01-MAR-2019 11:38:21:223] nsbasic_bsd: 31 2E 72 65 63 69 70 69  |1.recipi|\r\n" +
					"(7980) [01-MAR-2019 11:38:21:223] nsbasic_bsd: 65 6E 74 5F 6E 61 6D 65  |ent_name|\r\n" +
					"(7980) [01-MAR-2019 11:38:21:223] nsbasic_bsd: 20 3D 20 3A 34 20 41 4E  |.=.:4.AN|\r\n" +
					"(7980) [01-MAR-2019 11:38:21:223] nsbasic_bsd: 44 20 74 32 30 2E 40 64  |D.t20.@d|\r\n" +
					"(7980) [01-MAR-2019 11:38:21:223] nsbasic_bsd: 6F 63 5F 69 64 20 3D 20  |oc_id.=.|\r\n" +
					"(7980) [01-MAR-2019 11:38:21:223] nsbasic_bsd: 74 31 30 2E 64 6F 63 5F  |t10.doc_|\r\n" +
					"(7980) [01-MAR-2019 11:38:21:223] nsbasic_bsd: 69 64 20 41 4E 44 20 74  |id.AND.t|\r\n" +
					"(7980) [01-MAR-2019 11:38:21:223] nsbasic_bsd: 32 31 2E 64 6F 63 5F 69  |21.doc_i|\r\n" +
					"(7980) [01-MAR-2019 11:38:21:223] nsbasic_bsd: 64 20 3D 20 74 31 30 2E  |d.=.t10.|\r\n" +
					"(7980) [01-MAR-2019 11:38:21:223] nsbasic_bsd: 64 6F 63 5F 69 64 29 29  |doc_id))|\r\n" +
					"(7980) [01-MAR-2019 11:38:21:223] nsbasic_bsd: 20 20 4F 52 20 28 53 45  |..OR.(SE|\r\n" +
					"(7980) [01-MAR-2019 11:38:21:223] nsbasic_bsd: 4C 45 43 54 20 63 6F 40  |LECT.co@|\r\n" +
					"(7980) [01-MAR-2019 11:38:21:223] nsbasic_bsd: 75 6E 74 28 74 32 30 2E  |unt(t20.|\r\n" +
					"(7980) [01-MAR-2019 11:38:21:223] nsbasic_bsd: 73 65 6E 64 65 64 5F 74  |sended_t|\r\n" +
					"(7980) [01-MAR-2019 11:38:21:223] nsbasic_bsd: 6F 5F 74 69 6D 65 73 74  |o_timest|\r\n" +
					"(7980) [01-MAR-2019 11:38:21:223] nsbasic_bsd: 61 6D 70 29 20 46 52 4F  |amp).FRO|\r\n" +
					"(7980) [01-MAR-2019 11:38:21:223] nsbasic_bsd: 4D 20 66 6C 6F 77 5F 6C  |M.flow_l|\r\n" +
					"(7980) [01-MAR-2019 11:38:21:223] nsbasic_bsd: 6F 67 20 74 32 30 2C 20  |og.t20,.|\r\n" +
					"(7980) [01-MAR-2019 11:38:21:223] nsbasic_bsd: 66 6C 6F 77 5F 63 75 72  |flow_cur|\r\n" +
					"(7980) [01-MAR-2019 11:38:21:223] nsbasic_bsd: 72 65 6E 74 20 74 32 31  |rent.t21|\r\n" +
					"(7980) [01-MAR-2019 11:38:21:223] nsbasic_bsd: 40 20 57 48 45 52 45 20  |@.WHERE.|\r\n" +
					"(7980) [01-MAR-2019 11:38:21:223] nsbasic_bsd: 20 74 32 30 2E 61 63 74  |.t20.act|\r\n" +
					"(7980) [01-MAR-2019 11:38:21:223] nsbasic_bsd: 69 6F 6E 5F 69 6E 64 65  |ion_inde|\r\n" +
					"(7980) [01-MAR-2019 11:38:21:223] nsbasic_bsd: 78 20 3D 20 3A 35 20 41  |x.=.:5.A|\r\n" +
					"(7980) [01-MAR-2019 11:38:21:223] nsbasic_bsd: 4E 44 20 74 32 31 2E 72  |ND.t21.r|\r\n" +
					"(7980) [01-MAR-2019 11:38:21:223] nsbasic_bsd: 65 63 69 70 69 65 6E 74  |ecipient|\r\n" +
					"(7980) [01-MAR-2019 11:38:21:223] nsbasic_bsd: 5F 6E 61 6D 65 20 3D 20  |_name.=.|\r\n" +
					"(7980) [01-MAR-2019 11:38:21:224] nsbasic_bsd: 3A 36 20 41 4E 44 20 74  |:6.AND.t|\r\n" +
					"(7980) [01-MAR-2019 11:38:21:224] nsbasic_bsd: 32 38 30 2E 64 6F 63 5F  |280.doc_|\r\n" +
					"(7980) [01-MAR-2019 11:38:21:224] nsbasic_bsd: 69 64 20 3D 20 74 31 30  |id.=.t10|\r\n" +
					"(7980) [01-MAR-2019 11:38:21:224] nsbasic_bsd: 2E 64 6F 63 5F 69 64 20  |.doc_id.|\r\n" +
					"(7980) [01-MAR-2019 11:38:21:224] nsbasic_bsd: 41 4E 44 20 74 32 31 2E  |AND.t21.|\r\n" +
					"(7980) [01-MAR-2019 11:38:21:224] nsbasic_bsd: 64 6F 63 5F 69 64 20 3D  |doc_id.=|\r\n" +
					"(7980) [01-MAR-2019 11:38:21:224] nsbasic_bsd: 20 74 31 30 2E 64 6F 63  |.t10.doc|\r\n" +
					"(7980) [01-MAR-2019 11:38:21:224] nsbasic_bsd: 5F 69 64 29 20 3C 20 31  |_id).<.1|\r\n" +
					"(7980) [01-MAR-2019 11:38:21:224] nsbasic_bsd: 29 29 00 01 01 00 00 00  |))......|\r\n" +
					"(7980) [01-MAR-2019 11:38:21:224] nsbasic_bsd: 00 00 00 01 01 00 00 00  |........|\r\n" +
					"(7980) [01-MAR-2019 11:38:21:224] nsbasic_bsd: exit (0)\r\n" +
					"(7980) [01-MAR-2019 11:38:21:224] nsbasic_brc: entry: oln/tot=0\r\n" +
					"(7980) [01-MAR-2019 11:38:21:224] nttfprd: entry",
			},
			want:    "(SELECT DISTINCT t10.doc_id FROM flow_log t10, flow_current t11  WHERE t10.action_index = :1 AND t11.recipient_name = :2 AND t11.doc_id = t10.doc_id AND  ((t10.sended_to_timestamp > (SELECT MAX(t20.sended_to_timestamp) FROM flow_log t20, flow_current t21  WHERE t20.action_index = :3 AND t21.recipient_name = :4 AND t20.doc_id = t10.doc_id AND t21.doc_id = t10.doc_id))  OR (SELECT count(t20.sended_to_timestamp) FROM flow_log t20, flow_current t21 WHERE  t20.action_index = :5 AND t21.recipient_name = :6 AND t20.doc_id = t10.doc_id AND t21.doc_id = t10.doc_id) < 1))",
			wantErr: false,
		},
		{
			name: "select having width in it",
			args: args{
				trc: "(7980) [01-MAR-2019 11:37:47:910] nsbasic_bsd: packet dump\r\n" +
					"(7980) [01-MAR-2019 11:37:47:910] nsbasic_bsd: 01 BC 00 00 06 00 00 00  |........|\r\n" +
					"(7980) [01-MAR-2019 11:37:47:910] nsbasic_bsd: 00 00 11 69 15 01 01 01  |...i....|\r\n" +
					"(7980) [01-MAR-2019 11:37:47:910] nsbasic_bsd: 01 01 03 5E 16 02 80 69  |...^...i|\r\n" +
					"(7980) [01-MAR-2019 11:37:47:910] nsbasic_bsd: 00 01 02 03 3F 01 01 0D  |....?...|\r\n" +
					"(7980) [01-MAR-2019 11:37:47:910] nsbasic_bsd: 01 01 00 01 64 00 01 01  |....d...|\r\n" +
					"(7980) [01-MAR-2019 11:37:47:910] nsbasic_bsd: 03 00 01 00 01 01 01 00  |........|\r\n" +
					"(7980) [01-MAR-2019 11:37:47:910] nsbasic_bsd: 00 01 01 00 00 00 00 00  |........|\r\n" +
					"(7980) [01-MAR-2019 11:37:47:910] nsbasic_bsd: FE 40 73 65 6C 65 63 74  |.@select|\r\n" +
					"(7980) [01-MAR-2019 11:37:47:910] nsbasic_bsd: 20 66 69 65 6C 64 5F 6E  |.field_n|\r\n" +
					"(7980) [01-MAR-2019 11:37:47:910] nsbasic_bsd: 61 6D 65 2C 20 74 61 62  |ame,.tab|\r\n" +
					"(7980) [01-MAR-2019 11:37:47:910] nsbasic_bsd: 6C 65 5F 6E 61 6D 65 2C  |le_name,|\r\n" +
					"(7980) [01-MAR-2019 11:37:47:910] nsbasic_bsd: 20 66 69 65 6C 64 5F 77  |.field_w|\r\n" +
					"(7980) [01-MAR-2019 11:37:47:910] nsbasic_bsd: 69 64 74 68 2C 20 66 69  |idth,.fi|\r\n" +
					"(7980) [01-MAR-2019 11:37:47:910] nsbasic_bsd: 65 6C 64 5F 73 6F 72 74  |eld_sort|\r\n" +
					"(7980) [01-MAR-2019 11:37:47:910] nsbasic_bsd: 2C 20 73 70 65 63 69 61  |,.specia|\r\n" +
					"(7980) [01-MAR-2019 11:37:47:910] nsbasic_bsd: 6C 5F 40 66 69 65 6C 64  |l_@field|\r\n" +
					"(7980) [01-MAR-2019 11:37:47:910] nsbasic_bsd: 20 66 72 6F 6D 20 72 65  |.from.re|\r\n" +
					"(7980) [01-MAR-2019 11:37:47:910] nsbasic_bsd: 70 6F 72 74 5F 68 65 61  |port_hea|\r\n" +
					"(7980) [01-MAR-2019 11:37:47:910] nsbasic_bsd: 64 65 72 73 32 20 72 31  |ders2.r1|\r\n" +
					"(7980) [01-MAR-2019 11:37:47:910] nsbasic_bsd: 20 77 68 65 72 65 20 55  |.where.U|\r\n" +
					"(7980) [01-MAR-2019 11:37:47:910] nsbasic_bsd: 50 50 45 52 28 72 31 2E  |PPER(r1.|\r\n" +
					"(7980) [01-MAR-2019 11:37:47:910] nsbasic_bsd: 73 65 61 72 63 68 5F 69  |search_i|\r\n" +
					"(7980) [01-MAR-2019 11:37:47:910] nsbasic_bsd: 64 29 20 3D 20 20 3A 31  |d).=..:1|\r\n" +
					"(7980) [01-MAR-2019 11:37:47:910] nsbasic_bsd: 20 61 6E 40 64 20 72 31  |.an@d.r1|\r\n" +
					"(7980) [01-MAR-2019 11:37:47:910] nsbasic_bsd: 2E 75 73 65 72 5F 69 64  |.user_id|\r\n" +
					"(7980) [01-MAR-2019 11:37:47:910] nsbasic_bsd: 3D 27 44 45 46 41 55 4C  |='DEFAUL|\r\n" +
					"(7980) [01-MAR-2019 11:37:47:910] nsbasic_bsd: 54 27 20 61 6E 64 20 74  |T'.and.t|\r\n" +
					"(7980) [01-MAR-2019 11:37:47:910] nsbasic_bsd: 72 69 6D 28 66 69 65 6C  |rim(fiel|\r\n" +
					"(7980) [01-MAR-2019 11:37:47:910] nsbasic_bsd: 64 5F 6E 61 6D 65 29 20  |d_name).|\r\n" +
					"(7980) [01-MAR-2019 11:37:47:910] nsbasic_bsd: 20 6E 6F 74 20 69 6E 20  |.not.in.|\r\n" +
					"(7980) [01-MAR-2019 11:37:47:910] nsbasic_bsd: 28 73 65 6C 65 63 74 20  |(select.|\r\n" +
					"(7980) [01-MAR-2019 11:37:47:910] nsbasic_bsd: 74 72 69 6D 40 28 66 69  |trim@(fi|\r\n" +
					"(7980) [01-MAR-2019 11:37:47:910] nsbasic_bsd: 65 6C 64 5F 6E 61 6D 65  |eld_name|\r\n" +
					"(7980) [01-MAR-2019 11:37:47:910] nsbasic_bsd: 29 20 66 72 6F 6D 20 72  |).from.r|\r\n" +
					"(7980) [01-MAR-2019 11:37:47:910] nsbasic_bsd: 65 70 6F 72 74 5F 68 65  |eport_he|\r\n" +
					"(7980) [01-MAR-2019 11:37:47:910] nsbasic_bsd: 61 64 65 72 73 32 20 72  |aders2.r|\r\n" +
					"(7980) [01-MAR-2019 11:37:47:910] nsbasic_bsd: 32 20 77 68 65 72 65 20  |2.where.|\r\n" +
					"(7980) [01-MAR-2019 11:37:47:910] nsbasic_bsd: 55 50 50 45 52 28 72 32  |UPPER(r2|\r\n" +
					"(7980) [01-MAR-2019 11:37:47:910] nsbasic_bsd: 2E 73 65 61 72 63 68 5F  |.search_|\r\n" +
					"(7980) [01-MAR-2019 11:37:47:910] nsbasic_bsd: 69 64 29 3D 3A 15 32 20  |id)=:.2.|\r\n" +
					"(7980) [01-MAR-2019 11:37:47:910] nsbasic_bsd: 20 61 6E 64 20 72 32 2E  |.and.r2.|\r\n" +
					"(7980) [01-MAR-2019 11:37:47:910] nsbasic_bsd: 75 73 65 72 5F 69 64 3D  |user_id=|\r\n" +
					"(7980) [01-MAR-2019 11:37:47:910] nsbasic_bsd: 3A 33 29 00 01 01 00 00  |:3).....|\r\n" +
					"(7980) [01-MAR-2019 11:37:47:910] nsbasic_bsd: exit (0)",
			},
			want:    "select field_name, table_name, field_width, field_sort, special_field from report_headers2 r1 where UPPER(r1.search_id) =  :1 and r1.user_id='DEFAULT' and trim(field_name)  not in (select trim(field_name) from report_headers2 r2 where UPPER(r2.search_id)=:2  and r2.user_id=:3)",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getQueryFromTraceSnippet(tt.args.trc)
			if (err != nil) != tt.wantErr {
				t.Errorf("getQueryFromTraceSnippet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !significantCharsEqual(got, tt.want) {
				t.Errorf("getQueryFromTraceSnippet() = \n%v,\n want \n%v", got, tt.want)
			}
		})
	}
}
