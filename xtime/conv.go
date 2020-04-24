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
			return t.Format("2006/01/02 15:04:05")
		}
	} else if l == 12 {
		// Minute
		if t, err := time.Parse("200601021504", timestr); err == nil {
			return t.Format("2006/01/02 15:04")
		}
	} else if l == 8 {
		// Day
		if t, err := time.Parse("20060102", timestr); err == nil {
			return t.Format("2006/01/02")
		}
	}

	return ""
}

func StrToNum(s string) uint64 {
	if t, err := time.Parse("2006-01-02", s); err == nil {
		return atox.Uint64(t.Format("20060102"), 0)
	}
	if t, err := time.Parse("2006-01-02 15:04", s); err == nil {
		return atox.Uint64(t.Format("20060021504"), 0)
	}
	if t, err := time.Parse("2006-01-02 15:04:05", s); err == nil {
		return atox.Uint64(t.Format("2006002150405"), 0)
	}
	return 0
}
