package xgin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/kamasamikon/miego/atox"
	"github.com/kamasamikon/miego/in"
)

type PostMap map[string]interface{}

func Map(c *gin.Context) PostMap {
	var m PostMap
	if dat, err := ioutil.ReadAll(c.Request.Body); err != nil {
		return m
	} else {
		json.Unmarshal(dat, &m)
		return m
	}
}

func (pm PostMap) Get(name string) (string, error) {
	if x, ok := pm[name]; ok {
		return x.(string), nil
	}
	return "", fmt.Errorf("'%s' not found", name)
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
