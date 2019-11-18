package xgin

import (
	"encoding/json"
	"io/ioutil"

	"github.com/gin-gonic/gin"
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

func (pm PostMap) Get(name string) string {
	if x, ok := pm[name]; ok {
		return x.(string)
	}
	return ""
}
