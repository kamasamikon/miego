package conf

import (
	_ "embed"

	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/kamasamikon/miego/klog"
)

//go:embed main.cfg
var main_cfg string

func init() {
	LoadString(main_cfg, false)
}

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

var mapPathEntry = make(map[string]*confEntry)

// Delete an entry
func EntryRem(path string) {
	delete(mapPathEntry, path)
}

// Load file from configure
func EntryAdd(line string, overwrite bool) {
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
	} else if !overwrite {
		return
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

// LoadString : Load setting from string (lines of configuration)
func LoadString(s string, overwrite bool) {
	s = strings.Replace(s, "\r", "\n", -1)
	Lines := strings.Split(s, "\n")
	for _, Line := range Lines {
		EntryAdd(Line, overwrite)
	}
}

// Load : configure from a file.
func Load(fileName string, overwrite bool) error {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		klog.E("Load %s, Error is %s", fileName, err.Error())
		return err
	}
	klog.D("Load %s, Done.", fileName)
	LoadString(string(data), overwrite)
	return nil
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

// List : get a List entry. s:/names=:aaa:bbb first char is the seperator
func List(paths ...string) []string {
	// path: aaa/bbb
	var slice []string

	for _, path := range paths {
		key := "s:/" + path
		if v, ok := mapPathEntry[key]; ok {
			if len(v.vStr) > 0 {
				v.refGet++
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
func SetI(path string, value interface{}, force bool) {
	Set("I:/"+path, value, force)
}
func Seti(path string, value interface{}, force bool) {
	Set("i:/"+path, value, force)
}
func SetS(path string, value interface{}, force bool) {
	Set("S:/"+path, value, force)
}
func Sets(path string, value interface{}, force bool) {
	Set("s:/"+path, value, force)
}
func SetB(path string, value interface{}, force bool) {
	Set("B:/"+path, value, force)
}
func Setb(path string, value interface{}, force bool) {
	Set("b:/"+path, value, force)
}
func SetO(path string, value interface{}, force bool) {
	Set("O:/"+path, value, force)
}
func Seto(path string, value interface{}, force bool) {
	Set("o:/"+path, value, force)
}

func init() {
	{
		cfgList := os.Getenv("KCFG_FILES")
		files := strings.Split(cfgList, ":")
		for _, f := range files {
			if f != "" {
				if err := Load(f, true); err != nil {
					klog.E("LOAD KCFG_FILES Error: %s", err.Error())
				}
			}
		}
	}
	{
		cfgList := os.Getenv("KCFG_QQQ_FILES")
		files := strings.Split(cfgList, ":")
		for _, f := range files {
			if f != "" {
				if err := Load(f, true); err != nil {
					klog.E("LOAD KCFG_QQQ_FILES Error: %s", err.Error())
				}
				os.Remove(f)
			}
		}
	}

	{
		for _, argv := range os.Args {
			if strings.HasPrefix(argv, "--kfg=") {
				f := argv[6:]
				if f != "" {
					if err := Load(f, true); err != nil {
						klog.E("LOAD --kfg=xxx Error: %s", err.Error())
					}
				}
			}
		}
	}
	{
		for _, argv := range os.Args {
			if strings.HasPrefix(argv, "--kfg-qqq=") {
				f := argv[6:]
				if f != "" {
					if err := Load(f, true); err != nil {
						klog.E("LOAD --kfg-qqq=xxx Error: %s", err.Error())
					}
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
