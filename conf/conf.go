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
	_PathReady     = "e:/conf/ready"
	_Debug         = "i:/conf/debug"
	_MissedEntries = "s:/conf/missedEntries"
	_CCList        = "s:/conf/ccList"
)

//go:embed assets/*
var Assets embed.FS

// ///////////////////////////////////////////////////////////////////////
// TYPES
// ///////////////////////////////////////////////////////////////////////

type setter func(key string, val any) (vret any, ok bool) // e.val = vret if ok
type getter func(key string) (vnew any, ok bool)          // return vnew if ok else e.val

// Monitor is a Callback called when wathed entry modified.
type Monitor func(key string, vnow any, vnew any)

type confEntry struct {
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

	monitors []Monitor
}

type ConfCenter struct {
	Name string

	mutex sync.Mutex

	mapPathEntry     map[string]*confEntry
	mapMissedEntries map[string]int

	iItems map[string]*iItem
	sItems map[string]*sItem
	bItems map[string]*bItem
	eItems map[string]*eItem

	loadOKCount int
	loadNGCount int
	debug       int

	// OnReady : Called when all configure loaded.
	onReadys []func()

	// 防止扑空
	// path <> 1
	mapPendingSetter map[string]setter
	mapPendingGetter map[string]getter
}

// ///////////////////////////////////////////////////////////////////////
// Creator
// ///////////////////////////////////////////////////////////////////////
func New(Name string) *ConfCenter {
	cc := &ConfCenter{
		Name:         Name,
		mapPathEntry: make(map[string]*confEntry),

		iItems: make(map[string]*iItem),
		sItems: make(map[string]*sItem),
		bItems: make(map[string]*bItem),
		eItems: make(map[string]*eItem),

		loadOKCount:      0,
		loadNGCount:      0,
		mapMissedEntries: make(map[string]int),
		mutex:            sync.Mutex{},
		debug:            0,
		onReadys:         nil,
		mapPendingSetter: make(map[string]setter),
		mapPendingGetter: make(map[string]getter),
	}
	cc.EntryAdd("s:/conf/name", Name, true)

	tmpName := Name
	for i := 0; ; i++ {
		if _, ok := ccList[tmpName]; !ok {
			break
		}
		tmpName = fmt.Sprintf("%s-%d", Name, i+1)
	}

	ccList[tmpName] = cc
	cc.Name = tmpName

	return cc
}

// ///////////////////////////////////////////////////////////////////////
// Global
// ///////////////////////////////////////////////////////////////////////
// 默认值，这个必须存在
var Default *ConfCenter = New("default")

var ccList map[string]*ConfCenter = make(map[string]*ConfCenter)

func CCList() []string {
	var names []string
	for name := range ccList {
		names = append(names, name)
	}
	return names
}

func CCByName(name string) *ConfCenter {
	return ccList[name]
}

func SetDefault(name string) {
	cc := ccList[name]
	if cc != nil {
		Default = cc
		return
	}
	SetDefault("default")
}

func GetDefault() *ConfCenter {
	return Default
}

// ///////////////////////////////////////////////////////////////////////
// Helper
// ///////////////////////////////////////////////////////////////////////
// DebugPrint: 打印调试信息
func dp(formating string, args ...any) {
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

// ///////////////////////////////////////////////////////////////////////
// public members
// ///////////////////////////////////////////////////////////////////////

func (cc *ConfCenter) Clone(Name string) *ConfCenter {
	n := New(Name)

	cc.mutex.Lock()
	for p, e := range cc.mapPathEntry {
		var monitors []Monitor
		for _, m := range e.monitors {
			monitors = append(monitors, m)
		}

		n.mapPathEntry[p] = &confEntry{
			kind:     e.kind,
			path:     e.path,
			vInt:     e.vInt,
			vStr:     e.vStr,
			vBool:    e.vBool,
			vObj:     e.vObj,
			setter:   e.setter,
			getter:   e.getter,
			refGet:   e.refGet,
			refSet:   e.refSet,
			hidden:   e.hidden,
			monitors: monitors,
		}
	}
	cc.mutex.Unlock()

	// 覆盖一些特别的配置
	n.EntryAdd("s:/conf/name", n.Name, true)

	for e := range cc.mapMissedEntries {
		n.mapMissedEntries[e] = 1
	}

	return n
}

// ///////////////////////////////////////////////////////////////////////
// private members
// ///////////////////////////////////////////////////////////////////////
func (cc *ConfCenter) setMissedEntries(path string) {
	cc.mapMissedEntries[path] = 1
}

// called when locked
func (cc *ConfCenter) setByEntry(e *confEntry, value any) {
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
			if vv, ok := e.setter(e.path, vNew); ok {
				vNew = vv.(int64)
			}
		}
		e.vInt = vNew
		e.refSet++
		cc.monitorCall(e, vOld, vNew)

	case 's':
		vOld := e.vStr

		vNew := value.(string)
		if e.setter != nil {
			if vv, ok := e.setter(e.path, vNew); ok {
				vNew = vv.(string)
			}
		}
		e.vStr = vNew
		e.refSet++
		cc.monitorCall(e, vOld, vNew)

	case 'b':
		vOld := e.vBool

		vNew := value.(bool)
		if e.setter != nil {
			if vv, ok := e.setter(e.path, vNew); ok {
				vNew = vv.(bool)
			}
		}
		e.vBool = vNew
		e.refSet++
		cc.monitorCall(e, vOld, vNew)

	case 'o':
		vOld := e.vObj

		vNew := value
		if e.setter != nil {
			if vv, ok := e.setter(e.path, vNew); ok {
				vNew = vv
			}
		}
		e.vObj = vNew
		e.refSet++
		cc.monitorCall(e, vOld, vNew)

	case 'e':
		vNew := value.(string)
		if e.setter != nil {
			if vv, ok := e.setter(e.path, vNew); ok {
				vNew = vv.(string)
			}
		}
		e.refSet++
		cc.monitorCall(e, 0, vNew)
	}
}

