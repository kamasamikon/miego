package conf

import (
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

const (
	ColorTypeF     = "\x1b[1;31;40m"
	ColorTypeA     = "\x1b[91;40m"
	ColorTypeC     = "\x1b[1;36;40m"
	ColorTypeE     = "\x1b[96;40m"
	ColorTypeW     = "\x1b[1;33;40m"
	ColorTypeN     = "\x1b[93;40m"
	ColorTypeI     = "\x1b[1;32;40m"
	ColorTypeD     = "\x1b[92;40m"
	ColorTypeReset = "\x1b[0m"
)

const (
	PathReady     = "e:/conf/ready"
	Debug         = "i:/conf/debug"
	MissedEntries = "s:/conf/missedEntries"
)

type setter func(e *confEntry, v any) (vv any, ok bool) // e.v = vv if ok
type getter func(e *confEntry) (vv any, ok bool)        // return vv if ok else e.v

// See confcenter
type confEntry struct {
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
	vObj  any

	setter setter
	getter getter

	refGet int64
	refSet int64

	hidden bool
}

//go:embed assets/*
var Assets embed.FS

var mapPathEntry = make(map[string]*confEntry)
var LoadOKCount = 0
var LoadNGCount = 0

// path <> ref
var mapMissedEntries = make(map[string]int)

var mutex = &sync.Mutex{}

// 打印调试信息
var DEBUG = 0

func setMissedEntries(path string) {
	mapMissedEntries[path] = 1
}

// DebugPrint: 打印调试信息
func dp(formating string, args ...any) {
	if DEBUG == 0 {
		return
	}

	pc, filename, line, _ := runtime.Caller(1)

	funcname := runtime.FuncForPC(pc).Name()
	funcname = filepath.Ext(funcname)
	funcname = strings.TrimPrefix(funcname, ".")

	var sb strings.Builder
	sb.WriteRune('|')
	sb.WriteString(filename)
	sb.WriteRune('|')
	sb.WriteString(funcname)
	sb.WriteRune('|')
	sb.WriteString(strconv.Itoa(line))
	sb.WriteRune('|')
	sb.WriteRune(' ')
	sb.WriteString("\x1b[31;40m")
	sb.WriteString(fmt.Sprintf(formating, args...))
	sb.WriteString("\x1b[0m")
	sb.WriteRune('\n')

	fmt.Printf("%s", sb.String())
}

// Delete an entry
func EntryRem(path string) {
	mutex.Lock()
	defer mutex.Unlock()

	delete(mapPathEntry, path)
}

func EntryAddByLine(line string, overwrite bool) {
	mutex.Lock()
	defer mutex.Unlock()

	segs := strings.SplitN(strings.TrimSpace(line), "=", 2)
	if len(segs) < 2 {
		return
	}

	path, value := segs[0], segs[1]
	EntryAdd(path, value, overwrite)
}

func EntryAdd(path string, value string, overwrite bool) {
	kind, hidden, realpath := pathParse(path)
	if realpath == "" {
		return
	}

	e, exists := mapPathEntry[path]
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
			dp("BadValue.b: %s", value)
			return
		}

	case 'o':
		// line is json string
		o := make(map[string]any)
		if err := json.Unmarshal([]byte(value), &o); err != nil {
			dp("BadValue.o: %s", value)
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

func SetSetter(path string, setter setter) {
	if e, ok := mapPathEntry[path]; ok {
		e.setter = setter
	}
}

func SetGetter(path string, getter getter) {
	if e, ok := mapPathEntry[path]; ok {
		e.getter = getter
	}
}

// Load : configure from a file.
func LoadFromText(text string, overwrite bool) {
	Lines := strings.Split(strings.Replace(text, "\r", "\n", -1), "\n")
	size := len(Lines)
	i := 0
	for {
		if i >= size {
			break
		}

		Line := Lines[i]
		i++

		neat := strings.TrimSpace(Line)
		if neat == "" || neat[0] == '#' {
			continue
		}
		segs := strings.SplitN(neat, "=", 2)
		if len(segs) < 2 {
			continue
		}
		path, value := segs[0], segs[1]
		if len(path) < 4 || path[1] != ':' {
			continue
		}

		if strings.HasPrefix(value, "<<") {
			multiLineTag := value[2:]

			var sb strings.Builder
			for {
				if i >= size {
					break
				}

				Line = Lines[i]
				i++

				if Line == multiLineTag {
					break
				}

				sb.WriteString(Line)
				sb.WriteRune('\n')
			}

			EntryAdd(path, sb.String(), overwrite)
		} else {
			EntryAdd(path, value, overwrite)
		}
	}
}

// Load : configure from a file.
func LoadFromFile(fileName string, overwrite bool) error {
	const (
		NGName = "s:/conf/Load/NG/%d/Name=%s"
		NGWhy  = "s:/conf/Load/NG/%d/Why=%s"
		OKName = "s:/conf/Load/OK/%d=%s"
	)

	data, err := os.ReadFile(fileName)
	if err != nil {
		EntryAddByLine(fmt.Sprintf(NGName, LoadNGCount, fileName), false)
		EntryAddByLine(fmt.Sprintf(NGWhy, LoadNGCount, err.Error()), false)
		LoadNGCount++
		dp("Error:'%s', fileName:'%s'", err.Error(), fileName)
		return err
	}

	EntryAddByLine(fmt.Sprintf(OKName, LoadOKCount, fileName), false)
	LoadOKCount++

	LoadFromText(string(data), overwrite)
	return nil
}

// Ref : refGet, refSet
func Ref(path string) (int64, int64) {
	if e, ok := mapPathEntry[path]; ok {
		return e.refGet, e.refSet
	}
	return -1, -1
}

// Has : Check if a entry exists
func Has(path string) bool {
	mutex.Lock()
	defer mutex.Unlock()

	e, ok := mapPathEntry[path]
	if ok {
		e.refGet++
	}
	return ok
}

// Int : get a int typed configure
func Int(defval int64, paths ...string) int64 {
	mutex.Lock()
	defer mutex.Unlock()

	// path: aaa/bbb
	for _, path := range paths {
		if e, ok := mapPathEntry[path]; ok {
			e.refGet++
			if e.getter != nil {
				if vv, ok := e.getter(e); ok {
					return vv.(int64)
				}
			}
			return e.vInt
		}
		dp("Miss %s", path)
		setMissedEntries(path)
	}
	return defval
}

// Int : get a int typed configure
func IntX(paths ...string) (int64, bool) {
	mutex.Lock()
	defer mutex.Unlock()

	// path: aaa/bbb
	for _, path := range paths {
		if e, ok := mapPathEntry[path]; ok {
			e.refGet++
			if e.getter != nil {
				if vv, ok := e.getter(e); ok {
					return vv.(int64), true
				}
			}
			return e.vInt, true
		}
		dp("Miss %s", path)
		setMissedEntries(path)
	}
	return 0, false
}

// Inc : Increase or Decrease on int
func Inc(inc int64, path string) int64 {
	mutex.Lock()
	defer mutex.Unlock()

	if e, ok := mapPathEntry[path]; ok {
		e.refGet++
		vInt := e.vInt
		if e.getter != nil {
			if vv, ok := e.getter(e); ok {
				vInt = vv.(int64)
			}
		}
		vNew := vInt + inc
		setByEntry(e, vNew)
		return vNew
	}
	dp("Miss %s", path)
	setMissedEntries(path)
	return -1
}

// Flip : flip on bool
func Flip(path string) bool {
	mutex.Lock()
	defer mutex.Unlock()

	if e, ok := mapPathEntry[path]; ok {
		e.refGet++
		vBool := e.vBool
		if e.getter != nil {
			if vv, ok := e.getter(e); ok {
				vBool = vv.(bool)
			}
		}
		vNew := !vBool
		setByEntry(e, vNew)
		return vNew
	}
	dp("Miss %s", path)
	setMissedEntries(path)
	return false
}

// Str : get a str typed configure
func Str(defval string, paths ...string) string {
	mutex.Lock()
	defer mutex.Unlock()

	// path: aaa/bbb
	for _, path := range paths {
		if e, ok := mapPathEntry[path]; ok {
			e.refGet++
			if e.getter != nil {
				if vv, ok := e.getter(e); ok {
					return vv.(string)
				}
			}
			return e.vStr
		}
		dp("Miss %s", path)
		setMissedEntries(path)
	}
	return defval
}

// Str : get a str typed configure
func StrX(paths ...string) (string, bool) {
	mutex.Lock()
	defer mutex.Unlock()

	// path: aaa/bbb
	for _, path := range paths {
		if e, ok := mapPathEntry[path]; ok {
			e.refGet++
			if e.getter != nil {
				if vv, ok := e.getter(e); ok {
					return vv.(string), true
				}
			}
			return e.vStr, true
		}
		dp("Miss %s", path)
		setMissedEntries(path)
	}
	return "", false
}

// Bool : get a bool entry
func Bool(defval bool, paths ...string) bool {
	mutex.Lock()
	defer mutex.Unlock()

	// path: aaa/bbb
	for _, path := range paths {
		if e, ok := mapPathEntry[path]; ok {
			e.refGet++
			if e.getter != nil {
				if vv, ok := e.getter(e); ok {
					return vv.(bool)
				}
			}
			return e.vBool
		}
		dp("Miss %s", path)
		setMissedEntries(path)
	}
	return defval
}

// Bool : get a bool entry
func BoolX(paths ...string) (bool, bool) {
	mutex.Lock()
	defer mutex.Unlock()

	// path: aaa/bbb
	for _, path := range paths {
		if e, ok := mapPathEntry[path]; ok {
			e.refGet++
			if e.getter != nil {
				if vv, ok := e.getter(e); ok {
					return vv.(bool), true
				}
			}
			return e.vBool, true
		}
		dp("Miss %s", path)
		setMissedEntries(path)
	}
	return false, false
}

// Object : get a bool entry
func Obj(defval any, paths ...string) any {
	mutex.Lock()
	defer mutex.Unlock()

	// path: aaa/bbb
	for _, path := range paths {
		if e, ok := mapPathEntry[path]; ok {
			e.refGet++
			if e.getter != nil {
				if vv, ok := e.getter(e); ok {
					return vv
				}
			}
			return e.vObj
		}
		dp("Miss %s", path)
		setMissedEntries(path)
	}
	return defval
}

// Object : get a bool entry
func ObjX(paths ...string) (any, bool) {
	mutex.Lock()
	defer mutex.Unlock()

	// path: aaa/bbb
	for _, path := range paths {
		if e, ok := mapPathEntry[path]; ok {
			e.refGet++
			if e.getter != nil {
				if vv, ok := e.getter(e); ok {
					return vv, true
				}
			}
			return e.vObj, true
		}
		dp("Miss %s", path)
		setMissedEntries(path)
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
		if e, ok := mapPathEntry[path]; ok {
			if len(e.vStr) > 0 {
				e.refGet++
				vStr := e.vStr
				if e.getter != nil {
					if vv, ok := e.getter(e); ok {
						vStr = vv.(string)
					}
				}
				for _, s := range strings.Split(vStr, vStr[0:1]) {
					if s != "" {
						slice = append(slice, s)
					}
				}
			}
		}
		dp("Miss %s", path)
		setMissedEntries(path)
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
		dp("BadType: %s", path)
		realpath = ""
	}

	return kind, hidden, realpath
}

// Names : All Keys
func Names() []string {
	mutex.Lock()
	defer mutex.Unlock()

	var names []string
	for k := range mapPathEntry {
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

func Add(path string, value any) {
	Set(path, value, true)
}

func setByEntry(e *confEntry, value any) {
	switch e.kind {
	case 'i':
		vOld := e.vInt

		var vNew int64

		switch v := value.(type) {
		case int64:
			vNew = int64(v)
		case int32:
			vNew = int64(v)
		case int:
			vNew = int64(v)
		case int16:
			vNew = int64(v)
		case int8:
			vNew = int64(v)
		case uint64:
			vNew = int64(v)
		case uint32:
			vNew = int64(v)
		case uint:
			vNew = int64(v)
		case uint16:
			vNew = int64(v)
		case uint8:
			vNew = int64(v)
		}

		if e.setter != nil {
			if vv, ok := e.setter(e, vNew); ok {
				vNew = vv.(int64)
			}
		}
		e.vInt = vNew
		e.refSet++
		monitorCall(e, vOld, vNew)

	case 's':
		vOld := e.vStr

		vNew := value.(string)
		if e.setter != nil {
			if vv, ok := e.setter(e, vNew); ok {
				vNew = vv.(string)
			}
		}
		e.vStr = vNew
		e.refSet++
		monitorCall(e, vOld, vNew)

	case 'b':
		vOld := e.vBool

		vNew := value.(bool)
		if e.setter != nil {
			if vv, ok := e.setter(e, vNew); ok {
				vNew = vv.(bool)
			}
		}
		e.vBool = vNew
		e.refSet++
		monitorCall(e, vOld, vNew)

	case 'o':
		vOld := e.vObj

		vNew := value
		if e.setter != nil {
			if vv, ok := e.setter(e, vNew); ok {
				vNew = vv
			}
		}
		e.vObj = vNew
		e.refSet++
		monitorCall(e, vOld, vNew)

	case 'e':
		vNew := value.(string)
		if e.setter != nil {
			if vv, ok := e.setter(e, vNew); ok {
				vNew = vv.(string)
			}
		}
		e.refSet++
		monitorCall(e, 0, vNew)
	}
}

// Set : Modify or Add conf entry
func Set(path string, value any, create bool) {
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
		if create {
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

func init() {
	data, err := Assets.ReadFile("assets/main.cfg")
	if err == nil {
		LoadFromText(string(data), false)
	}

	//
	// Some builtin entries
	//
	Set(PathReady, "", true)
	Set(MissedEntries, "", true)

	SetGetter(MissedEntries, func(_ *confEntry) (vv any, ok bool) {
		var sb strings.Builder
		for p := range mapMissedEntries {
			sb.WriteString(p)
			sb.WriteRune(';')
		}
		return sb.String(), true
	})

	if os.Getenv("MG_CONF_DEBUG") == "debug" {
		Set(Debug, 1, true)
		DEBUG = 1
	} else {
		Set(Debug, 0, true)
		DEBUG = 0
	}

	SetSetter(Debug, func(_ *confEntry, v any) (vv any, ok bool) {
		DEBUG = v.(int)
		return nil, false
	})

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
				LoadFromFile(f, true)
			}
		}
	}
	{
		cfgList := os.Getenv("KCFG_QQQ_FILES")
		files := strings.Split(cfgList, ":")
		for _, f := range files {
			if f != "" {
				LoadFromFile(f, true)
				os.Remove(f)
			}
		}
	}
}

func LoadFromArg() {
	argc := len(os.Args)

	// --kfg abc.cfg --kfg=xyz.cfg
	{
		for i, argv := range os.Args {
			if argv == "--kfg" {
				i++
				if i < argc {
					LoadFromFile(os.Args[i], true)
				}
				continue
			}
			if strings.HasPrefix(argv, "--kfg=") {
				f := argv[6:]
				if f != "" {
					LoadFromFile(f, true)
				}
				continue
			}
		}
	}

	// --kfg-qqq abc.cfg --kfg-qqq=xyz.cfg
	{
		for i, argv := range os.Args {
			if argv == "--kfg-qqq" {
				i++
				if i < argc {
					LoadFromFile(os.Args[i], true)
					os.Remove(os.Args[i])
				}
				continue
			}
			if strings.HasPrefix(argv, "--kfg-qqq=") {
				f := argv[6:]
				if f != "" {
					LoadFromFile(f, true)
					os.Remove(f)
				}
			}
		}
	}

	// --kfg-item i:/abc=777 --kfg-item=s:/xyz=abc
	for i, argv := range os.Args {
		if argv == "--kfg-item" {
			i++
			if i < argc {
				EntryAddByLine(os.Args[i], true)
			}
			continue
		}
		if strings.HasPrefix(argv, "--kfg-item=") {
			item := argv[11:]
			EntryAddByLine(item, true)
		}
	}
}

// ///////////////////////////////////////////////////////////////////////
// OnReady : Called when all configure loaded.
var onReadys []func()

func OnReady(cb func()) {
	mutex.Lock()
	defer mutex.Unlock()
	onReadys = append(onReadys, cb)
}

// Last to call
func Go() {
	for _, cb := range onReadys {
		go cb()
	}
	if os.Getenv("MG_CONF_DUMP") == "1" {
		fmt.Println(Dump(false, "\n"))
	}
}
