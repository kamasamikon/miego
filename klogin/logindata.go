package klogin

import (
	"miego/xmap"

	"github.com/gin-gonic/gin"
)

// 从请求中获取Login需要的字段
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
