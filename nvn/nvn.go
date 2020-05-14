package nvn

import (
	"fmt"
)

// Name VS Number

var mapInt map[string]int
var mapStr map[string]string

func N(s interface{}, class string) int {
	var key string

	if x, ok := s.(string); ok {
		key = fmt.Sprintf("%s___%s", class, x)
	} else {
		key = fmt.Sprintf("%s___%d", class, s.(int))
	}

	if n, ok := mapInt[key]; ok {
		return n
	}
	return -1
}

func S(n interface{}, class string) string {
	var key string

	if x, ok := n.(string); ok {
		key = fmt.Sprintf("%s___%s", class, x)
	} else {
		key = fmt.Sprintf("%s___%d", class, n.(int))
	}

	if s, ok := mapStr[key]; ok {
		return s
	}
	return ""
}

func Add(n int, s string, class string) {
	var key string

	key = fmt.Sprintf("%s___%s", class, s)
	mapInt[key] = n

	key = fmt.Sprintf("%s___%d", class, n)
	mapStr[key] = s
}

func init() {
	mapInt = make(map[string]int)
	mapStr = make(map[string]string)
}
