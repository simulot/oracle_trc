package queries

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"strconv"
	"time"
)

func GetInt64(r *bytes.Buffer, size int, compress bool, bigEndian bool) (int64, error) {
	var ret int64

	negFlag := false
	if compress {
		b, err := r.ReadByte()
		if err != nil {
			return 0, err
		}
		if (b & 0x80) != 0 {
			negFlag = true
		}
		size = int(b & 0x7F)
		bigEndian = true
	}
	if size == 0 {
		return 0, nil
	}
	buff := r.Next(size)
	temp := make([]byte, 8)
	if bigEndian {
		copy(temp[8-size:], buff)
		ret = int64(binary.BigEndian.Uint64(temp))
	} else {
		copy(temp[:size], buff)
		ret = int64(binary.LittleEndian.Uint64(temp))
	}
	if negFlag {
		ret = ret * -1
	}
	return ret, nil
}

func GetInt(r *bytes.Buffer, size int, compress bool, bigEndian bool) (int32, error) {
	temp, err := GetInt64(r, size, compress, bigEndian)
	if err != nil {
		return 0, err
	}
	return int32(temp), nil
}

func GetUInt(r *bytes.Buffer, size int, compress bool, bigEndian bool) (uint32, error) {
	i, err := GetInt(r, size, compress, bigEndian)
	return uint32(i), err
}

