package conf

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync"
)

const (
	ColorType_F     = "\x1b[1;31;40m"
	ColorType_A     = "\x1b[91;40m"
	ColorType_C     = "\x1b[1;36;40m"
	ColorType_E     = "\x1b[96;40m"
	ColorType_W     = "\x1b[1;33;40m"
	ColorType_N     = "\x1b[93;40m"
	ColorType_I     = "\x1b[1;32;40m"
	ColorType_D     = "\x1b[92;40m"
	ColorType_Reset = "\x1b[0m"
)

const (
	PathReady = "e:/conf/ready"
)

// See confcenter
type confEntry struct {
	// kind: a:arr, b:bool, d:dat(len+dat), e:event, i:int, s:str, p:ptr
	//
	// path: i:/aaa/bbb; b:/xxx/zzz
	//
	// vXxx: value for each type
	//
	//
	// hidden: Not show by dump
	kind byte
	path string

	vInt  int64
	vStr  string
	vBool bool
	vObj  interface{}

	hidden bool
}

var mapPathEntry = make(map[string]*confEntry)
var LoadOKCount = 0
var LoadNGCount = 0

var mutex = &sync.Mutex{}

// Delete an entry
func EntryRem(path string) {
	mutex.Lock()
	defer mutex.Unlock()

	delete(mapPathEntry, path)
}

// Load file from configure
func EntryAdd(line string, overwrite bool) {
	mutex.Lock()
	defer mutex.Unlock()

	line = strings.TrimSpace(line)

	segs := strings.SplitN(line, "=", 2)
	if len(segs) == 0 {
		return
	}

	path := segs[0]
	if len(path) < 4 || path[1] != ':' {
		return
	}
	kind, hidden, realpath := pathParse(path)
	if realpath == "" {
		return
	}

	var value string
	if len(segs) == 2 {
		value = segs[1]
	}

	e, exists := mapPathEntry[path]
	if exists && !overwrite {
		return
	}

	var vNew interface{}

	switch kind {
	case 'i':
		if vInt, err := strconv.ParseInt(value, 10, 64); err == nil {
			vNew = vInt
		} else {
			return
		}

	case 's':
		vNew = value

	case 'b':
		// true: 1, t, T
		// false: 0, f, F
		x := value[0]
		if x == '1' || x == 't' || x == 'T' || x == 'y' || x == 'Y' {
			vNew = true
		} else if x == '0' || x == 'f' || x == 'F' || x == 'n' || x == 'N' {
			vNew = false
		} else {
			return
		}

	case 'o':
		// line is json string
		o := make(map[string]interface{})
		if err := json.Unmarshal([]byte(value), &o); err != nil {
			return
		}
		vNew = o

	case 'e':
		// event, no data at all, value treated as a parameter
		vNew = value

	default:
		return
	}

	if !exists {
		e = &confEntry{
			kind:   kind,
			hidden: hidden,
			path:   realpath,
		}
		mapPathEntry[e.path] = e
	}
	setByEntry(e, vNew)
}

// LoadString : Load setting from string (lines of configuration)
func LoadString(s string, overwrite bool) {
	s = strings.Replace(s, "\r", "\n", -1)
	Lines := strings.Split(s, "\n")
	for _, Line := range Lines {
		EntryAdd(Line, overwrite)
	}
}

// Load : configure from a file.
func LoadFile(fileName string, overwrite bool) error {
	const (
		NGName = "s:/conf/Load/NG/%d/Name=%s"
		NGWhy  = "s:/conf/Load/NG/%d/Why=%s"
		OKName = "s:/conf/Load/OK/%d=%s"
	)

	sp := fmt.Sprintf
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		EntryAdd(sp(NGName, LoadNGCount, fileName), false)
		EntryAdd(sp(NGWhy, LoadNGCount, err.Error()), false)
		LoadNGCount++
		return err
	}

	EntryAdd(sp(OKName, LoadOKCount, fileName), false)
	LoadOKCount++

	LoadString(string(data), overwrite)
	return nil
}

