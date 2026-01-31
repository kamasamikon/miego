package version

import (
	"fmt"
	"strconv"
	"strings"
)

// "30065229831" > 30065229831
// "7.7.7" > 30065229831
// 30065229831 > 30065229831
func N(v interface{}) uint64 {
	vv := fmt.Sprintf("%v", v)

	// "30065229831" > 30065229831
	// "7.7.7" > 30065229831
	segs := strings.Split(vv, ".")

	// "30065229831" > 30065229831
	if len(segs) == 1 {
		if num, err := strconv.ParseUint(vv, 0, 64); err != nil {
			return 0
		} else {
			return num
		}
	}

	// "7.7.7" > 30065229831
	if len(segs) == 3 {
		a, err := strconv.ParseUint(segs[0], 0, 64)
		if err != nil {
			return 0
		}
		b, err := strconv.ParseUint(segs[1], 0, 64)
		if err != nil {
			return 0
		}
		c, err := strconv.ParseUint(segs[2], 0, 64)
		if err != nil {
			return 0
		}

		var n uint64
		n += (a << 32) & 0xffff00000000
		n += (b << 16) & 0xffff0000
		n += (c << 0) & 0xffff

		return n
	}

	return 0
}

// "30065229831" > "7.7.7"
// "7.7.7" > "7.7.7"
// 30065229831 > "7.7.7"
func S(v interface{}) string {
	n := N(v)
	a := (n >> 32) & 0xffff
	b := (n >> 16) & 0xffff
	c := (n >> 0) & 0xffff
	return fmt.Sprintf("%d.%d.%d", a, b, c)
}

// a.b.c => number
func S2N(s string) (uint64, error) {
	// XXX: a.b.c 版本转成内部格式
	segs := strings.Split(s, ".")
	if len(segs) < 3 {
		return 0, fmt.Errorf("bad version")
	}
	a, err := strconv.ParseUint(segs[0], 0, 64)
	if err != nil {
		return 0, fmt.Errorf("bad version")
	}
	b, err := strconv.ParseUint(segs[1], 0, 64)
	if err != nil {
		return 0, fmt.Errorf("bad version")
	}
	c, err := strconv.ParseUint(segs[2], 0, 64)
	if err != nil {
		return 0, fmt.Errorf("bad version")
	}

	var n uint64
	n += (a << 32) & 0xffff00000000
	n += (b << 16) & 0xffff0000
	n += (c << 0) & 0xffff

	return n, nil
}

// number => a.b.c
func N2S(n uint64) string {
	a := (n >> 32) & 0xffff
	b := (n >> 16) & 0xffff
	c := (n >> 0) & 0xffff
	return fmt.Sprintf("%d.%d.%d", a, b, c)
}
