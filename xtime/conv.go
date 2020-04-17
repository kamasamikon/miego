package xtime

import (
	"fmt"
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
		if t, err := time.Parse("20060102150405", timestr); err == nil {
			return t.Format("2006/01/02 15:04:05")
		}
	} else if l == 8 {
		if t, err := time.Parse("20060102", timestr); err == nil {
			return t.Format("2006/01/02")
		}
	}

	return ""
}
