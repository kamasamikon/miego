package conf

import (
	"strconv"
	"strings"
)

func (cc *ConfCenter) EntryRem(path string) {
	kind, key := pathParse(path)
	if key == "" {
		return
	}

	switch kind {
	case 'i':
		cc.IRem(key)
	case 's':
		cc.SRem(key)
	case 'b':
		cc.BRem(key)
	case 'e':
		cc.ERem(key)
	default:
		return
	}
}

func (cc *ConfCenter) EntryAddByLine(line string, overwrite bool) {
	segs := strings.SplitN(strings.TrimSpace(line), "=", 2)
	if len(segs) < 2 {
		return
	}

	path, value := segs[0], segs[1]
	cc.EntryAdd(path, value, overwrite)
}

func (cc *ConfCenter) EntryAdd(path string, value string, overwrite bool) {
	kind, key := pathParse(path)
	if key == "" {
		return
	}

	switch kind {
	case 'i':
		if vInt, err := strconv.ParseInt(value, 10, 64); err == nil {
			if overwrite {
				cc.ISetf(key, vInt)
			} else {
				cc.ISet(key, vInt)
			}
		} else {
			dp("BadValue.i: %s", value)
			return
		}

	case 's':
		if overwrite {
			cc.SSetf(key, value)
		} else {
			cc.SSet(key, value)
		}

	case 'b':
		// true: 1, t, T
		// false: 0, f, F
		x := value[0]
		if x == '1' || x == 't' || x == 'T' || x == 'y' || x == 'Y' {
			if overwrite {
				cc.BSetf(key, true)
			} else {
				cc.BSet(key, true)
			}
		} else if x == '0' || x == 'f' || x == 'F' || x == 'n' || x == 'N' {
			if overwrite {
				cc.BSetf(key, false)
			} else {
				cc.BSet(key, false)
			}
		} else {
			dp("BadValue.b: %s", value)
			return
		}

	case 'e':
		// event, no data at all, value treated as a parameter
		if overwrite {
			cc.ESendf(key, value)
		} else {
			cc.ESendf(key, value)
		}

	default:
		return
	}
}
