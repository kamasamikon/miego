package xgin

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/kamasamikon/miego/conf"
	"github.com/kamasamikon/miego/page"
)

// Default :Only and default Engine
var Default *gin.Engine

// XXX: Copied from gin/examples/graceful-shutdown/...
func gracefulRun(Engine *gin.Engine, addr string) {
	srv := &http.Server{
		Addr:    addr,
		Handler: Engine,
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

// Go :Default listening on localhost:8888
func Go(Engine *gin.Engine, addr string) {
	if addr == "" {
		port := conf.Int(8888, "i:/ms/port")
		addr = fmt.Sprintf(":%d", port)
	}
	if Engine == nil {
		Engine = Default
	}
	gracefulRun(Engine, addr)
}

//
// Init
//

func init() {
	// gin.SetMode(gin.ReleaseMode)
	// gin.SetMode(gin.ReleaseMode)
	// Default = gin.New()
	Default = gin.Default()
	Default.SetFuncMap(template.FuncMap{
		"FormatAsDate":  FormatAsDate,
		"ToHTML":        ToHTML,
		"ToJS":          ToJS,
		"ToCSS":         ToCSS,
		"Choice":        Choice,
		"ToAttr":        ToAttr,
		"NtimeToString": NtimeToString,
		"MapGet":        MapGet,
		"MapChoice":     MapChoice,
		"SubStr":        SubStr,
	})
}

func DebugSettings(Engine *gin.Engine, prefix string, xRouters int64, xReadme int64, xConf int64) {
	if Engine == nil {
		Engine = Default
	}
	if xRouters == -1 {
		xRouters = conf.Int(1, "i:/gin/debug/routers")
	}
	if xRouters == 1 {
		Engine.GET(prefix+"/debug/routers", func(c *gin.Context) {
			if c.Query("html") == "1" {
				var lines []string
				lines = append(lines, "| Method | Path |")
				lines = append(lines, "| ---- | ---- |")
				for _, x := range Engine.Routes() {
					lines = append(lines, fmt.Sprintf("| %s | %s |", x.Method, x.Path))
				}
				html := page.Markdown("", "", strings.Join(lines, "\\n"))
				c.Data(200, "text/html", []byte(html))
			} else {
				var routers []gin.H
				for _, x := range Engine.Routes() {
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

	if xConf == -1 {
		xConf = conf.Int(1, "i:/gin/debug/conf")
	}
	if xConf == 1 {
		Engine.GET(prefix+"/debug/conf", func(c *gin.Context) {
			c.String(200, conf.DumpRaw(true))
		})
	}
}
