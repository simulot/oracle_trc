package ts

import (
	"reflect"
	"testing"
	"time"
)

func TestOracleTS_DD_MON_YYYY_HH_MI_SS_FF3(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		args    args
		want    time.Time
		wantErr bool
	}{
		{
			name:    "26-FEB-2019 11:46:40:939",
			args:    args{[]byte("26-FEB-2019 11:46:40:939")},
			want:    time.Date(2019, 02, 26, 11, 46, 40, 939, time.Local),
			wantErr: false,
		},
		{
			name:    "26-Feb-2019 11:46:40:939",
			args:    args{[]byte("26-Feb-2019 11:46:40:939")},
			want:    time.Date(2019, 02, 26, 11, 46, 40, 939, time.Local),
			wantErr: false,
		},
		{
			name:    "26-AAA-2019 11:46:40:939",
			args:    args{[]byte("26-AAA-2019 11:46:40:939")},
			want:    time.Time{},
			wantErr: true,
		},
		{
			name:    "AE-FEB-2019 11:46:40:939",
			args:    args{[]byte("AE-FEB-2019 11:46:40:939")},
			want:    time.Time{},
			wantErr: true,
		},
		{
			name:    "26-FEB-2019",
			args:    args{[]byte("26-FEB-2019")},
			want:    time.Date(2019, 02, 26, 0, 0, 0, 0, time.Local),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := OracleTS_DD_MON_YYYY_HH_MI_SS_FF9(tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("OracleTS_DD_MON_YYYY_HH_MI_SS_FF3() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OracleTS_DD_MON_YYYY_HH_MI_SS_FF3() = %v, want %v", got, tt.want)
			}
		})
	}
}
