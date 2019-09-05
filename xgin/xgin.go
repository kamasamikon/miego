package xgin

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"

	"github.com/kamasamikon/miego/conf"
)

type ioInfo struct {
	Ping interface{}
	Pong interface{}
}

var ioInfoDict = make(map[string]*ioInfo)

var htmlHead = []byte(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>README</title>
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
	// Default.Run(addr)
	gracefulRun(addr)
}

//
// Set Ping Pong debug information
//

// PPSet : Set PingPong debug info
func PPSet(method string, path string, ping interface{}, pong interface{}) {
	key := fmt.Sprintf("%s@%s", method, path)
	ioInfoDict[key] = &ioInfo{
		Ping: ping,
		Pong: pong,
	}
}

// PPSetPOST : PPSet("POST", ...)
func PPSetPOST(path string, ping interface{}, pong interface{}) {
	PPSet("POST", path, ping, pong)
}

// PPSetGET : PPSet("POST", ...)
func PPSetGET(path string, ping interface{}, pong interface{}) {
	PPSet("GET", path, ping, pong)
}

// PPSetUPDATE : PPSet("POST", ...)
func PPSetUPDATE(path string, ping interface{}, pong interface{}) {
	PPSet("UPDATE", path, ping, pong)
}

// PPSetDELETE : PPSet("POST", ...)
func PPSetDELETE(path string, ping interface{}, pong interface{}) {
	PPSet("DELETE", path, ping, pong)
}

// PPGet : Get PingPong debug info
func _PPGet(method string, path string) *ioInfo {
	key := fmt.Sprintf("%s@%s", method, path)
	if info, ok := ioInfoDict[key]; ok {
		return info
	}
	return nil
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

	if conf.Int(1, "gin/debug/routers") == 1 {
		Default.GET("/debug/routers", func(c *gin.Context) {
			var routers []gin.H
			for _, x := range Default.Routes() {
				var ping interface{}
				var pong interface{}
				pp := _PPGet(x.Method, x.Path)
				if pp != nil {
					ping = pp.Ping
					pong = pp.Pong
				}

				r := gin.H{
					"Method": x.Method,
					"Path":   x.Path,
				}
				if ping != nil {
					r["Ping"] = ping
				}
				if pong != nil {
					r["Pong"] = pong
				}

				routers = append(routers, r)
			}

			if c.Query("html") == "1" {
				if data, err := json.MarshalIndent(routers, "", "  "); err == nil {
					c.Data(200, binding.MIMEHTML, htmlHead)
					c.Data(200, binding.MIMEHTML, []byte("<pre>"))
					c.Data(200, binding.MIMEHTML, data)
					c.Data(200, binding.MIMEHTML, []byte("</pre>"))
					c.Data(200, binding.MIMEHTML, htmlFoot)
					return
				}
			}

			c.JSON(200, routers)
		})
	}

	if conf.Int(1, "gin/debug/readme") == 1 {
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

	if conf.Int(1, "gin/debug/conf") == 1 {
		Default.GET("/debug/conf", func(c *gin.Context) {
			c.String(200, conf.Dump())
		})
	}
}
