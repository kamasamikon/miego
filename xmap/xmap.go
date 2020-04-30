package xmap

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kamasamikon/miego/atox"
	"github.com/kamasamikon/miego/in"
)

type Map map[string]interface{}

// Make : Make("a", "b", "c", 222) => {"a":"b", "c":222}
func Make(v ...interface{}) Map {
	xm := make(Map)

	for i := 0; i < len(v)/2; i++ {
		xm[v[2*i].(string)] = v[2*i+1]
	}

	return xm
}

// Map : Convert gin's request to Map
func MapBody(c *gin.Context) Map {
	xm := make(Map)
	if dat, err := ioutil.ReadAll(c.Request.Body); err != nil {
		return xm
	} else {
		json.Unmarshal(dat, &xm)
		return xm
	}
}

// Map : Convert gin's request to Map
func MapQuery(c *gin.Context, useLast bool) Map {
	xm := make(Map)
	for k, a := range c.Request.URL.Query() {
		if useLast {
			v := a[len(a)-1]
			xm[k] = v
		} else {
			xm[k] = a[0]
		}
	}
	return xm
}

func MapAll(c *gin.Context, overwrite bool) Map {
	xm := make(Map)
	if dat, err := ioutil.ReadAll(c.Request.Body); err == nil {
		json.Unmarshal(dat, &xm)
	}
	for k, a := range c.Request.URL.Query() {
		if overwrite {
			v := a[len(a)-1]
			xm[k] = v
		} else {
			xm[k] = a[0]
		}
	}
	return xm
}

func (xm Map) String() string {
	if data, err := json.Marshal(xm); err == nil {
		return string(data)
	} else {
		return ""
	}
}

func (xm Map) Merge(other Map, safe bool) {
	if safe {
		for k, v := range other {
			if _, ok := xm[k]; !ok {
				xm[k] = v
			}
		}
	} else {
		for k, v := range other {
			xm[k] = v
		}
	}
}

func (xm Map) Append(safe bool, args ...interface{}) Map {
	for i := 0; i < len(args)/2; i++ {
		k := args[2*i].(string)
		v := args[2*i+1]

		if safe {
			if _, ok := xm[k]; !ok {
				xm[k] = v
			}
		} else {
			xm[k] = v
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
		return x.(string), true
	}
	return "", false
}

func (xm Map) Str(name string, defv string) string {
	if x, ok := xm[name]; ok {
		return x.(string)
	}
	return defv
}
func (xm Map) S(name string) string {
	return xm.Str(name, "")
}

func (xm Map) Int(name string, defv int) int {
	if x, ok := xm[name]; ok {
		return atox.Int(x.(string), defv)
	}
	return defv
}
func (xm Map) I(name string) int {
	return xm.Int(name, 0)
}

func (xm Map) Uint(name string, defv uint) uint {
	if x, ok := xm[name]; ok {
		return atox.Uint(x.(string), defv)
	}
	return defv
}
func (xm Map) U(name string) uint {
	return xm.Uint(name, 0)
}

func (xm Map) Int64(name string, defv int64) int64 {
	if x, ok := xm[name]; ok {
		return atox.Int64(x.(string), defv)
	}
	return defv
}

func (xm Map) Uint64(name string, defv uint64) uint64 {
	if x, ok := xm[name]; ok {
		return atox.Uint64(x.(string), defv)
	}
	return defv
}

func (xm Map) Bool(name string, defv bool) bool {
	if x, ok := xm[name]; ok {
		if b, ok := x.(bool); ok {
			return b
		} else if s, ok := x.(string); ok {
			c := s[0]
			if in.C(c, 'T', 't', 'Y', 'y', '1') {
				return true
			}
			if in.C(c, 'F', 'f', 'N', 'n', '0') {
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
		key := x.(string)
		return strings.Split(key, sep)
	}
	return nil
}

// List : xm["aa"]="a;c;c" ==> ["a", "c"]
func (xm Map) Set(name string, sep string) map[string]int {
	set := make(map[string]int)
	if x, ok := xm[name]; ok {
		key := x.(string)
		for _, v := range strings.Split(key, sep) {
			set[v] = 1
		}
	}
	return set
}
