package ts

import (
	"bytes"
	"time"

	"github.com/pkg/errors"
)

/*
  This package converts oracle times format
*/

type TimeParserFn func([]byte) (time.Time, error)

func GetParser(layout string) (TimeParserFn, error) {
	f, ok := parser[layout]
	if !ok {
		return nil, errors.Errorf("Unsupported time format :'%s'", layout)
	}
	return f, nil
}

var parser = map[string]TimeParserFn{
	"DD-MON-YYYY HH:MI:SS:FF1": OracleTS_DD_MON_YYYY_HH_MI_SS_FF9,
	"DD-MON-YYYY HH:MI:SS:FF2": OracleTS_DD_MON_YYYY_HH_MI_SS_FF9,
	"DD-MON-YYYY HH:MI:SS:FF3": OracleTS_DD_MON_YYYY_HH_MI_SS_FF9,
	"DD-MON-YYYY HH:MI:SS:FF4": OracleTS_DD_MON_YYYY_HH_MI_SS_FF9,
	"DD-MON-YYYY HH:MI:SS:FF5": OracleTS_DD_MON_YYYY_HH_MI_SS_FF9,
	"DD-MON-YYYY HH:MI:SS:FF6": OracleTS_DD_MON_YYYY_HH_MI_SS_FF9,
	"DD-MON-YYYY HH:MI:SS:FF7": OracleTS_DD_MON_YYYY_HH_MI_SS_FF9,
	"DD-MON-YYYY HH:MI:SS:FF8": OracleTS_DD_MON_YYYY_HH_MI_SS_FF9,
	"DD-MON-YYYY HH:MI:SS:FF9": OracleTS_DD_MON_YYYY_HH_MI_SS_FF9,
}

func OracleTS_DD_MON_YYYY_HH_MI_SS_FF9(b []byte) (time.Time, error) {
	day, month, year, hour, minute, second, millisecond := 0, 0, 0, 0, 0, 0, 0

	if len(b) >= 2 {
		day, b = eatDigits(b, 2)
		if day < 0 {
			return time.Time{}, errors.New("Can't parse timestamp's day")
		}
	}
	if len(b) > 1 {
		b = eatSeparators(b)
	}
	if len(b) >= 3 {
		var m []byte
		m, b = eatLetters(b, 3)
		for i := 0; i < len(monthShort); i++ {
			if bytes.Equal(m, monthShort[i]) {
				month = i + 1
				break
			}
		}
		if month == 0 {
			return time.Time{}, errors.New("Can't parse timestamp's month")
		}
	}
	if len(b) > 1 {
		b = eatSeparators(b)
	}
	if len(b) >= 4 {
		year, b = eatDigits(b, 4)
		if year < 0 {
			return time.Time{}, errors.New("Can't parse timestamp's year")
		}
	}
	if len(b) > 1 {
		b = eatSeparators(b)
	}

	if len(b) >= 2 {
		hour, b = eatDigits(b, 2)
		if hour < 0 {
			return time.Time{}, errors.New("Can't parse timestamp's hour")
		}
	}
	if len(b) > 1 {
		b = eatSeparators(b)
	}
	if len(b) >= 2 {
		minute, b = eatDigits(b, 2)
		if minute < 0 {
			return time.Time{}, errors.New("Can't parse timestamp's minute")
		}
	}
	if len(b) > 1 {
		b = eatSeparators(b)
	}
	if len(b) >= 2 {
		second, b = eatDigits(b, 2)
		if second < 0 {
			return time.Time{}, errors.New("Can't parse timestamp's second")
		}
	}
	if len(b) > 1 {
		b = eatSeparators(b)
	}
	if len(b) >= 1 {
		millisecond, b = eatDigits(b, 9)
		if millisecond < 0 {
			return time.Time{}, errors.New("Can't parse timestamp's millisecond")
		}
	}
	return time.Date(year, time.Month(month), day, hour, minute, second, millisecond, time.Local), nil
}

var monthShort = [12][]byte{
	[]byte("jan"),
	[]byte("feb"),
	[]byte("mar"),
	[]byte("apr"),
	[]byte("may"),
	[]byte("jun"),
	[]byte("jul"),
	[]byte("aug"),
	[]byte("sep"),
	[]byte("oct"),
	[]byte("nov"),
	[]byte("dec"),
}

func eatDigits(b []byte, maxlen int) (int, []byte) {
	n := 0
	i := 0
	for i = 0; i < len(b); i++ {
		c := b[i]
		if c < '0' || c > '9' {
			break
		}
		n = (n * 10) + int(c) - '0'
	}
	if i == 0 {
		return -1, b

	}
	return n, b[i:]
}
func eatLetters(b []byte, maxlen int) (m []byte, suffix []byte) {
	m = make([]byte, maxlen)
	i := 0
	for i = 0; i < len(b); i++ {
		c := b[i] | ('a' - 'A')
		if c < 'a' || c > 'z' {
			break
		}
		m[i] = c
	}
	return m, b[i:]
}

func eatSeparators(b []byte) []byte {
	i := 0
	for i < len(b) {
		switch b[i] {
		case ':', '-', '/', '.', ' ':
			i++
		default:
			return b[i:]
		}
	}
	return []byte{}

}
