package xgin

import (
	"github.com/gin-gonic/gin"
)

func New(debug bool) *gin.Engine {
	g := gin.New()

	if debug {
		g.GET("/debug/urls", func(c *gin.Context) {
			c.JSON(200, g.Routes())
		})
		g.GET("/debug/doc", func(c *gin.Context) {
			c.String(200, "TODO: Load README.md and show.")
		})
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	return g
}
