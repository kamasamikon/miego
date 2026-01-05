package xmap

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"miego/klog"

	"github.com/davecgh/go-spew/spew"
)

type Map map[string]interface{}

func First(mpArr []Map) Map {
	if mpArr == nil {
		return nil
	}
	return mpArr[0]
}

func Nth(mpArr []Map, n int) Map {
	size := len(mpArr)
	if size == 0 {
		return nil
	}
	if n < 0 {
		if (-1 * n) > size {
			return nil
		}
		n += size
		return mpArr[n]
	}
	if n >= size {
		return nil
	}
	return mpArr[n]
}

// Make : Make("a", "b", "c", 222) => {"a":"b", "c":222}
func Make(v ...interface{}) Map {
	xm := make(Map)

	for i := 0; i < len(v)/2; i++ {
		xm[v[2*i].(string)] = v[2*i+1]
	}

	return xm
}

// MapData : Convert json string
func MapData(data []byte) Map {
	xm := make(Map)
	if err := json.Unmarshal(data, &xm); err != nil {
		klog.E(err.Error())
	}
	return xm
}

// MapBody : Convert gin's request to Map
func MapBody(r *http.Request) Map {
	xm := make(Map)
	if dat, err := io.ReadAll(r.Body); err != nil {
		return xm
	} else {
		json.Unmarshal(dat, &xm)
		return xm
	}
}

// MapQuery : Convert gin's request to Map
func MapQuery(r *http.Request, useLast bool) Map {
	xm := make(Map)
	for k, a := range r.URL.Query() {
		if useLast {
			v := a[len(a)-1]
			xm[k] = v
		} else {
			xm[k] = a[0]
		}
	}
	return xm
}

// MapAll : MapBody then MapQuery
func MapAll(r *http.Request, overwrite bool) Map {
	xm := make(Map)
	if dat, err := io.ReadAll(r.Body); err == nil {
		json.Unmarshal(dat, &xm)
	}
	for k, a := range r.URL.Query() {
		if overwrite {
			v := a[len(a)-1]
			xm[k] = v
		} else {
			xm[k] = a[0]
		}
	}
	return xm
}

// Marshal : xmap -> string
func (xm Map) Marshal() string {
	if data, err := json.MarshalIndent(xm, "", "\t"); err == nil {
		return string(data)
	} else {
		return ""
	}
}

func (xm Map) Dump(title string, wlist string, blist string) {
	setGen := func(s string) map[string]int {
		set := make(map[string]int)
		for _, v := range strings.Split(s, ":") {
			if v != "" {
				set[v] = 1
			}
		}
		return set
	}
	setHas := func(s string, arr map[string]int) bool {
		_, ok := arr[s]
		return ok
	}

	// While and Black
	wSet := setGen(wlist)
	bSet := setGen(blist)

	// Keys in use
	var keys []string

	for k := range xm {
		if len(bSet) != 0 {
			if setHas(k, bSet) {
				continue
			}
		}

		if len(wSet) != 0 {
			if setHas(k, wSet) {
				keys = append(keys, k)
			}
		} else {
			keys = append(keys, k)
		}
	}

	var wKeys []string
	if len(wSet) != 0 {
		for _, s := range strings.Split(wlist, ":") {
			if setHas(s, wSet) {
				wKeys = append(wKeys, s)
			}
		}
	} else {
		sort.Strings(keys)
		wKeys = keys
	}

	//
	// Calc the width of key column
	//
	width := 1
	for _, s := range wKeys {
		if len(s) > width {
			width = len(s)
		}
	}
	if width > 20 {
		width = 20
	}

	fmtLine := fmt.Sprintf(" %%%ds : %%s", width)

	//
	// Collect fields to output string
	//
	var lines []string
	lines = append(lines, title)
	lines = append(lines, "\r\n")
	for _, k := range wKeys {
		if k == "" {
			continue
		}
		if len(wSet) != 0 {
			if _, ok := wSet[k]; !ok {
				continue
			}
		}
		if v, ok := xm[k]; ok {
			sdump := spew.Sdump(v)
			lines = append(lines, fmt.Sprintf(fmtLine, k, sdump))
		}
	}

	//
	// Print out
	//
	klog.KLog(2, 0, klog.ColorType_W, 'X', strings.Join(lines, ""))
}

// Duplicate on to another
// names = [fromName, toName, fromName, toName, ...]
func (xm Map) Dup(names ...string) {
	for i := 0; i < len(names)/2; i++ {
		frName := names[2*i]
		if frValue, ok := xm[frName]; ok {
			toName := names[2*i+1]
			if toName != "" {
				xm[toName] = frValue
			}
		}
	}
}

// names = [fromName, toName, fromName, toName, ...]
func (xm Map) SafeDup(names ...string) {
	for i := 0; i < len(names)/2; i++ {
		frName := names[2*i]
		if frValue, ok := xm[frName]; ok {
			toName := names[2*i+1]
			if _, ok := xm[toName]; !ok {
				xm[toName] = frValue
			}
		}
	}
}

func (xm Map) Rem(names ...string) {
	for _, name := range names {
		delete(xm, name)
	}
}

func (xm Map) Merge(other Map) {
	for k, v := range other {
		xm[k] = v
	}
}

func (xm Map) SafeMerge(other Map) int {
	cnt := 0
	for k, v := range other {
		if _, ok := xm[k]; !ok {
			xm[k] = v
			cnt += 1
		}
	}
	return cnt
}

func (xm Map) MergeSome(other Map, names ...string) int {
	cnt := 0
	for _, name := range names {
		if _, ok := other[name]; ok {
			klog.W("%s", name)
			xm[name] = other[name]
			cnt += 1
		}
	}
	return cnt
}

