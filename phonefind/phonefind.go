package phonefind

import (
	"bytes"
	"errors"
	"io/ioutil"
	"path"
	"runtime"

	"github.com/kamasamikon/miego/conf"
)

const (
	CMCC               byte = iota + 0x01 //中国移动
	CUCC                                  //中国联通
	CTCC                                  //中国电信
	CTCC_v                                //电信虚拟运营商
	CUCC_v                                //联通虚拟运营商
	CMCC_v                                //移动虚拟运营商
	INT_LEN            = 4
	CHAR_LEN           = 1
	HEAD_LENGTH        = 8
	PHONE_INDEX_LENGTH = 9
	PHONE_DAT          = "phone.dat"
)

type PhoneRecord struct {
	PhoneNum string
	Province string
	City     string
	ZipCode  string
	AreaZone string
	CardType string
}

var (
	content     []byte
	CardTypemap = map[byte]string{
		CMCC:   "中国移动",
		CUCC:   "中国联通",
		CTCC:   "中国电信",
		CTCC_v: "中国电信虚拟运营商",
		CUCC_v: "中国联通虚拟运营商",
		CMCC_v: "中国移动虚拟运营商",
	}
	total_len, firstoffset int32
)

func Init(phoneData string) (int32, error) {
	if phoneData == "" {
		_, fulleFilename, _, _ := runtime.Caller(0)
		phoneData = path.Join(path.Dir(fulleFilename), PHONE_DAT)
	}
	var err error
	content, err = ioutil.ReadFile(phoneData)
	if err != nil {
		return 0, err
	}
	total_len = int32(len(content))
	firstoffset = get4(content[INT_LEN : INT_LEN*2])

	return totalRecord(), nil
}

func get4(b []byte) int32 {
	if len(b) < 4 {
		return 0
	}
	return int32(b[0]) | int32(b[1])<<8 | int32(b[2])<<16 | int32(b[3])<<24
}

func getN(s string) (uint32, error) {
	var n, cutoff, maxVal uint32
	i := 0
	base := 10
	cutoff = (1<<32-1)/10 + 1
	maxVal = 1<<uint(32) - 1
	for ; i < len(s); i++ {
		var v byte
		d := s[i]
		switch {
		case '0' <= d && d <= '9':
			v = d - '0'
		case 'a' <= d && d <= 'z':
			v = d - 'a' + 10
		case 'A' <= d && d <= 'Z':
			v = d - 'A' + 10
		default:
			return 0, errors.New("invalid syntax")
		}
		if v >= byte(base) {
			return 0, errors.New("invalid syntax")
		}

		if n >= cutoff {
			// n*base overflows
			n = (1<<32 - 1)
			return n, errors.New("value out of range")
		}
		n *= uint32(base)

		n1 := n + uint32(v)
		if n1 < n || n1 > maxVal {
			// n+v overflows
			n = (1<<32 - 1)
			return n, errors.New("value out of range")
		}
		n = n1
	}
	return n, nil
}

func version() string {
	return string(content[0:INT_LEN])
}

func totalRecord() int32 {
	return (int32(len(content)) - firstRecordOffset()) / PHONE_INDEX_LENGTH
}

func firstRecordOffset() int32 {
	return get4(content[INT_LEN : INT_LEN*2])
}

// 二分法查询phone数据
func Find(phone_num string, checkMode bool) (pr *PhoneRecord, err error) {
	if len(phone_num) < 7 || len(phone_num) > 11 {
		return nil, errors.New("illegal phone length")
	}

	var left int32
	phone_seven_int, err := getN(phone_num[0:7])
	if err != nil {
		return nil, errors.New("illegal phone number")
	}
	phone_seven_int32 := int32(phone_seven_int)
	right := (total_len - firstoffset) / PHONE_INDEX_LENGTH
	for {
		if left > right {
			break
		}
		mid := (left + right) / 2
		offset := firstoffset + mid*PHONE_INDEX_LENGTH
		if offset >= total_len {
			break
		}
		cur_phone := get4(content[offset : offset+INT_LEN])
		record_offset := get4(content[offset+INT_LEN : offset+INT_LEN*2])
		card_type := content[offset+INT_LEN*2 : offset+INT_LEN*2+CHAR_LEN][0]
		switch {
		case cur_phone > phone_seven_int32:
			right = mid - 1
		case cur_phone < phone_seven_int32:
			left = mid + 1
		default:
			if checkMode {
				return nil, nil
			}

			cbyte := content[record_offset:]
			end_offset := int32(bytes.Index(cbyte, []byte("\000")))
			data := bytes.Split(cbyte[:end_offset], []byte("|"))
			card_str, ok := CardTypemap[card_type]
			if !ok {
				card_str = "未知电信运营商"
			}
			pr = &PhoneRecord{
				PhoneNum: phone_num,
				Province: string(data[0]),
				City:     string(data[1]),
				ZipCode:  string(data[2]),
				AreaZone: string(data[3]),
				CardType: card_str,
			}
			err = nil
			return
		}
	}
	return nil, errors.New("phone's data not found")
}

func init() {
	phoneData := conf.Str("", "s:/phonefind/datafile")
	Init(phoneData)
}
