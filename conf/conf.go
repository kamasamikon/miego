package conf

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
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
	//
	// Safe: Not show by dump
	kind byte
	path string

	vInt  int64
	vStr  string
	vBool bool

	refGet int64
	refSet int64

	safe bool
}

// KConfMonitor is a Callback called when wathed entry modified.
type KConfMonitor func(path string, oVal interface{}, nVal interface{})

type confMonitors struct {
	funcList []KConfMonitor
}

var mapPathEntry = make(map[string]*confEntry)
var mapPathMonitors = make(map[string]*confMonitors)

// Load file from configure
func EntryAdd(line string) {
	line = strings.TrimSpace(line)
	segs := strings.SplitN(line, "=", 2)
	if len(segs) < 2 {
		return
	}

	path, value := segs[0], segs[1]
	if path[1] != ':' {
		return
	}

	kind, safe, realpath := pathParse(path)
	if realpath == "" {
		return
	}

	e := confEntry{
		kind: kind,
		safe: safe,
		path: realpath,
	}

	switch kind {
	case 'i':
		if vInt, err := strconv.ParseInt(value, 10, 64); err == nil {
			e.vInt = vInt
		} else {
			return
		}

	case 's':
		e.vStr = value

	case 'b':
		// true: 1, t, T
		// false: 0, f, F
		x := value[0]
		if x == '1' || x == 't' || x == 'T' || x == 'y' || x == 'Y' {
			e.vBool = true
		} else if x == '0' || x == 'f' || x == 'F' || x == 'n' || x == 'N' {
			e.vBool = false
		} else {
			return
		}

	default:
		return
	}

	// Simply overwrite the old value.
	mapPathEntry[e.path] = &e
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
		if line, err := buf.ReadString('\n'); err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		} else {
			EntryAdd(line)
		}
	}
}

// Has : Check if a entry exists
func Has(path string) bool {
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

// List : get a List entry
func List(sep string, paths ...string) []string {
	// path: aaa/bbb
	for _, path := range paths {
		key := "s:/" + path
		if v, ok := mapPathEntry[key]; ok {
			v.refGet++
			return strings.Split(v.vStr, sep)
		}
	}
	return nil
}

func pathParse(path string) (kind byte, safe bool, realpath string) {
	switch path[0] {
	case 'i', 'I':
		kind = 'i'
		safe = path[0] == 'I'
		realpath = "i" + path[1:]

	case 's', 'S':
		kind = 's'
		safe = path[0] == 'S'
		realpath = "s" + path[1:]

	case 'b', 'B':
		kind = 'b'
		safe = path[0] == 'B'
		realpath = "b" + path[1:]

	default:
		realpath = ""
	}

	return kind, safe, realpath
}

func Names() []string {
	var names []string
	for k, _ := range mapPathEntry {
		names = append(names, k)
	}
	return names
}

// Set : Modify or Add conf entry
func Set(path string, value interface{}, force bool) {
	var e *confEntry
	var ok bool

	kind, safe, realpath := pathParse(path)
	if realpath == "" {
		return
	}

	// FIXME: Set("B:/aaa/bbb", ...) will makes query failed.
	e, ok = mapPathEntry[realpath]
	if !ok {
		if force {
			e = &confEntry{
				kind: kind,
				safe: safe,
				path: realpath,
			}
			mapPathEntry[e.path] = e
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
	var cList []*confEntry
	for _, v := range mapPathEntry {
		cList = append(cList, v)
	}
	sort.Slice(cList, func(i int, j int) bool {
		return strings.Compare(cList[i].path[1:], cList[j].path[1:]) < 0
	})

	var lines []string
	for _, v := range cList {
		if v.safe {
			continue
		}

		switch v.kind {
		case 'i':
			lines = append(lines, fmt.Sprintf("(%d/%d) \t%-20s \t%d", v.refGet, v.refSet, v.path, v.vInt))

		case 's':
			lines = append(lines, fmt.Sprintf("(%d/%d) \t%-20s \t\"%s\"", v.refGet, v.refSet, v.path, v.vStr))

		case 'b':
			lines = append(lines, fmt.Sprintf("(%d/%d) \t%-20s \t%t", v.refGet, v.refSet, v.path, v.vBool))
		}
	}

	// Add the last \n
	lines = append(lines, "")

	return strings.Join(lines, "\n")
}

func DumpRaw() string {
	var cList []*confEntry
	for _, v := range mapPathEntry {
		cList = append(cList, v)
	}
	sort.Slice(cList, func(i int, j int) bool {
		return strings.Compare(cList[i].path[1:], cList[j].path[1:]) < 0
	})

	var lines []string
	for _, v := range cList {
		if v.safe {
			continue
		}

		switch v.kind {
		case 'i':
			lines = append(lines, fmt.Sprintf("%s=%d", v.refGet, v.refSet, v.path, v.vInt))

		case 's':
			lines = append(lines, fmt.Sprintf("%s=%s", v.refGet, v.refSet, v.path, v.vStr))

		case 'b':
			lines = append(lines, fmt.Sprintf("%s=%t", v.refGet, v.refSet, v.path, v.vBool))
		}
	}

	// Add the last \n
	lines = append(lines, "")

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
	for _, argv := range os.Args {
		if strings.HasPrefix(argv, "--kfg-item=") {
			item := argv[11:]
			EntryAdd(item)
		}
	}
}
