package atox

import (
	"strconv"
)

// Atoi : atoi, if fail return default value
func Int64(a string, def int64) int64 {
	x, e := strconv.ParseInt(a, 0, 64)
	if e != nil {
		return def
	}
	return x
}

// Atoi : atoi, if fail return default value
func Uint64(a string, def uint64) uint64 {
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
