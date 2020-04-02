package xgin

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kamasamikon/miego/atox"
	"github.com/kamasamikon/miego/in"
)

type PostMap map[string]interface{}

// MakeMap : MakeMap("a", "b", "c", 222) => {"a":"b", "c":222}
func MakeMap(v ...interface{}) PostMap {
	m := make(PostMap)

	for i := 0; i < len(v)/2; i++ {
		m[v[2*i].(string)] = v[2*i+1]
	}

	return m
}

// Map : Convert gin's request to Map
func Map(c *gin.Context) PostMap {
	m := make(PostMap)
	if dat, err := ioutil.ReadAll(c.Request.Body); err != nil {
		return m
	} else {
		json.Unmarshal(dat, &m)
		return m
	}
}

// Map : Convert gin's request to Map
func MapFromQuery(c *gin.Context, useLast bool) PostMap {
	m := make(PostMap)
	for k, a := range c.Request.URL.Query() {
		if useLast {
			v := a[len(a)-1]
			m[k] = v
		} else {
			m[k] = a[0]
		}
	}
	return m
}

func (pm PostMap) Merge(other PostMap, safe bool) {
	if safe {
		for k, v := range other {
			if _, ok := pm[k]; !ok {
				pm[k] = v
			}
		}
	} else {
		for k, v := range other {
			pm[k] = v
		}
	}
}

func (pm PostMap) Has(name string) bool {
	if _, ok := pm[name]; ok {
		return true
	}
	return false
}

func (pm PostMap) Get(name string) (string, bool) {
	if x, ok := pm[name]; ok {
		return x.(string), true
	}
	return "", false
}

func (pm PostMap) Str(name string, defv string) string {
	if x, ok := pm[name]; ok {
		return x.(string)
	}
	return defv
}

func (pm PostMap) Int(name string, defv int) int {
	if x, ok := pm[name]; ok {
		return atox.Int(x.(string), defv)
	}
	return defv
}

func (pm PostMap) Uint(name string, defv uint) uint {
	if x, ok := pm[name]; ok {
		return atox.Uint(x.(string), defv)
	}
	return defv
}

func (pm PostMap) Int64(name string, defv int64) int64 {
	if x, ok := pm[name]; ok {
		return atox.Int64(x.(string), defv)
	}
	return defv
}

func (pm PostMap) Uint64(name string, defv uint64) uint64 {
	if x, ok := pm[name]; ok {
		return atox.Uint64(x.(string), defv)
	}
	return defv
}

func (pm PostMap) Bool(name string, defv bool) bool {
	if x, ok := pm[name]; ok {
		c := x.(string)[0]
		if in.C(c, 'T', 't', 'Y', 'y', '1') {
			return true
		}
		if in.C(c, 'F', 'f', 'N', 'n', '0') {
			return false
		}
	}
	return defv
}

// List : pm["aa"]="a;b;c" ==> ["a", "b", "c"]
func (pm PostMap) List(name string, sep string) []string {
	if x, ok := pm[name]; ok {
		key := x.(string)
		return strings.Split(key, sep)
	}
	return nil
}

// List : pm["aa"]="a;c;c" ==> ["a", "c"]
func (pm PostMap) Set(name string, sep string) map[string]int {
	set := make(map[string]int)
	if x, ok := pm[name]; ok {
		key := x.(string)
		for _, v := range strings.Split(key, sep) {
			set[v] = 1
		}
	}
	return set
}
