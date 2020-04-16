package roast

import (
	"fmt"
	"strings"
)

func likeParse(s string, fmtMode bool) (out string, isLike bool) {
	var ss string

	ss = strings.Replace(s, "*", "%", -1)
	ss = strings.Replace(ss, ".", "_", -1)

	if fmtMode {
		ss = strings.Replace(ss, "%", "%%", -1)
	}

	if strings.IndexByte(ss, '%') >= 0 {
		return ss, true
	}
	if strings.IndexByte(ss, '_') >= 0 {
		return ss, true
	}

	return ss, false
}

func LIKE(qList []string, name string, field string) []string {
	if s, like := likeParse(name, false); like {
		return append(qList, fmt.Sprintf("%s LIKE '%s'", field, s))
	} else {
		return append(qList, fmt.Sprintf("%s = '%s'", field, s))
	}
}