// ///////////////////////////////////////////////////////////////////////
// Public members
// ///////////////////////////////////////////////////////////////////////
// Delete an entry
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
		cc.mapPathEntry[e.path] = e
		if getter, ok := cc.mapPendingGetter[e.path]; ok {
			e.getter = getter
			delete(cc.mapPendingGetter, e.path)
		}
		if setter, ok := cc.mapPendingSetter[e.path]; ok {
			e.setter = setter
			delete(cc.mapPendingSetter, e.path)
		}
	}
	cc.setByEntry(e, vNew)
}

// Load : configure from a file.
func (cc *ConfCenter) LoadFromText(text string, overwrite bool) {
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

			cc.EntryAdd(path, sb.String(), overwrite)
		} else {
			cc.EntryAdd(path, value, overwrite)
		}
	}
}

// Load : configure from a file.
func (cc *ConfCenter) LoadFromFile(fileName string, overwrite bool) error {
	const (
		NGName = "s:/conf/Load/NG/%d/Name=%s"
		NGWhy  = "s:/conf/Load/NG/%d/Why=%s"
		OKName = "s:/conf/Load/OK/%d=%s"
	)

	data, err := os.ReadFile(fileName)
	if err != nil {
		cc.EntryAddByLine(fmt.Sprintf(NGName, cc.loadNGCount, fileName), false)
		cc.EntryAddByLine(fmt.Sprintf(NGWhy, cc.loadNGCount, err.Error()), false)
		cc.loadNGCount++
		dp("Error:'%s', fileName:'%s'", err.Error(), fileName)
		return err
	}

	cc.EntryAddByLine(fmt.Sprintf(OKName, cc.loadOKCount, fileName), false)
	cc.loadOKCount++

	cc.LoadFromText(string(data), overwrite)
	return nil
}

// Ref : refGet, refSet
func (cc *ConfCenter) Ref(path string) (int64, int64) {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	if e, ok := cc.mapPathEntry[path]; ok {
		return e.refGet, e.refSet
	}
	return -1, -1
}

// Has : Check if a entry exists
func (cc *ConfCenter) Has(path string) bool {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	e, ok := cc.mapPathEntry[path]
	if ok {
		e.refGet++
	}
	return ok
}

// Get value as string
func (cc *ConfCenter) Raw(name string) (string, bool) {
	cc.mutex.Lock()
	e, ok := cc.mapPathEntry[name]
	if !ok {
		cc.mutex.Unlock()
		return "", false
	}
	cc.mutex.Unlock()

	switch e.kind {
	case 'i':
		vInt := e.vInt
		if e.getter != nil {
			if vv, ok := e.getter(e.path); ok {
				vInt = vv.(int64)
			}
		}
		return fmt.Sprintf("%v", vInt), true

	case 's':
		vStr := e.vStr
		if e.getter != nil {
			if vv, ok := e.getter(e.path); ok {
				vStr = vv.(string)
			}
		}
		return fmt.Sprintf("%v", vStr), true

	case 'b':
		vBool := e.vBool
		if e.getter != nil {
			if vv, ok := e.getter(e.path); ok {
				vBool = vv.(bool)
			}
		}
		return fmt.Sprintf("%v", vBool), true

	case 'o':
		return fmt.Sprintf("%v", "..."), true

	case 'e':
		return fmt.Sprintf("%v", "..."), true
	}

	return "", false
}

