package nvn

import (
	"fmt"
	"strings"

	"miego/conf"
)

// String VS String

// "aaa" VS "AAA"
// Class__aaa => AAA
// Class__AAA => aaa
var mapsvs map[string]string

func S(s string, class string) string {
	key := class + "__" + s
	if s, ok := mapsvs[key]; ok {
		return s
	}
	return ""
}

func Add(sa string, sb string, class string) {
	a := class + "__" + sa
	b := class + "__" + sb
	mapsvs[a] = sb
	mapsvs[b] = sa
}

func LoadConf(prefix string) {
	if prefix == "" {
		prefix = "svs"
	}
	for _, name := range conf.Names() {
		segs := strings.Split(name, "/")
		// segs := ["s:", "svs", "<class>", "<name>"]
		// s:/svs/Role/管理员=1	=> svs.Add(1, "管理员", "Role")
		if len(segs) == 4 && segs[0] == "i:" && segs[1] == prefix {
			sa := conf.S(name)
			if sa == "" {
				continue
			}

			class := segs[2]
			sb := segs[3]

			Add(sa, sb, class)
		}
	}
}

func Dump() string {
	var lines []string

	// k: cls_a v: cls_b
	for a, b := range mapsvs {
		segs_a := strings.Split(a, "___")
		class_a, sa := segs_a[0], segs_a[1]

		segs_b := strings.Split(b, "___")
		class_b, sb := segs_b[0], segs_b[1]

		if class_a != class_b {
			continue
		}

		lines = append(lines, fmt.Sprintf("%s\t %s:%s", class_a, sa, sb))
	}

	lines = append(lines, "")
	return strings.Join(lines, "\n")
}

func init() {
	// i:/nvn/<class>/<name>=v
	mapsvs = make(map[string]string)
}
