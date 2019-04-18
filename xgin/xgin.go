package xgin

import (
	"github.com/gin-gonic/gin"
)

func New(debug bool) *gin.Engine {
	g := gin.New()

	if debug {
		g.GET("/gin/routers", func(c *gin.Context) {
			c.JSON(200, g.Routes())
		})
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	return g
}