// Flip : flip on bool
func (cc *ConfCenter) Flip(path string) bool {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	if e, ok := cc.mapPathEntry[path]; ok {
		e.refGet++
		vBool := e.vBool
		if e.getter != nil {
			if vv, ok := e.getter(e.path); ok {
				vBool = vv.(bool)
			}
		}
		vNew := !vBool
		cc.setByEntry(e, vNew)
		return vNew
	}
	cc.setMissedEntries(path)
	return false
}

// Str : get a str typed configure
func (cc *ConfCenter) Str(defval string, paths ...string) string {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	// path: aaa/bbb
	for _, path := range paths {
		if e, ok := cc.mapPathEntry[path]; ok {
			e.refGet++
			if e.getter != nil {
				if vv, ok := e.getter(e.path); ok {
					return vv.(string)
				}
			}
			return e.vStr
		}
		cc.setMissedEntries(path)
	}
	return defval
}

// Str : get a str typed configure
func (cc *ConfCenter) StrX(paths ...string) (string, bool) {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	// path: aaa/bbb
	for _, path := range paths {
		if e, ok := cc.mapPathEntry[path]; ok {
			e.refGet++
			if e.getter != nil {
				if vv, ok := e.getter(e.path); ok {
					return vv.(string), true
				}
			}
			return e.vStr, true
		}
		cc.setMissedEntries(path)
	}
	return "", false
}

// Bool : get a bool entry
func (cc *ConfCenter) Bool(defval bool, paths ...string) bool {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	// path: aaa/bbb
	for _, path := range paths {
		if e, ok := cc.mapPathEntry[path]; ok {
			e.refGet++
			if e.getter != nil {
				if vv, ok := e.getter(e.path); ok {
					return vv.(bool)
				}
			}
			return e.vBool
		}
		cc.setMissedEntries(path)
	}
	return defval
}

// Bool : get a bool entry
func (cc *ConfCenter) BoolX(paths ...string) (bool, bool) {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	// path: aaa/bbb
	for _, path := range paths {
		if e, ok := cc.mapPathEntry[path]; ok {
			e.refGet++
			if e.getter != nil {
				if vv, ok := e.getter(e.path); ok {
					return vv.(bool), true
				}
			}
			return e.vBool, true
		}
		cc.setMissedEntries(path)
	}
	return false, false
}

// Object : get a bool entry
func (cc *ConfCenter) Obj(defval any, paths ...string) any {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	// path: aaa/bbb
	for _, path := range paths {
		if e, ok := cc.mapPathEntry[path]; ok {
			e.refGet++
			if e.getter != nil {
				if vv, ok := e.getter(e.path); ok {
					return vv
				}
			}
			return e.vObj
		}
		cc.setMissedEntries(path)
	}
	return defval
}

// Object : get a bool entry
func (cc *ConfCenter) ObjX(paths ...string) (any, bool) {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	// path: aaa/bbb
	for _, path := range paths {
		if e, ok := cc.mapPathEntry[path]; ok {
			e.refGet++
			if e.getter != nil {
				if vv, ok := e.getter(e.path); ok {
					return vv, true
				}
			}
			return e.vObj, true
		}
		cc.setMissedEntries(path)
	}
	return nil, false
}

// List : get a List entry. s:/names=:aaa:bbb first char is the seperator
func (cc *ConfCenter) List(sep string, paths ...string) []string {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	// path: aaa/bbb
	var slice []string
	for _, path := range paths {
		if e, ok := cc.mapPathEntry[path]; ok {
			if len(e.vStr) > 0 {
				e.refGet++
				vStr := e.vStr
				if e.getter != nil {
					if vv, ok := e.getter(e.path); ok {
						vStr = vv.(string)
					}
				}
				for _, s := range strings.Split(vStr, sep) {
					if s != "" {
						slice = append(slice, s)
					}
				}
			}
		} else {
			cc.setMissedEntries(path)
		}
	}
	return slice
}

// Names : All Keys
func (cc *ConfCenter) Names() []string {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	var names []string
	for k := range cc.mapPathEntry {
		names = append(names, k)
	}
	return names
}

// SafeNames : Keys not hidden
func (cc *ConfCenter) SafeNames() []string {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	var names []string
	for k, e := range cc.mapPathEntry {
		if !e.hidden {
			names = append(names, k)
		}
	}
	return names
}

