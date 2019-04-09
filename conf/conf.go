package conf

import (
	"bufio"
	"fmt"
	"github.com/kamasamikon/miego/klog"
	"io"
	"os"
	"strconv"
	"strings"
)

// See confcenter
type confEntry struct {
	// kind: i:int, s:str
	// kind: a:arr, b:bool, d:dat(len+dat), e:event, i:int, s:str, p:ptr
	// path: i:/aaa/bbb
	// vXxx: value for each type
	kind    byte
	path    string
	vInt    int64
	vStr    string
	vBool   bool
	deleted bool
}

// KConfMonitor is a Callback called when wathed entry modified.
type KConfMonitor func(path string, oVal interface{}, nVal interface{})

type confMonitors struct {
	funcList []KConfMonitor
}

var mapPathEntry = make(map[string]*confEntry)
var mapPathMonitors = make(map[string]*confMonitors)

// Load file from configure
func entryAdd(line string) {
	segs := strings.SplitN(line, "=", 2)
	if len(segs) < 2 {
		return
	}

	path, value := segs[0], segs[1]
	kind := path[0]

	e := &confEntry{
		kind: kind,
		path: path,
	}

	switch kind {
	case 'i':
		if vInt, err := strconv.ParseInt(value, 10, 64); err == nil {
			e.vInt = vInt
		}
	case 's':
		e.vStr = value
	case 'b':
		// true: 1, t, T
		// false: 0, f, F
		x := value[0]
		if x == '1' || x == 't' || x == 'T' {
			e.vBool = true
		} else {
			e.vBool = false
		}
	}

	mapPathEntry[path] = e
}

// Load configure from a file.
func Load(fileName string) error {
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	buf := bufio.NewReader(f)
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		entryAdd(line)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
	}
}

// Exists : Check if a entry exists
func Exists(path string) bool {
	_, ok := mapPathEntry[path]
	return ok
}

// Int : get a int typed configure
func Int(path string, defval int64) int64 {
	// path: aaa/bbb
	key := "i:/" + path
	if v, ok := mapPathEntry[key]; ok {
		return v.vInt
	}
	return defval
}

// Str : get a str typed configure
func Str(path string, defval string) string {
	// path: aaa/bbb
	key := "s:/" + path
	if v, ok := mapPathEntry[key]; ok {
		return v.vStr
	}
	return defval
}

// Bool : get a bool entry
func Bool(path string, defval bool) bool {
	// path: aaa/bbb
	key := "b:/" + path
	if v, ok := mapPathEntry[key]; ok {
		return v.vBool
	}
	return defval
}

// Set : Modify or Add conf entry
func Set(path string, value interface{}, force bool) {
	var e *confEntry
	var ok bool

	kind := path[0]

	e, ok = mapPathEntry[path]
	if !ok {
		if force {
			e = &confEntry{
				kind: kind,
				path: path,
			}
		} else {
			klog.D("path:%s and force:false", path)
			return
		}
	}

	switch kind {
	case 'i':
		v := value.(int64)
		monitorCall(e, e.vInt, v)
		e.vInt = v
	case 's':
		v := value.(string)
		monitorCall(e, e.vStr, v)
		e.vStr = v
	case 'b':
		v := value.(bool)
		monitorCall(e, e.vBool, v)
		e.vBool = v
	default:
		klog.E("Bad Kind: %c", kind)
	}
}

// Monitor : Callback when entry changed.
func Monitor(path string, callback func(path string, oVal interface{}, nVal interface{})) {
	// path: i:/aaa/bbb

	var m *confMonitors
	if monitors, ok := mapPathMonitors[path]; ok {
		m = monitors
	} else {
		m = &confMonitors{}
		mapPathMonitors[path] = m
	}

	m.funcList = append(m.funcList, callback)
}

func monitorCall(e *confEntry, oVal interface{}, nVal interface{}) {
	if monitors, ok := mapPathMonitors[e.path]; ok {
		for _, f := range monitors.funcList {
			if f != nil {
				go f(e.path, oVal, nVal)
			}
		}
	}
}

// Dump : Print all entries
func Dump() {
	for p, v := range mapPathEntry {
		switch v.kind {
		case 'i':
			fmt.Printf("%20s : %c : %d\n", p, v.kind, v.vInt)
		case 's':
			fmt.Printf("%20s : %c : %s\n", p, v.kind, v.vStr)
		case 'b':
			fmt.Printf("%20s : %c : %t\n", p, v.kind, v.vBool)
		}
	}
}

func init() {
	path := "./aster.cfg"
	Load(path)
}
