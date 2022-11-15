package version

import (
	"fmt"
	"strconv"
	"strings"
)

// a.b.c => number
func S2N(s string) (uint64, error) {
	// XXX: a.b.c 版本转成内部格式
	segs := strings.Split(s, ".")
	if len(segs) < 3 {
		return 0, fmt.Errorf("Bad SWVersion")
	}
	a, err := strconv.ParseUint(segs[0], 0, 64)
	if err != nil {
		return 0, fmt.Errorf("Bad SWVersion")
	}
	b, err := strconv.ParseUint(segs[0], 0, 64)
	if err != nil {
		return 0, fmt.Errorf("Bad SWVersion")
	}
	c, err := strconv.ParseUint(segs[0], 0, 64)
	if err != nil {
		return 0, fmt.Errorf("Bad SWVersion")
	}

	var n uint64 = 0
	n += (a << 8) & 0xffff00000000
	n += (b << 4) & 0xffff0000
	n += (c << 0) & 0xffff

	return n, nil
}

// number => a.b.c
func N2S(n uint64) string {
	a := (n >> 8) & 0xffff
	b := (n >> 4) & 0xffff
	c := (n >> 0) & 0xffff
	return fmt.Sprintf("%d.%d.%d", a, b, c)
}
