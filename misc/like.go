package misc

import (
	"strings"
)

func KrLike(s string) (out string, isLike bool) {
	var ss string

	ss = strings.Replace(s, "*", "%", -1)
	ss = strings.Replace(ss, ".", "_", -1)

	if strings.IndexByte(ss, '%') >= 0 {
		return ss, true
	}
	if strings.IndexByte(ss, '_') >= 0 {
		return ss, true
	}

	return ss, false
}
