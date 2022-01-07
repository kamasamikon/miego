package conf

import (
	"bufio"
	"encoding/json"
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
	// hidden: Not show by dump
	kind byte
	path string

	vInt  int64
	vStr  string
	vBool bool
	vObj  interface{}

	refGet int64
	refSet int64

	hidden bool
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

	kind, hidden, realpath := pathParse(path)
	if realpath == "" {
		return
	}

	e, ok := mapPathEntry[path]
	if !ok {
		e = &confEntry{
			kind:   kind,
			hidden: hidden,
			path:   realpath,
		}
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

	case 'o':
		// line is json string
		o := make(map[string]interface{})
		if err := json.Unmarshal([]byte(value), &o); err != nil {
			klog.E(err.Error())
		}
		e.vObj = o

	default:
		return
	}

	if ok {
		e.refSet++
	}
	// Simply overwrite the old value.
	mapPathEntry[e.path] = e
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

// Ref : refGet, refSet
func Ref(path string) (int64, int64) {
	if v, ok := mapPathEntry[path]; ok {
		return v.refGet, v.refSet
	}
	return -1, -1
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

// Inc : Increase or Decrease on int
func Inc(inc int64, path string) {
	if v, ok := mapPathEntry[path]; ok {
		v.refSet++
		v.vInt += inc
	}
}

// Flip : flip on bool
func Flip(path string) {
	if v, ok := mapPathEntry[path]; ok {
		v.refSet++
		v.vBool = !v.vBool
	}
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

// Object : get a bool entry
func Obj(defval interface{}, paths ...string) interface{} {
	// path: aaa/bbb
	for _, path := range paths {
		key := "o:/" + path
		if v, ok := mapPathEntry[key]; ok {
			v.refGet++
			return v.vObj
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

func pathParse(path string) (kind byte, hidden bool, realpath string) {
	switch path[0] {
	case 'i', 'I':
		kind = 'i'
		hidden = path[0] == 'I'
		realpath = "i" + path[1:]

	case 's', 'S':
		kind = 's'
		hidden = path[0] == 'S'
		realpath = "s" + path[1:]

	case 'b', 'B':
		kind = 'b'
		hidden = path[0] == 'B'
		realpath = "b" + path[1:]

	case 'o', 'O':
		kind = 'o'
		hidden = path[0] == 'O'
		realpath = "o" + path[1:]

	default:
		realpath = ""
	}

	return kind, hidden, realpath
}

// Names : All Keys
func Names() []string {
	var names []string
	for k, _ := range mapPathEntry {
		names = append(names, k)
	}
	return names
}

// SafeNames : Keys not hidden
func SafeNames() []string {
	var names []string
	for k, e := range mapPathEntry {
		if !e.hidden {
			names = append(names, k)
		}
	}
	return names
}

func Add(path string, value interface{}) {
	Set(path, value, true)
}

// Set : Modify or Add conf entry
func Set(path string, value interface{}, force bool) {
	var e *confEntry
	var ok bool

	kind, hidden, realpath := pathParse(path)
	if realpath == "" {
		return
	}

	// FIXME: Set("B:/aaa/bbb", ...) will makes query failed.
	e, ok = mapPathEntry[realpath]
	if !ok {
		if force {
			e = &confEntry{
				kind:   kind,
				hidden: hidden,
				path:   realpath,
			}
			mapPathEntry[e.path] = e
		} else {
			klog.D("path:%s and force:false", path)
			return
		}
	}

	switch kind {
	case 'i':
		vOld := e.vInt

		var vNew int64
		if x, ok := value.(int64); ok {
			vNew = int64(x)
		} else if x, ok := value.(int32); ok {
			vNew = int64(x)
		} else if x, ok := value.(int); ok {
			vNew = int64(x)
		} else if x, ok := value.(int16); ok {
			vNew = int64(x)
		} else if x, ok := value.(int8); ok {
			vNew = int64(x)
		} else if x, ok := value.(uint64); ok {
			vNew = int64(x)
		} else if x, ok := value.(uint32); ok {
			vNew = int64(x)
		} else if x, ok := value.(uint); ok {
			vNew = int64(x)
		} else if x, ok := value.(uint16); ok {
			vNew = int64(x)
		} else if x, ok := value.(uint8); ok {
			vNew = int64(x)
		}
		e.vInt = vNew
		e.refSet++
		monitorCall(e, vOld, vNew)

	case 's':
		vOld := e.vStr

		vNew := value.(string)
		e.vStr = vNew
		e.refSet++
		monitorCall(e, vOld, vNew)

	case 'b':
		vOld := e.vBool

		vNew := value.(bool)
		e.vBool = vNew
		e.refSet++
		monitorCall(e, vOld, vNew)

	case 'o':
		vOld := e.vObj

		vNew := value.(interface{})
		e.vObj = vNew
		e.refSet++
		monitorCall(e, vOld, vNew)

	default:
		klog.E("Bad Kind: %c", kind)
	}
}

// Monitor : Callback AFTER entry changed.
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
		if v.hidden {
			continue
		}

		switch v.kind {
		case 'i':
			lines = append(lines, fmt.Sprintf("(%04d/%04d) \t%-20s \t%d", v.refGet, v.refSet, v.path, v.vInt))

		case 's':
			lines = append(lines, fmt.Sprintf("(%04d/%04d) \t%-20s \t\"%s\"", v.refGet, v.refSet, v.path, v.vStr))

		case 'b':
			lines = append(lines, fmt.Sprintf("(%04d/%04d) \t%-20s \t%t", v.refGet, v.refSet, v.path, v.vBool))

		case 'o':
			lines = append(lines, fmt.Sprintf("(%04d/%04d) \t%-20s \t%s", v.refGet, v.refSet, v.path, "..."))
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
		if v.hidden {
			continue
		}

		switch v.kind {
		case 'i':
			lines = append(lines, fmt.Sprintf("%s=%d", v.path, v.vInt))

		case 's':
			lines = append(lines, fmt.Sprintf("%s=%s", v.path, v.vStr))

		case 'b':
			lines = append(lines, fmt.Sprintf("%s=%t", v.path, v.vBool))

		case 'o':
			lines = append(lines, fmt.Sprintf("%s=%s", v.path, "..."))
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
