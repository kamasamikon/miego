package xgin

import (
	"fmt"
	"html/template"
	"time"
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
func ToAttr(s string) template.HTMLAttr {
	return template.HTMLAttr(s)
}

// Add to number
func NumAdd(a int, others ...int) int {
	for _, x := range others {
		a += x
	}
	return a
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

// 返回obj[name]，如果不存在，返回defval
func MapGet(obj map[string]interface{}, name string, defval string) string {
	if x, ok := obj[name]; ok {
		if s, ok := x.(string); ok {
			return s
		}
	}
	return defval
}

// 结合MapGet和Choice
func MapChoice(obj map[string]interface{}, name string, check string, eqstr string, nestr string) string {
	val := ""
	if x, ok := obj[name]; ok {
		if s, ok := x.(string); ok {
			val = s
		}
	}
	return Choice(val, check, eqstr, nestr)
}

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

// "20060102150405" => "2006-01-02 15:04:05"
func NtimeToString(s string) string {
	if len(s) == 14 {
		// "20060102150405"
		if t, err := time.Parse("20060102150405", s); err == nil {
			return t.Format("2006-01-02 15:04:05")
		}
	}
	if len(s) == 8 {
		// "20060102"
		if t, err := time.Parse("20060102", s); err == nil {
			return t.Format("2006-01-02")
		}
	}
	return "NA"
}
