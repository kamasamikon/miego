package xgin

import (
	"encoding/json"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
)

type KRouter struct {
	Method  string
	Path    string
	Handler string
}

func New(debug bool) *gin.Engine {
	g := gin.New()

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

	if debug {
		gin.SetMode(gin.DebugMode)

		g.GET("/debug/routers", func(c *gin.Context) {
			var routers []KRouter
			for _, x := range g.Routes() {
				routers = append(routers, KRouter{
					Method:  x.Method,
					Path:    x.Path,
					Handler: x.Handler,
				})
			}
			if data, err := json.MarshalIndent(routers, "", "  "); err == nil {
				c.Data(200, binding.MIMEHTML, head)
				c.Data(200, binding.MIMEHTML, []byte("<pre>"))
				c.Data(200, binding.MIMEHTML, data)
				c.Data(200, binding.MIMEHTML, []byte("</pre>"))
				c.Data(200, binding.MIMEHTML, foot)
			} else {
				c.JSON(200, g.Routes())
			}
		})

		g.GET("/debug/readme", func(c *gin.Context) {
			htmlFlags := html.CommonFlags | html.HrefTargetBlank
			opts := html.RendererOptions{Flags: htmlFlags}
			renderer := html.NewRenderer(opts)

			md, err := ioutil.ReadFile("README.md")
			if err != nil {
				c.String(200, err.Error())
				return
			}

			c.Data(200, binding.MIMEHTML, head)

			body := markdown.ToHTML(md, nil, renderer)
			c.Data(200, binding.MIMEHTML, body)

			c.Data(200, binding.MIMEHTML, foot)
		})
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	return g
}
