package klogin

import (
	"github.com/gin-gonic/gin"
	"github.com/kamasamikon/miego/xmap"
)

func LoginData(c *gin.Context, keys ...string) []string {
	var arr []string

	if c.ContentType() == "application/json" {
		mp := xmap.MapBody(c)
		for _, key := range keys {
			arr = append(arr, mp.S(key))
		}
	} else {
		for _, key := range keys {
			arr = append(arr, c.PostForm(key))
		}
	}

	return arr
}
