package xgin

import (
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
)

func New(debug bool) *gin.Engine {
	g := gin.New()

	if debug {
		gin.SetMode(gin.DebugMode)

		g.GET("/debug/routers", func(c *gin.Context) {
			c.JSON(200, g.Routes())
		})
		g.GET("/debug/readme", func(c *gin.Context) {
			htmlFlags := html.CommonFlags | html.HrefTargetBlank
			opts := html.RendererOptions{Flags: htmlFlags}
			renderer := html.NewRenderer(opts)

			md, err := ioutil.ReadFile("README.md")
			if err != nil {
				c.String(200, err.Error())
			}

			head := []byte(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>README</title>
</head>

<body>
`)

			foot := []byte(`
</body>
</html>
`)
			html := markdown.ToHTML(md, nil, renderer)

			page := head
			page = append(page, html...)
			page = append(page, foot...)

			c.Data(200, binding.MIMEHTML, page)
		})
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	return g
}