func (xm Map) Put(args ...interface{}) Map {
	for i := 0; i < len(args)/2; i++ {
		k := args[2*i].(string)
		v := args[2*i+1]
		xm[k] = v
	}
	return xm
}

// PutKeys : Put with default value
func (xm Map) PutKeys(defval string, args ...interface{}) Map {
	for i := 0; i < len(args); i++ {
		k := args[i].(string)
		xm[k] = defval
	}
	return xm
}

// SafePut : Put only when not exists
func (xm Map) SafePut(args ...interface{}) Map {
	for i := 0; i < len(args)/2; i++ {
		k := args[2*i].(string)
		v := args[2*i+1]

		if _, ok := xm[k]; !ok {
			xm[k] = v
		}
	}
	return xm
}

// ReplaceKeys : Rename Keys, ReplaceKeys(oldKeyA, newKeyA, oldKeyB, newKeyB, ...)
func (xm Map) ReplaceKeys(keys ...string) Map {
	for i := 0; i < len(keys)/2; i++ {
		oldKey := keys[2*i]
		if v, ok := xm[oldKey]; ok {
			if newKey := keys[2*i+1]; newKey != "" {
				xm[newKey] = v
			}
			delete(xm, oldKey)
		}
	}
	return xm
}

func (xm Map) Has(name string) bool {
	if _, ok := xm[name]; ok {
		return true
	}
	return false
}

func (xm Map) Get(name string) (string, bool) {
	if x, ok := xm[name]; ok {
		if s, ok := x.(string); ok {
			return s, true
		}
	}
	return "", false
}

func (xm Map) AsStr(name string, defv string) string {
	if x, ok := xm[name]; ok {
		return fmt.Sprintf("%v", x)
	}
	return defv
}
func (xm Map) AsS(name string) string {
	if x, ok := xm[name]; ok {
		return fmt.Sprintf("%v", x)
	}
	return ""
}

func (xm Map) AsInt(name string, defv int64) int64 {
	if x, ok := xm[name]; ok {
		s := fmt.Sprintf("%v", x)
		xx, e := strconv.ParseInt(s, 0, 64)
		if e != nil {
			return defv
		}
		return xx
	}
	return defv
}
func (xm Map) AsUint(name string, defv uint64) uint64 {
	if x, ok := xm[name]; ok {
		s := fmt.Sprintf("%v", x)
		xx, e := strconv.ParseUint(s, 0, 64)
		if e != nil {
			return defv
		}
		return xx
	}
	return defv
}
func (xm Map) AsFloat(name string, defv float64) float64 {
	if x, ok := xm[name]; ok {
		s := fmt.Sprintf("%v", x)
		xx, e := strconv.ParseFloat(s, 64)
		if e != nil {
			return defv
		}
		return xx
	}
	return defv
}

func (xm Map) StrMust(name string, defv string) string {
	if x, ok := xm[name]; ok {
		if s, ok := x.(string); ok {
			if s == "" {
				return defv
			}
			return s
		}
	}
	return defv
}
func (xm Map) Str(name string, defv string) string {
	if x, ok := xm[name]; ok {
		if s, ok := x.(string); ok {
			return s
		}
	}
	return defv
}
func (xm Map) S(name string) string {
	if x, ok := xm[name]; ok {
		if s, ok := x.(string); ok {
			return s
		}
	}
	return ""
}

func (xm Map) Int(name string, defv int) int {
	if x, ok := xm[name]; ok {
		xx, e := strconv.ParseInt(x.(string), 0, 64)
		if e != nil {
			return defv
		}
		return int(xx)
	}
	return defv
}
func (xm Map) I(name string) int {
	return xm.Int(name, 0)
}

func (xm Map) Uint(name string, defv uint) uint {
	if x, ok := xm[name]; ok {
		xx, e := strconv.ParseUint(x.(string), 0, 64)
		if e != nil {
			return defv
		}
		return uint(xx)
	}
	return defv
}
func (xm Map) U(name string) uint {
	return xm.Uint(name, 0)
}

func (xm Map) Int64(name string, defv int64) int64 {
	if x, ok := xm[name]; ok {
		xx, e := strconv.ParseInt(x.(string), 0, 64)
		if e != nil {
			return defv
		}
		return xx
	}
	return defv
}

func (xm Map) Uint64(name string, defv uint64) uint64 {
	if x, ok := xm[name]; ok {
		xx, e := strconv.ParseUint(x.(string), 0, 64)
		if e != nil {
			return defv
		}
		return xx
	}
	return defv
}

func (xm Map) Bool(name string, defv bool) bool {
	if x, ok := xm[name]; ok {
		if b, ok := x.(bool); ok {
			return b
		} else if s, ok := x.(string); ok {
			switch s[0] {
			case 'T', 't', 'Y', 'y', '1':
				return true
			case 'F', 'f', 'N', 'n', '0':
				return false
			}
		} else if i, ok := x.(int); ok {
			return i != 0
		}
	}
	return defv
}

// List : xm["aa"]="a;b;c" ==> ["a", "b", "c"]
func (xm Map) List(name string, sep string) []string {
	if x, ok := xm[name]; ok {
		if key, ok := x.(string); ok {
			return strings.Split(key, sep)
		}
	}
	return nil
}

// List : xm["aa"]="a;c;c" ==> ["a", "c"]
func (xm Map) Set(name string, sep string) map[string]int {
	set := make(map[string]int)
	if x, ok := xm[name]; ok {
		if key, ok := x.(string); ok {
			for _, v := range strings.Split(key, sep) {
				set[v] = 1
			}
		}
	}
	return set
}
