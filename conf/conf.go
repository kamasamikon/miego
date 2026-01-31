package conf

import (
	"embed"
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
	_bPathReady = "conf/ready"
	_bDebug     = "conf/debug"
)

//go:embed assets/*
var Assets embed.FS

// ///////////////////////////////////////////////////////////////////////
// TYPES
// ///////////////////////////////////////////////////////////////////////

// Monitor is a Callback called when wathed entry modified.
type Monitor func(key string, vnow any, vnew any)

type confEntry struct {
	kind byte
	path string

	vInt  int64
	vStr  string
	vBool bool
	vObj  any

	hidden bool

	monitors []Monitor
}

type ConfCenter struct {
	Name string

	mutex sync.Mutex

	mapPathEntry map[string]*confEntry

	iItems map[string]*iItem
	sItems map[string]*sItem
	bItems map[string]*bItem
	eItems map[string]*eItem

	loadOKCount int
	loadNGCount int
	debug       int

	// OnReady : Called when all configure loaded.
	onReadys []func()
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

		loadOKCount: 0,
		loadNGCount: 0,
		mutex:       sync.Mutex{},
		debug:       0,
		onReadys:    nil,
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
var Default = New("default")

var ccList = make(map[string]*ConfCenter)

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

// XXX: path = kind + key
func pathParse(path string) (kind byte, key string) {
	switch path[0] {
	case 'i', 'I':
		kind = 'i'
		key = path[3:]

	case 's', 'S':
		kind = 's'
		key = path[3:]

	case 'b', 'B':
		kind = 'b'
		key = path[3:]

	case 'o', 'O':
		kind = 'o'
		key = path[3:]

	case 'e', 'E':
		kind = 'e'
		key = path[3:]

	default:
		dp("BadType: %s", path)
		key = ""
	}

	return kind, key
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
			hidden:   e.hidden,
			monitors: monitors,
		}
	}
	cc.mutex.Unlock()

	// 覆盖一些特别的配置
	n.EntryAdd("s:/conf/name", n.Name, true)

	return n
}

// Get value as string
func (cc *ConfCenter) Raw(path string) (string, bool) {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	kind, key := pathParse(path)

	switch kind {
	case 'i':
		if item, ok := cc.iItems[key]; ok {
			return fmt.Sprintf("%v", item.value), true
		}

	case 's':
		if item, ok := cc.sItems[key]; ok {
			return fmt.Sprintf("%v", item.value), true
		}

	case 'b':
		if item, ok := cc.bItems[key]; ok {
			return fmt.Sprintf("%v", item.value), true
		}

	case 'e':
		if _, ok := cc.eItems[key]; ok {
			return fmt.Sprintf("%v", "..."), true
		}
	}

	return "", false
}

// Names : All Keys
func (cc *ConfCenter) Names() []string {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()

	var names []string

	for k := range cc.iItems {
		names = append(names, "i:/"+k)
	}
	for k := range cc.sItems {
		names = append(names, "s:/"+k)
	}
	for k := range cc.bItems {
		names = append(names, "b:/"+k)
	}
	for k := range cc.eItems {
		names = append(names, "e:/"+k)
	}

	return names
}

func (cc *ConfCenter) OnReady(cb func()) {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()
	cc.onReadys = append(cc.onReadys, cb)
}

func (cc *ConfCenter) Ready() {
	cc.BSet(_bPathReady, true)
}

// Last to call
func (cc *ConfCenter) Go() {
	for _, cb := range cc.onReadys {
		go cb()
	}
	if os.Getenv("MG_CONF_DUMP") == "1" {
		fmt.Println(cc.Dump("\n"))
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
	Default.BSetf(_bPathReady, false)

	if os.Getenv("MG_CONF_DEBUG") == "debug" {
		Default.BSetf(_bDebug, true)
		Default.debug = 1
	} else {
		Default.BSetf(_bDebug, false)
		Default.debug = 0
	}

	// 优先级: 命令行 > 环境变量
	Default.LoadFromEnv()
	Default.LoadFromArg()

	// Load environment to conf
	for _, env := range os.Environ() {
		segs := strings.SplitN(env, "=", 2)
		if len(segs) == 2 {
			Default.EntryAdd("s:/env/"+segs[0], segs[1], true)
		}
	}
}
