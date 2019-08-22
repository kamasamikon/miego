package conf

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/kamasamikon/miego/klog"
)

// See confcenter
type confEntry struct {
	// kind: i:int, s:str
	// kind: a:arr, b:bool, d:dat(len+dat), e:event, i:int, s:str, p:ptr
	//
	// path: i:/aaa/bbb; b:/xxx/zzz
	//
	// vXxx: value for each type
	//
	// refGet/refSet: count by Read or Write
	kind byte
	path string

	vInt  int64
	vStr  string
	vBool bool

	refGet int64
	refSet int64
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
		if x == '1' || x == 't' || x == 'T' || x == 'y' || x == 'Y' {
			e.vBool = true
		} else {
			e.vBool = false
		}
	}

	// Simply overwrite the old value.
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
func Int(defval int64, paths ...string) int64 {
	// path: aaa/bbb
	for _, path := range paths {
		key := "i:/" + path
		if v, ok := mapPathEntry[key]; ok {
			v.refGet++
			return v.vInt
		}
	}
	return defval
}

// Str : get a str typed configure
func Str(defval string, paths ...string) string {
	// path: aaa/bbb
	for _, path := range paths {
		key := "s:/" + path
		if v, ok := mapPathEntry[key]; ok {
			v.refGet++
			return v.vStr
		}
	}
	return defval
}

// Bool : get a bool entry
func Bool(defval bool, paths ...string) bool {
	// path: aaa/bbb
	for _, path := range paths {
		key := "b:/" + path
		if v, ok := mapPathEntry[key]; ok {
			v.refGet++
			return v.vBool
		}
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
		e.refSet++

	case 's':
		v := value.(string)
		monitorCall(e, e.vStr, v)
		e.vStr = v
		e.refSet++

	case 'b':
		v := value.(bool)
		monitorCall(e, e.vBool, v)
		e.vBool = v
		e.refSet++

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
func Dump() string {
	var lines []string

	for _, v := range mapPathEntry {
		switch v.kind {
		case 'i':
			lines = append(lines, fmt.Sprintf("(%d/%d) \t%-20s \t\"%d\"", v.refGet, v.refSet, v.path, v.vInt))

		case 's':
			lines = append(lines, fmt.Sprintf("(%d/%d) \t%-20s \t\"%s\"", v.refGet, v.refSet, v.path, v.vStr))

		case 'b':
			lines = append(lines, fmt.Sprintf("(%d/%d) \t%-20s \t\"%t\"", v.refGet, v.refSet, v.path, v.vBool))
		}
	}

	return strings.Join(lines, "\n")
}

func init() {
	cfgList := os.Getenv("KCFG_FILES")
	files := strings.Split(cfgList, ":")
	for _, f := range files {
		if f != "" {
			if err := Load(f); err != nil {
				klog.E("LOAD KCFG_FILES Error: %s", err.Error())
			}
		}
	}

	for _, argv := range os.Args {
		if strings.HasPrefix(argv, "--kfg=") {
			f := argv[6:]
			if f != "" {
				if err := Load(f); err != nil {
					klog.E("LOAD --kfg=xxx Error: %s", err.Error())
				}
			}
		}
	}
}