// Has : Check if a entry exists
func Has(path string) bool {
	mutex.Lock()
	defer mutex.Unlock()

	_, ok := mapPathEntry[path]
	return ok
}

// Int : get a int typed configure
func Int(defval int64, paths ...string) int64 {
	mutex.Lock()
	defer mutex.Unlock()

	// path: aaa/bbb
	for _, path := range paths {
		if v, ok := mapPathEntry[path]; ok {
			return v.vInt
		}
	}
	return defval
}

// Int : get a int typed configure
func IntX(paths ...string) (int64, bool) {
	mutex.Lock()
	defer mutex.Unlock()

	// path: aaa/bbb
	for _, path := range paths {
		if v, ok := mapPathEntry[path]; ok {
			return v.vInt, true
		}
	}
	return 0, false
}

// Inc : Increase or Decrease on int
func Inc(inc int64, path string) {
	mutex.Lock()
	defer mutex.Unlock()

	if e, ok := mapPathEntry[path]; ok {
		vNew := e.vInt + 1
		setByEntry(e, vNew)
	}
}

// Flip : flip on bool
func Flip(path string) {
	mutex.Lock()
	defer mutex.Unlock()

	if e, ok := mapPathEntry[path]; ok {
		vNew := !e.vBool
		setByEntry(e, vNew)
	}
}

// Str : get a str typed configure
func Str(defval string, paths ...string) string {
	mutex.Lock()
	defer mutex.Unlock()

	// path: aaa/bbb
	for _, path := range paths {
		if v, ok := mapPathEntry[path]; ok {
			return v.vStr
		}
	}
	return defval
}

// Str : get a str typed configure
func StrX(paths ...string) (string, bool) {
	mutex.Lock()
	defer mutex.Unlock()

	// path: aaa/bbb
	for _, path := range paths {
		if v, ok := mapPathEntry[path]; ok {
			return v.vStr, true
		}
	}
	return "", false
}

// Bool : get a bool entry
func Bool(defval bool, paths ...string) bool {
	mutex.Lock()
	defer mutex.Unlock()

	// path: aaa/bbb
	for _, path := range paths {
		if v, ok := mapPathEntry[path]; ok {
			return v.vBool
		}
	}
	return defval
}

// Bool : get a bool entry
func BoolX(paths ...string) (bool, bool) {
	mutex.Lock()
	defer mutex.Unlock()

	// path: aaa/bbb
	for _, path := range paths {
		if v, ok := mapPathEntry[path]; ok {
			return v.vBool, true
		}
	}
	return false, false
}

// Object : get a bool entry
func Obj(defval interface{}, paths ...string) interface{} {
	mutex.Lock()
	defer mutex.Unlock()

	// path: aaa/bbb
	for _, path := range paths {
		if v, ok := mapPathEntry[path]; ok {
			return v.vObj
		}
	}
	return defval
}

// Object : get a bool entry
func ObjX(paths ...string) (interface{}, bool) {
	mutex.Lock()
	defer mutex.Unlock()

	// path: aaa/bbb
	for _, path := range paths {
		if v, ok := mapPathEntry[path]; ok {
			return v.vObj, true
		}
	}
	return nil, false
}

// List : get a List entry. s:/names=:aaa:bbb first char is the seperator
func List(paths ...string) []string {
	mutex.Lock()
	defer mutex.Unlock()

	// path: aaa/bbb
	var slice []string
	for _, path := range paths {
		if v, ok := mapPathEntry[path]; ok {
			if len(v.vStr) > 0 {
				for _, s := range strings.Split(v.vStr, v.vStr[0:1]) {
					if s != "" {
						slice = append(slice, s)
					}
				}
			}
		}
	}
	return slice
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

	case 'e', 'E':
		kind = 'e'
		hidden = path[0] == 'E'
		realpath = "e" + path[1:]

	default:
		realpath = ""
	}

	return kind, hidden, realpath
}

