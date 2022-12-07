package xtime

import (
	"fmt"
	"strconv"
	"strings"
	"time"
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

	return
}

// 2016-01-02 15:03:05
// 2016-01-02
// 2016/01/02 15:03:05
// 2016/01/02
// 20060102150405
func AnyToNum(g interface{}) uint64 {
	s := fmt.Sprintf("%v", g)

	s = strings.Replace(s, ":", " ", -1) // 15:03:05
	s = strings.Replace(s, "-", " ", -1) // 2016-01-02
	s = strings.Replace(s, "/", " ", -1) // 2016/01/02
	s = strings.Replace(s, "T", " ", -1) // 2017-03-04T01:23:43.000Z
	s = strings.Replace(s, ".", " ", -1) // 2017-03-04T01:23:43.000Z

	// s = "2017 03 04 01 23 43 000Z"
	// s = "20170304012343
	var segarr []string
	for _, s := range strings.Split(s, " ") {
		if s != "" {
			segarr = append(segarr, s)
		}
	}

	segcnt := len(segarr)

	timefmt := ""

	// 2016 01
	if segcnt >= 2 {
		timefmt += "2006"
		tmp := len(segarr[1]) // 月
		if tmp == 1 {
			timefmt += "1"
		} else if tmp == 2 {
			timefmt += "01"
		}
	}

	// 2016 01 02
	if segcnt >= 3 {
		tmp := len(segarr[2]) // 日
		if tmp == 1 {
			timefmt += "2"
		} else if tmp == 2 {
			timefmt += "02"
		}
	}

	// 2016 01 02 15
	if segcnt >= 4 {
		tmp := len(segarr[3]) // 时
		if tmp == 1 {
			timefmt += "3"
		} else if tmp == 2 {
			timefmt += "15"
		}
	}

	// 2016 01 02 15 04
	if segcnt >= 5 {
		tmp := len(segarr[4]) // 分
		if tmp == 1 {
			timefmt += "4"
		} else if tmp == 2 {
			timefmt += "04"
		}
	}

	// 2016 01 02 15 04 05
	if segcnt >= 6 {
		tmp := len(segarr[5]) // 秒
		if tmp == 1 {
			timefmt += "5"
		} else if tmp == 2 {
			timefmt += "05"
		}
	}

	// 201601021503
	// 2016010215
	// 20160102
	// 201601
	// 2016
	if segcnt == 1 {
		tmp := len(segarr[0])

		// 20160102150305
		if tmp >= 14 {
			timefmt = "20060102150405"
		}

		// 201601021503
		if tmp == 12 {
			timefmt = "200601021504"
		}
		// 2016010215
		if tmp == 10 {
			timefmt = "2006010215"
		}
		// 20160102
		if tmp == 8 {
			timefmt = "20060102"
		}
		// 201601
		if tmp == 6 {
			timefmt = "200601"
		}
		// 2016
		if tmp == 4 {
			timefmt = "2006"
		}
	}

	stime := strings.Join(segarr, "")
	if t, err := time.Parse(timefmt, stime[0:len(timefmt)]); err == nil {
		return TimeToNum(t)
	}

	return 0
}

// OffsetDate : 把绝对日期或者日期偏移量调整为真正的日期
func OffsetDate(sDate string) uint64 {
	now := time.Now()

	if sDate[0] == '+' || sDate[0] == '-' {
		offset, _ := strconv.ParseInt(sDate, 0, 64)
		tmp := now.AddDate(0, 0, int(offset)).Format("20060102")

		x, _ := strconv.ParseUint(tmp, 0, 64)
		return x
	}
	x, _ := strconv.ParseUint(sDate, 0, 64)
	return x
}

// FormatDuration : 转成 X 天 x 小时 ...
func FormatDuration(seconds uint64) string {
	MM := seconds % 60
	FF := (seconds / 60) % 60
	SS := (seconds / 60 / 60) % 24
	RR := (seconds / 60 / 60 * 24)

	s := ""
	if RR != 0 {
		s += fmt.Sprintf("%d天", RR)
	}
	if SS != 0 {
		s += fmt.Sprintf("%d小时", SS)
	}
	if FF != 0 {
		s += fmt.Sprintf("%d分", FF)
	}
	if MM != 0 {
		s += fmt.Sprintf("%d秒", MM)
	}

	return s
}