func readBytes(buff *bytes.Buffer) ([]byte, error) {
	out := make([]byte, 0, 40)
	var l, b byte
	var err error
	hasChunks := false

	l, err = buff.ReadByte()
	if err != nil {
		return nil, err
	}
	if l == 0xFE {
		// Marker for buffer bigger than 0x40
		l, err = buff.ReadByte()
		if err != nil {
			return nil, err
		}
		hasChunks = true
	}

extern:
	for l > 0 {
		for l > 0 {
			l--
			b, err = buff.ReadByte()
			if err != nil {
				break extern
			}
			out = append(out, b)
		}
		if hasChunks {
			l, err = buff.ReadByte()
			if err != nil {
				break extern
			}
		}
	}
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Information found on GO-ORA project by Samy Sultan

type OracleType uint
type ParameterDirection uint

const (
	Input  ParameterDirection = 1
	Output ParameterDirection = 2
	InOut  ParameterDirection = 3
	RetVal ParameterDirection = 9
)

//go:generate stringer -type=OracleType

const (
	NCHAR            OracleType = 1
	NUMBER           OracleType = 2
	SB1              OracleType = 3
	SB2              OracleType = 3
	SB4              OracleType = 3
	FLOAT            OracleType = 4
	NullStr          OracleType = 5
	VarNum           OracleType = 6
	LONG             OracleType = 8
	VARCHAR          OracleType = 9
	ROWID            OracleType = 11
	DATE             OracleType = 12
	VarRaw           OracleType = 15
	BFloat           OracleType = 21
	BDouble          OracleType = 22
	RAW              OracleType = 23
	LongRaw          OracleType = 24
	UINT             OracleType = 68
	LongVarChar      OracleType = 94
	LongVarRaw       OracleType = 95
	CHAR             OracleType = 96
	CHARZ            OracleType = 97
	IBFloat          OracleType = 100
	IBDouble         OracleType = 101
	RefCursor        OracleType = 102
	NOT              OracleType = 108
	XMLType          OracleType = 108
	OCIRef           OracleType = 110
	OCIClobLocator   OracleType = 112
	OCIBlobLocator   OracleType = 113
	OCIFileLocator   OracleType = 114
	ResultSet        OracleType = 116
	OCIString        OracleType = 155
	OCIDate          OracleType = 156
	TimeStampDTY     OracleType = 180
	TimeStampTZ_DTY  OracleType = 181
	IntervalYM_DTY   OracleType = 182
	IntervalDS_DTY   OracleType = 183
	TimeTZ           OracleType = 186
	TimeStamp        OracleType = 187
	TimeStampTZ      OracleType = 188
	IntervalYM       OracleType = 189
	IntervalDS       OracleType = 190
	UROWID           OracleType = 208
	TimeStampLTZ_DTY OracleType = 231
	TimeStampeLTZ    OracleType = 232
)

type ParameterType int

const (
	Number ParameterType = 1
	String ParameterType = 2
)

type ParameterInfo struct {
	Name                 string
	Direction            ParameterDirection
	IsNull               bool
	AllowNull            bool
	ColAlias             string
	DataType             OracleType
	IsXmlType            bool
	Flag                 uint8
	Precision            uint8
	Scale                uint8
	MaxLen               uint32
	MaxCharLen           uint32
	MaxNoOfArrayElements uint32
	ContFlag             uint32
	ToID                 []byte
	Version              uint32
	CharsetID            uint32
	CharsetForm          uint8
	Value                []byte
	getDataFromServer    bool
}

func (p ParameterInfo) String() string {
	if p.IsNull {
		return "(null)"
	}
	switch p.DataType {
	case CHAR:
		return string(p.Value)
	case DATE, TimeStamp, TimeStampDTY, TimeStampeLTZ, TimeStampLTZ_DTY, TimeStampTZ, TimeStampTZ_DTY:
		d, err := DecodeDate(p.Value)
		if err != nil {
			return "( " + err.Error() + ")"
		}
		return d.Format(time.RFC3339)
	case NUMBER:
		return strconv.FormatFloat(DecodeDouble(p.Value), 'g', -1, 64)
	default:
		return "(" + p.DataType.String() + ")"
	}
}

func GetParamInfo(buff *bytes.Buffer) (*ParameterInfo, error) {
	p := &ParameterInfo{}
	var err error
	var b byte
	b, err = buff.ReadByte()
	if err != nil {
		return nil, err
	}
	p.DataType = OracleType(b)

	p.Flag, err = buff.ReadByte()
	if err != nil {
		return nil, err
	}

	p.Precision, err = buff.ReadByte()
	if err != nil {
		return nil, err
	}

	p.Scale, err = buff.ReadByte()
	if err != nil {
		return nil, err
	}

	p.MaxLen, err = GetUInt(buff, 4, true, true)
	if err != nil {
		return nil, err
	}
	p.MaxNoOfArrayElements, err = GetUInt(buff, 4, true, true)
	if err != nil {
		return nil, err
	}

	p.ContFlag, err = GetUInt(buff, 4, true, true)
	if err != nil {
		return nil, err
	}

	//  ToID is present ?
	b, err = buff.ReadByte()
	if err != nil {
		return nil, err
	}
	if b > 0 {
		var blen uint32
		blen, err = GetUInt(buff, 4, true, true)
		if err != nil {
			return nil, err
		}
		p.ToID, err = readBytes(buff)
		if err != nil {
			return nil, err
		}
		if len(p.ToID) != int(blen) {
			if err != nil {
				return nil, fmt.Errorf("GetParamInfo unexpected len of ToID, got %d, expected %d", len(p.ToID), int(blen))
			}
		}
	}

	p.Version, err = GetUInt(buff, 2, true, true)
	if err != nil {
		return nil, err
	}

	p.CharsetID, err = GetUInt(buff, 2, true, true)
	if err != nil {
		return nil, err
	}

	p.CharsetForm, err = buff.ReadByte()
	if err != nil {
		return nil, err
	}

	p.MaxCharLen, err = GetUInt(buff, 2, true, true)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func DecodeDate(data []byte) (time.Time, error) {
	if len(data) < 7 {
		return time.Now(), errors.New("abnormal data representation for date")
	}
	year := (int(data[0]) - 100) * 100
	year += int(data[1]) - 100
	nanoSec := 0
	if len(data) > 7 {
		nanoSec = int(binary.BigEndian.Uint32(data[7:10]))
	}
	tzHour := 0
	tzMin := 0
	if len(data) > 11 {
		tzHour = int(data[11]) - 20
		tzMin = int(data[12]) - 60
	}

	return time.Date(year, time.Month(data[2]), int(data[3]),
		int(data[4]-1)+tzHour, int(data[5]-1)+tzMin, int(data[6]-1), nanoSec, time.UTC), nil
}

func DecodeInt(inputData []byte) int {
	// take a copy of input
	input := make([]byte, len(inputData))
	copy(input, inputData)
	if input[0] == 0x80 {
		return 0
	}
	length, neg := decodeSign(input)
	if length > len(input[1:]) {
		input = append(input, make([]byte, length-len(input[1:]))...)
	}
	data := input[1 : 1+length]
	ret := 0
	for x := 0; x < len(data); x++ {
		ret = (ret * 100) + int(data[x])
	}
	if neg {
		return ret * -1
	} else {
		return ret
	}
}

func decodeSign(input []byte) (length int, neg bool) {
	if input[0] > 0x80 {
		length = int(input[0]) - 0x80 - 0x40
		for x := 1; x < len(input); x++ {
			input[x] = input[x] - 1
		}
		neg = false
	} else {
		length = 0xFF - int(input[0]) - 0x80 - 0x40
		if len(input) <= 20 && input[len(input)-1] == 102 {
			input = input[:len(input)-1]
		}
		for x := 1; x < len(input); x++ {
			input[x] = uint8(101 - input[x])
		}
		neg = true
	}
	return
}

// ProtectAddFigure check if adding digit d overflows the int64 capacity.
// Return true when overflow
func ProtectAddFigure(m *int64, d int64) bool {
	r := *m * 10
	if r < 0 {
		return true
	}
	r += d
	if r < 0 {
		return true
	}
	*m = r
	return false
}

// DecodeDouble decode Oracle binary representation of numbers into float64
//
// Some documentation:
//	https://gotodba.com/2015/03/24/how-are-numbers-saved-in-oracle/
//  https://www.orafaq.com/wiki/Number

func DecodeDouble(inputData []byte) float64 {

	if len(inputData) == 0 {
		return math.NaN()
	}
	if inputData[0] == 0x80 {
		return 0
	}
	var (
		negative bool
		exponent int
		mantissa int64
	)

	negative = inputData[0]&0x80 == 0
	if negative {
		exponent = int(inputData[0]^0x7f) - 64
	} else {
		exponent = int(inputData[0]&0x7f) - 64
	}

	buf := inputData[1:]
	// When negative, strip the last byte if equal 0x66
	if negative && inputData[len(inputData)-1] == 0x66 {
		buf = inputData[1 : len(inputData)-1]
	}

	// Loop on mantissa digits, stop with the capacity of int64 is reached
	mantissaDigits := 0
	for _, digit100 := range buf {
		digit100--
		if negative {
			digit100 = 100 - digit100
		}
		if ProtectAddFigure(&mantissa, int64(digit100/10)) {
			break
		}
		mantissaDigits++
		if ProtectAddFigure(&mantissa, int64(digit100%10)) {
			break
		}
		mantissaDigits++
	}

	exponent = exponent*2 - mantissaDigits // Adjust exponent to the retrieved mantissa
	if negative {
		mantissa = -mantissa
	}

	ret := float64(mantissa) * math.Pow10(exponent)
	return ret
}
