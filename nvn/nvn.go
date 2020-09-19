package nvn

import (
	"fmt"

	"github.com/kamasamikon/miego/conf"
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

func LoadConf() {
	for _, name := range conf.Name() {
		segs := strings.Split(name, "/")
		// segs := ["i:", "nvn", "<class>", "<name>"]
		// i:/nvn/Role/管理员=1
		if len(segs) == 4 && segs[0] == "i:" && segs[1] == "nvn" {
			n := conf.Int(0, name)
			class := segs[2]
			s := segs[3]

			Add(int(n), s, class)
		}
	}
}

func Dump() {
	var lines []string

	for k, v := range mapInt {
		segs := strings.Split(k, "___")
		class, s := segs[0], segs[1]
		lines = append(lines, fmt.Sprintf("%s\t %s:%d", class, s, v)
	}

	lines = append(lines, "")
	return strings.Join(lines, "\n")
}

func init() {
	// i:/nvn/<class>/<name>=v
	mapInt = make(map[string]int)
	mapStr = make(map[string]string)
}