// cc.Set : Modify or Add conf entry
func (cc *ConfCenter) Set(path string, value any, create bool) {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	var e *confEntry
	var ok bool

	kind, hidden, realpath := pathParse(path)
	if realpath == "" {
		return
	}

	e, ok = cc.mapPathEntry[realpath]
	if !ok {
		if create {
			e = &confEntry{
				kind:   kind,
				hidden: hidden,
				path:   realpath,
			}
			cc.mapPathEntry[e.path] = e
			if getter, ok := cc.mapPendingGetter[e.path]; ok {
				e.getter = getter
				delete(cc.mapPendingGetter, e.path)
			}
			if setter, ok := cc.mapPendingSetter[e.path]; ok {
				e.setter = setter
				delete(cc.mapPendingSetter, e.path)
			}
		} else {
			return
		}
	}

	cc.setByEntry(e, value)
}

func (cc *ConfCenter) LoadFromEnv() {
	{
		cfgList := os.Getenv("KCFG_FILES")
		files := strings.Split(cfgList, ":")
		for _, f := range files {
			if f != "" {
				cc.LoadFromFile(f, true)
			}
		}
	}

	{
		cfgList := os.Getenv("KCFG_QQQ_FILES")
		files := strings.Split(cfgList, ":")
		for _, f := range files {
			if f != "" {
				cc.LoadFromFile(f, true)
				os.Remove(f)
			}
		}
	}
}

func (cc *ConfCenter) LoadFromArg() {
	argc := len(os.Args)

	// --kfg abc.cfg --kfg=xyz.cfg
	{
		for i, argv := range os.Args {
			if argv == "--kfg" {
				i++
				if i < argc {
					cc.LoadFromFile(os.Args[i], true)
				}
				continue
			}
			if strings.HasPrefix(argv, "--kfg=") {
				f := argv[6:]
				if f != "" {
					cc.LoadFromFile(f, true)
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
					cc.LoadFromFile(os.Args[i], true)
					os.Remove(os.Args[i])
				}
				continue
			}
			if strings.HasPrefix(argv, "--kfg-qqq=") {
				f := argv[6:]
				if f != "" {
					cc.LoadFromFile(f, true)
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
				cc.EntryAddByLine(os.Args[i], true)
			}
			continue
		}
		if strings.HasPrefix(argv, "--kfg-item=") {
			item := argv[11:]
			cc.EntryAddByLine(item, true)
		}
	}
}

func (cc *ConfCenter) OnReady(cb func()) {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()
	cc.onReadys = append(cc.onReadys, cb)
}

func (cc *ConfCenter) Ready() {
	cc.Set(_PathReady, "", true)
}

// Last to call
func (cc *ConfCenter) Go() {
	for _, cb := range cc.onReadys {
		go cb()
	}
	if os.Getenv("MG_CONF_DUMP") == "1" {
		fmt.Println(cc.Dump(false, "\n"))
	}
}

// ///////////////////////////////////////////////////////////////////////
// INIT
// ///////////////////////////////////////////////////////////////////////
func init() {
	data, err := Assets.ReadFile("assets/main.cfg")
	if err == nil {
		Default.LoadFromText(string(data), false)
	}

	//
	// Some builtin entries
	//
	Default.Set(_PathReady, "", true)
	Default.Set(_MissedEntries, "", true)
	Default.Set(_CCList, "", true)

	Default.SetGetter(_MissedEntries, func(_ string) (vv any, ok bool) {
		var sb strings.Builder
		for p := range Default.mapMissedEntries {
			sb.WriteString(p)
			sb.WriteRune(';')
		}
		return sb.String(), true
	})
	Default.SetGetter(_CCList, func(_ string) (vv any, ok bool) {
		return strings.Join(CCList(), ";"), true
	})

	if os.Getenv("MG_CONF_DEBUG") == "debug" {
		Default.Set(_Debug, 1, true)
		Default.debug = 1
	} else {
		Default.Set(_Debug, 0, true)
		Default.debug = 0
	}

	Default.SetSetter(_Debug, func(_ string, v any) (vv any, ok bool) {
		Default.debug = v.(int)
		return nil, false
	})

	// 优先级: 命令行 > 环境变量
	Default.LoadFromEnv()
	Default.LoadFromArg()

	// Load environment to conf
	for _, env := range os.Environ() {
		segs := strings.SplitN(env, "=", 2)
		if len(segs) == 2 {
			Default.Set("s:/env/"+segs[0], segs[1], true)
		}
	}
}
