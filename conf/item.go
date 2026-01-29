package conf

import (
	"encoding/json"
	"strconv"
	"strings"
)

func (cc *ConfCenter) EntryRem(path string) {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	delete(cc.mapPathEntry, path)
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
	kind, hidden, realpath := pathParse(path)
	if realpath == "" {
		return
	}

	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	e, exists := cc.mapPathEntry[path]
	if exists && !overwrite {
		return
	}

	var vNew any

	switch kind {
	case 'i':
		if vInt, err := strconv.ParseInt(value, 10, 64); err == nil {
			vNew = vInt
		} else {
			dp("BadValue.i: %s", value)
			return
		}
		cc.ISet(realpath, vNew)

	case 's':
		vNew = value
		cc.ISet(realpath, vNew)

	case 'b':
		// true: 1, t, T
		// false: 0, f, F
		x := value[0]
		if x == '1' || x == 't' || x == 'T' || x == 'y' || x == 'Y' {
			vNew = true
		} else if x == '0' || x == 'f' || x == 'F' || x == 'n' || x == 'N' {
			vNew = false
		} else {
			dp("BadValue.b: %s", value)
			return
		}
		cc.ISet(realpath, vNew)

	case 'o':
		// line is json string
		o := make(map[string]any)
		if err := json.Unmarshal([]byte(value), &o); err != nil {
			dp("BadValue.o: %s", value)
			return
		}
		vNew = o
		cc.ISet(realpath, vNew)

	case 'e':
		// event, no data at all, value treated as a parameter
		vNew = value
		cc.ISet(realpath, vNew)

	default:
		return
	}

	if !exists {
		e = &confEntry{
			kind:   kind,
			hidden: hidden,
			path:   realpath,
		}
		cc.mapPathEntry[e.path] = e
	}
	cc.setByEntry(e, vNew)
}
