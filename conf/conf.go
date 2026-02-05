package conf

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
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

type ConfCenter struct {
	Name string

	mutex sync.Mutex

	// i,s,b 是配置数据
	// e是广播消息
	// x是函数方式获取数据，类型是any
	iItems map[string]*iItem
	sItems map[string]*sItem
	bItems map[string]*bItem
	eItems map[string]*eItem
	xItems map[string]*xItem

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
		iItems: make(map[string]*iItem),
		sItems: make(map[string]*sItem),
		bItems: make(map[string]*bItem),
		eItems: make(map[string]*eItem),
		xItems: make(map[string]*xItem),

		loadOKCount: 0,
		loadNGCount: 0,
		mutex:       sync.Mutex{},
		debug:       0,
		onReadys:    nil,
	}

	tmpName := Name
	for i := 0; ; i++ {
		if _, ok := ccList[tmpName]; !ok {
			break
		}
		tmpName = fmt.Sprintf("%s-%d", Name, i+1)
	}

	ccList[tmpName] = cc
	cc.Name = tmpName
	cc.SSetf("conf/name", Name)

	var ccNames []string
	for ccName := range ccList {
		ccNames = append(ccNames, ccName)
	}
	sort.Strings(ccNames)
	cc.SSetf("conf/names", strings.Join(ccNames, ","))

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
	newcc := New(Name)

	cc.mutex.Lock()
	for key, item := range cc.iItems {
		iItemsNew := &iItem{
			value: item.value,
		}
		if item.monitors != nil {
			monitors := make(map[string]iMonitor)
			for n, m := range item.monitors {
				monitors[n] = m
			}
			iItemsNew.monitors = monitors
		}
		newcc.iItems[key] = iItemsNew
	}
	for key, item := range cc.sItems {
		sItemsNew := &sItem{
			value: item.value,
		}
		if item.monitors != nil {
			monitors := make(map[string]sMonitor)
			for n, m := range item.monitors {
				monitors[n] = m
			}
			sItemsNew.monitors = monitors
		}
		newcc.sItems[key] = sItemsNew
	}
	for key, item := range cc.bItems {
		bItemsNew := &bItem{
			value: item.value,
		}
		if item.monitors != nil {
			monitors := make(map[string]bMonitor)
			for n, m := range item.monitors {
				monitors[n] = m
			}
			bItemsNew.monitors = monitors
		}
		newcc.bItems[key] = bItemsNew
	}
	for key, item := range cc.eItems {
		eItemsNew := &eItem{}
		if item.listeners != nil {
			listeners := make(map[string]eListener)
			for n, m := range item.listeners {
				listeners[n] = m
			}
			eItemsNew.listeners = listeners
		}
		newcc.eItems[key] = eItemsNew
	}
	for key, item := range cc.xItems {
		xItemsNew := &xItem{
			setter: item.setter,
			getter: item.getter,
		}
		newcc.xItems[key] = xItemsNew
	}

	newcc.loadOKCount = cc.loadOKCount
	newcc.loadNGCount = cc.loadNGCount
	newcc.debug = cc.debug
	cc.mutex.Unlock()

	return newcc
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
	if cc.BTrue(_bPathReady) {
		return
	}

	cc.Ready()
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
			Default.SSetf("env/"+segs[0], segs[1])
		}
	}
}
