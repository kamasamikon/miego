package xtime

import (
	"fmt"
	"github.com/kamasamikon/miego/atox"
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

func StrToNum(s string, flag byte) int64 {
	var tmp string

	for {
		if t, err := time.Parse("20060102", s); err == nil {
			tmp = t.Format("20060102")
			break
		}
		if t, err := time.Parse("200601021504", s); err == nil {
			tmp = t.Format("200601021504")
			break
		}

		if t, err := time.Parse("2006-01-02", s); err == nil {
			tmp = t.Format("20060102")
			break
		}
		if t, err := time.Parse("2006-01-02 15:04", s); err == nil {
			tmp = t.Format("200601021504")
			break
		}
		if t, err := time.Parse("2006-01-02  15:04", s); err == nil {
			tmp = t.Format("200601021504")
			break
		}
		if t, err := time.Parse("2006-01-02 15:04:05", s); err == nil {
			tmp = t.Format("20060102150405")
			break
		}
		if t, err := time.Parse("2006-01-02  15:04:05", s); err == nil {
			tmp = t.Format("200601021504")
			break
		}

		if t, err := time.Parse("2006-1-2", s); err == nil {
			tmp = t.Format("20060102")
			break
		}
		if t, err := time.Parse("2006-1-2 3:4", s); err == nil {
			tmp = t.Format("200601021504")
			break
		}
		if t, err := time.Parse("2006-1-2  3:4", s); err == nil {
			tmp = t.Format("200601021504")
			break
		}
		if t, err := time.Parse("2006-1-2 3:4:5", s); err == nil {
			tmp = t.Format("20060102150405")
			break
		}
		if t, err := time.Parse("2006-1-2  3:4:5", s); err == nil {
			tmp = t.Format("200601021504")
			break
		}

		return 0
	}

	tmp += "000000000"
	switch flag {
	case 'N', 'n':
		tmp = tmp[0:4]
	case 'Y', 'y':
		tmp = tmp[0:6]
	case 'R', 'r':
		tmp = tmp[0:8]
	case 'S', 's':
		tmp = tmp[0:10]
	case 'F', 'f':
		tmp = tmp[0:12]
	case 'M', 'm':
		tmp = tmp[0:14]
	}

	return atox.Int64(tmp, 0)
}

func AnyToNum(g interface{}) int64 {
	stime := fmt.Sprintf("%v000000", g)

	l := len(stime)

	if l < 8 {
		return 0
	}

	if stime[4] == '-' {
		if l >= 19 {
			if t, err := time.Parse("2006-01-02 15:04:05", stime[0:19]); err == nil {
				return TimeToNum(t)
			}
		} else {
			if t, err := time.Parse("2006-01-02", stime[0:10]); err == nil {
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
