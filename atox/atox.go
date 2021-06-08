package atox

import (
	"strconv"
	"strings"
)

func Prepare(s string) string {
	os := strings.ReplaceAll(s, " ", "")
	os = strings.ReplaceAll(os, "\t", "")
	os = strings.ReplaceAll(os, "\r", "")
	os = strings.ReplaceAll(os, "\n", "")
	if os == "" {
		return ""
	}

	if len(os) == 1 {
		return os
	}

	var prefix string
	if os[0] == '+' || os[0] == '-' {
		prefix = os[0:1]
		os = os[1:]
		if len(os) == 1 {
			return prefix + os
		}
	}

	if os[0] == '0' && (os[1] == 'x' || os[1] == 'X') {
		return prefix + os
	}

	ns := strings.TrimLeft(os, "0")
	if ns == "" {
		return prefix + "0"
	}

	return prefix + ns
}

// Atoi : atoi, if fail return default value
func Int64(a string, def int64) int64 {
	a = Prepare(a)
	x, e := strconv.ParseInt(a, 0, 64)
	if e != nil {
		return def
	}
	return x
}

// Atoi : atoi, if fail return default value
func Uint64(a string, def uint64) uint64 {
	a = Prepare(a)
	x, e := strconv.ParseUint(a, 0, 64)
	if e != nil {
		return def
	}
	return x
}

// Atoi : atoi, if fail return default value
func Int(a string, def int) int {
	return int(Int64(a, int64(def)))
}

// Atoi : atoi, if fail return default value
func Uint(a string, def uint) uint {
	return uint(Uint64(a, uint64(def)))
}

func Float(a string, def float64) float64 {
	if f, err := strconv.ParseFloat(a, 64); err == nil {
		return f
	}
	return def
}
