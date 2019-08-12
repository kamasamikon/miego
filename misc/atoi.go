package misc

import (
	"strconv"
)

// Atoi : atoi, if fail return default value
func AtoInt(a string, def int) int {
	x, e := strconv.ParseInt(a, 0, 64)
	if e != nil {
		return def
	}
	return int(x)
}

// Atoi : atoi, if fail return default value
func AtoUint(a string, def uint) uint {
	x, e := strconv.ParseInt(a, 0, 64)
	if e != nil {
		return def
	}
	return uint(x)
}

// Atoi : atoi, if fail return default value
func AtoInt64(a string, def int64) int64 {
	x, e := strconv.ParseInt(a, 0, 64)
	if e != nil {
		return def
	}
	return int64(x)
}

// Atoi : atoi, if fail return default value
func AtoUint64(a string, def uint64) uint64 {
	x, e := strconv.ParseInt(a, 0, 64)
	if e != nil {
		return def
	}
	return uint64(x)
}
