package xtime

import (
	"fmt"
	"github.com/kamasamikon/miego/atox"
	"time"
)

func NumToStr(o interface{}) string {
	var timestr string

	if s, ok := o.(string); ok {
		timestr = s
	} else if n, ok := o.(uint64); ok {
		timestr = fmt.Sprintf("%d", n)
	}

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

func StrToNum(s string, flag byte) uint64 {
	var tmp string

	if t, err := time.Parse("2006-01-02", s); err == nil {
		tmp = t.Format("20060102")
	}
	if t, err := time.Parse("2006-01-02 15:04", s); err == nil {
		tmp = t.Format("200601021504")
	}
	if t, err := time.Parse("2006-01-02 15:04:05", s); err == nil {
		tmp = t.Format("20060102150405")
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

	return atox.Uint64(tmp, 0)
}
