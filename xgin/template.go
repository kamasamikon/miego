package xgin

import (
	"fmt"
	"html/template"
	"time"

	"github.com/kamasamikon/miego/atox"
	"github.com/kamasamikon/miego/xmap"
	"github.com/kamasamikon/miego/xtime"
)

func ToHTML(x string) interface{} {
	return template.HTML(x)
}
func ToJS(x string) interface{} {
	return template.JS(x)
}
func ToCSS(x string) interface{} {
	return template.CSS(x)
}

func SubStr(s string, beg int, end int) string {
	slen := len(s) + 1

	beg = ((beg % slen) + slen) % slen
	end = ((end % slen) + slen) % slen

	if beg >= end {
		return ""
	}
	return s[beg:end]
}

// 比较a和b，相等就返回eqstr，否则就返回nestr
func Choice(a string, b string, eqstr string, nestr string) string {
	if a == b {
		return eqstr
	}
	return nestr
}
func ToAttr(s string) template.HTMLAttr {
	return template.HTMLAttr(s)
}

// 返回obj[name]，如果不存在，返回defval
func MPGet(obj xmap.Map, name string, defval string) string {
	return obj.Str(name, defval)
}

// 结合MPGet和Choice
func MapChoice(obj xmap.Map, name string, check string, eqstr string, nestr string) string {
	return Choice(obj.Str(name, ""), check, eqstr, nestr)
}

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func NtimeToString(t string) string {
	nt, _ := xtime.NumTimeToTime(atox.Uint64(t, 0))
	return nt.Format("2006-01-02 15:04:05")
}
