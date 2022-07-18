package xtime

import (
	"fmt"
	"time"

	"github.com/kamasamikon/miego/atox"
	"github.com/kamasamikon/miego/klog"
)

func NumToStr(o interface{}) string {
	timestr := fmt.Sprintf("%v", o)

	l := len(timestr)
	if l == 14 {
		// Second
		if t, err := time.Parse("20060102150405", timestr); err == nil {
			return t.Format("2006-01-02 15:04:05")
		}
	} else if l == 12 {
		// Minute
		if t, err := time.Parse("200601021504", timestr); err == nil {
			return t.Format("2006-01-02 15:04")
		}
	} else if l == 8 {
		// Day
		if t, err := time.Parse("20060102", timestr); err == nil {
			return t.Format("2006-01-02")
		}
	}

	return ""
}

func StrToNum(s string, flag byte) (Num uint64) {
	var NNNN uint64
	var YY uint64
	var RR uint64
	var SS uint64
	var FF uint64
	var MM uint64

	curStage := 0 // Stage of NNNN YY RR SS FF MM
	maxStage := 0
	switch flag {
	case 'N':
		maxStage = 0
	case 'Y':
		maxStage = 1
	case 'R':
		maxStage = 2
	case 'S':
		maxStage = 3
	case 'F':
		maxStage = 4
	case 'M':
		maxStage = 5
	}

	dcount := 0 // digital number count
	slen := len(s)

	for i := 0; i < slen; i++ {
		if curStage > maxStage {
			break
		}

		n := uint64(s[i]) - '0'
		if n < 0 || n > 9 {
			// 不是数字。
			// 当前字段已经有数字了？
			//     已经有数字了，转到下一个字段
			//     如果没有数字，继续等
			if dcount > 0 {
				// 有数字，下一步
				dcount = 0
				curStage++
			}
			continue
		}

		switch curStage {
		case 0: // NNNN
			NNNN = NNNN*10 + n
			dcount++

			if dcount == 4 {
				dcount = 0
				curStage++
			}

		case 1: // YY
			YY = YY*10 + n
			dcount++

			if dcount == 2 {
				dcount = 0
				curStage++
			}

		case 2: // RR
			RR = RR*10 + n
			dcount++

			if dcount == 2 {
				dcount = 0
				curStage++
			}

		case 3: // SS
			SS = SS*10 + n
			dcount++

			if dcount == 2 {
				dcount = 0
				curStage++
			}

		case 4: // FF
			FF = FF*10 + n
			dcount++

			if dcount == 2 {
				dcount = 0
				curStage++
			}

		case 5: // MM
			MM = MM*10 + n
			dcount++

			if dcount == 2 {
				dcount = 0
				curStage++
			}
		}
	}

	if YY == 0 {
		YY = 1
	}
	if RR == 0 {
		RR = 1
	}
	switch flag {
	case 'N':
		Num = NNNN
	case 'Y':
		Num = NNNN*100 + YY
	case 'R':
		Num = NNNN*10000 + YY*100 + RR
	case 'S':
		Num = NNNN*1000000 + YY*10000 + RR*100 + SS
	case 'F':
		Num = NNNN*100000000 + YY*1000000 + RR*10000 + SS*100 + FF
	case 'M':
		Num = NNNN*10000000000 + YY*100000000 + RR*1000000 + SS*10000 + FF*100 + MM
	}

	klog.F("FLAG:%c\tS:(%v)\tN:%v\n", flag, s, Num)
	return
}

func AnyToNum(g interface{}) uint64 {
	stime := fmt.Sprintf("%v000000", g)

	l := len(stime)
	if l < 8 {
		return 0
	}

	if stime[4] == '-' {
		// 2016-01-02 15:03:05 or 2016-01-02
		if l >= 19 {
			if t, err := time.Parse("2006-01-02 15:04:05", stime[0:19]); err == nil {
				return TimeToNum(t)
			}
		} else {
			if t, err := time.Parse("2006-01-02", stime[0:10]); err == nil {
				return TimeToNum(t)
			}
		}
	} else if stime[4] == '/' {
		// 2016/01/02 15:03:05 or 2016/01/02
		if l >= 19 {
			if t, err := time.Parse("2006/01/02 15:04:05", stime[0:19]); err == nil {
				return TimeToNum(t)
			}
		} else {
			if t, err := time.Parse("2006/01/02", stime[0:10]); err == nil {
				return TimeToNum(t)
			}
		}
	} else {
		if l >= 18 {
			// 20160102150305
			if t, err := time.Parse("20060102150405", stime[0:14]); err == nil {
				return TimeToNum(t)
			}
		} else {
			// 20160102
			if t, err := time.Parse("20060102", stime[0:8]); err == nil {
				return TimeToNum(t)
			}
		}
	}

	return 0
}

// OffsetDate : 把绝对日期或者日期偏移量调整为真正的日期
func OffsetDate(sDate string) uint64 {
	now := time.Now()

	if sDate[0] == '+' || sDate[0] == '-' {
		offset := atox.Int(sDate, 0)
		tmp := now.AddDate(0, 0, offset).Format("20060102")
		return atox.Uint64(tmp, 0)
	}
	return atox.Uint64(sDate, 0)
}