// Names : All Keys
func Names() []string {
	mutex.Lock()
	defer mutex.Unlock()

	var names []string
	for k, _ := range mapPathEntry {
		names = append(names, k)
	}
	return names
}

// SafeNames : Keys not hidden
func SafeNames() []string {
	mutex.Lock()
	defer mutex.Unlock()

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

func setByEntry(e *confEntry, value interface{}) {
	switch e.kind {
	case 'i':
		vOld := e.vInt

		var vNew int64
		switch value.(type) {
		case int64:
			vNew = int64(value.(int64))
		case int32:
			vNew = int64(value.(int32))
		case int:
			vNew = int64(value.(int))
		case int16:
			vNew = int64(value.(int16))
		case int8:
			vNew = int64(value.(int8))
		case uint64:
			vNew = int64(value.(uint64))
		case uint32:
			vNew = int64(value.(uint32))
		case uint:
			vNew = int64(value.(uint))
		case uint16:
			vNew = int64(value.(uint16))
		case uint8:
			vNew = int64(value.(uint8))
		}

		e.vInt = vNew
		monitorCall(e, vOld, vNew)

	case 's':
		vOld := e.vStr

		vNew := value.(string)
		e.vStr = vNew
		monitorCall(e, vOld, vNew)

	case 'b':
		vOld := e.vBool

		vNew := value.(bool)
		e.vBool = vNew
		monitorCall(e, vOld, vNew)

	case 'o':
		vOld := e.vObj

		vNew := value.(interface{})
		e.vObj = vNew
		monitorCall(e, vOld, vNew)

	case 'e':
		vNew := value.(string)
		monitorCall(e, 0, vNew)
	}
}

// Set : Modify or Add conf entry
func Set(path string, value interface{}, force bool) {
	mutex.Lock()
	defer mutex.Unlock()

	var e *confEntry
	var ok bool

	kind, hidden, realpath := pathParse(path)
	if realpath == "" {
		return
	}

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
			return
		}
	}

	setByEntry(e, value)
}

func Ready() {
	Set(PathReady, "", true)
}

//go:embed main.mgc
var main_cfg string

func init() {
	//
	// Some builtin entries
	//
	EntryAdd(PathReady, false)
	LoadString(main_cfg, false)

	// 优先级: 命令行 > 环境变量
	LoadFromEnv()
	LoadFromArg()

	// Load environment to conf
	for _, env := range os.Environ() {
		segs := strings.SplitN(env, "=", 2)
		if len(segs) == 2 {
			Set("s:/env/"+segs[0], segs[1], true)
		}
	}
}

func LoadFromEnv() {
	{
		cfgList := os.Getenv("KCFG_FILES")
		files := strings.Split(cfgList, ":")
		for _, f := range files {
			if f != "" {
				LoadFile(f, true)
			}
		}
	}
	{
		cfgList := os.Getenv("KCFG_QQQ_FILES")
		files := strings.Split(cfgList, ":")
		for _, f := range files {
			if f != "" {
				LoadFile(f, true)
				os.Remove(f)
			}
		}
	}
}

func LoadFromArg() {
	{
		for _, argv := range os.Args {
			if strings.HasPrefix(argv, "--kfg=") {
				f := argv[6:]
				if f != "" {
					LoadFile(f, true)
				}
			}
		}
	}
	{
		for _, argv := range os.Args {
			if strings.HasPrefix(argv, "--kfg-qqq=") {
				f := argv[6:]
				if f != "" {
					LoadFile(f, true)
					os.Remove(f)
				}
			}
		}
	}

	for _, argv := range os.Args {
		if strings.HasPrefix(argv, "--kfg-item=") {
			item := argv[11:]
			EntryAdd(item, true)
		}
	}
}

// Last to call
func Go() {
	Ready()
}
