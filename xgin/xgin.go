package xgin

import (
	"context"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"

	"github.com/kamasamikon/miego/conf"
	"github.com/kamasamikon/miego/page"
)

var htmlHead = []byte(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
	<meta content="width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=0" name="viewport" />
    <title>README</title>
    <style>
      html { font-size: 14px; background-color: var(--bg-color); color: var(--text-color); font-family: "Helvetica Neue", Helvetica, Arial, sans-serif; -webkit-font-smoothing: antialiased; }
    </style>
</head>
<body>
`)

var htmlFoot = []byte(`
</body>
</html>
`)

// Default :Only and default Engine
var Default *gin.Engine

// XXX: Copied from gin/examples/graceful-shutdown/...
func gracefulRun(addr string) {
	srv := &http.Server{
		Addr:    addr,
		Handler: Default,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT)
	<-quit
	fmt.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		fmt.Println("Server Shutdown:", err)
	}
	fmt.Println("Server exiting")
}

// Run :Default listening on localhost:8888
func Run(addr string) {
	if addr == "" {
		port := conf.Int(8888, "ms/port")
		addr = fmt.Sprintf(":%d", port)
	}
	gracefulRun(addr)
}

//
// Init
//

func init() {
	if conf.Int(0, "gin/releaseMode") == 1 {
		gin.SetMode(gin.ReleaseMode)
		Default = gin.New()
	} else {
		gin.SetMode(gin.DebugMode)
		Default = gin.Default()
	}

	Default.SetFuncMap(template.FuncMap{
		"FormatAsDate":  FormatAsDate,
		"ToHTML":        ToHTML,
		"ToJS":          ToJS,
		"ToCSS":         ToCSS,
		"Choice":        Choice,
		"ToAttr":        ToAttr,
		"NtimeToString": NtimeToString,
		"MPGet":         MPGet,
		"MapChoice":     MapChoice,
		"SubStr":        SubStr,
	})
}

func DebugSettings(xRouters int64, xReadme int64, xConf int64) {
	if xRouters == -1 {
		xRouters = conf.Int(1, "gin/debug/routers")
	}
	if xRouters == 1 {
		Default.GET("/debug/routers", func(c *gin.Context) {
			if c.Query("html") == "1" {
				var lines []string
				lines = append(lines, "| Method | Path |")
				lines = append(lines, "| ---- | ---- |")
				for _, x := range Default.Routes() {
					lines = append(lines, fmt.Sprintf("| %s | %s |", x.Method, x.Path))
				}
				html := page.Markdown("", "", strings.Join(lines, "\\n"))
				c.Data(200, "text/html", []byte(html))
			} else {
				var routers []gin.H
				for _, x := range Default.Routes() {
					r := gin.H{
						"Method": x.Method,
						"Path":   x.Path,
					}
					routers = append(routers, r)
				}
				c.JSON(200, routers)
			}
		})
	}

	if xReadme == -1 {
		xReadme = conf.Int(1, "gin/debug/readme")
	}
	if xReadme == 1 {
		Default.GET("/debug/readme", func(c *gin.Context) {
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

	if xConf == -1 {
		xConf = conf.Int(1, "gin/debug/routers")
	}
	if xConf == 1 {
		Default.GET("/debug/conf", func(c *gin.Context) {
			c.String(200, conf.Dump())
		})
	}
}
