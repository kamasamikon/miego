package xgin

import (
	"encoding/json"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"

	"github.com/kamasamikon/miego/conf"
	_ "github.com/kamasamikon/miego/xconf"
)

type KRouter struct {
	Method  string
	Path    string
	Handler string
}

var htmlHead []byte = []byte(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>README</title>
</head>
<body>
`)

var htmlFoot []byte = []byte(`
</body>
</html>
`)

var Engine *gin.Engine

func init() {
	if conf.Int("gin/releaseMode", 0) == 1 {
		Engine = gin.New()
		gin.SetMode(gin.ReleaseMode)
	} else {
		Engine = gin.Default()
		gin.SetMode(gin.DebugMode)
	}

	if conf.Int("gin/debug/routers", 1) == 1 {
		Engine.GET("/debug/routers", func(c *gin.Context) {
			var routers []KRouter
			for _, x := range Engine.Routes() {
				routers = append(routers, KRouter{
					Method:  x.Method,
					Path:    x.Path,
					Handler: x.Handler,
				})
			}
			if data, err := json.MarshalIndent(routers, "", "  "); err == nil {
				c.Data(200, binding.MIMEHTML, htmlHead)
				c.Data(200, binding.MIMEHTML, []byte("<pre>"))
				c.Data(200, binding.MIMEHTML, data)
				c.Data(200, binding.MIMEHTML, []byte("</pre>"))
				c.Data(200, binding.MIMEHTML, htmlFoot)
			} else {
				c.JSON(200, Engine.Routes())
			}
		})
	}

	if conf.Int("gin/debug/readme", 1) == 1 {
		Engine.GET("/debug/readme", func(c *gin.Context) {
			htmlFlags := html.CommonFlags | html.HrefTargetBlank
			opts := html.RendererOptions{Flags: htmlFlags}
			renderer := html.NewRenderer(opts)

			md, err := ioutil.ReadFile("README.md")
			if err != nil {
				c.String(200, err.Error())
				return
			}

			c.Data(200, binding.MIMEHTML, htmlHead)

			body := markdown.ToHTML(md, nil, renderer)
			c.Data(200, binding.MIMEHTML, body)

			c.Data(200, binding.MIMEHTML, htmlFoot)
		})
	}
}
