package xgin

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/kamasamikon/miego/conf"
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
	if Engine == nil {
		Engine = Default
	}

	Engine.SetFuncMap(template.FuncMap{
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

	if conf.Int(1, "i:/gin/releaseMode") == 1 {
		gin.SetMode(gin.ReleaseMode)
	}

	for i, x := range Engine.Routes() {
		conf.Set(
			fmt.Sprintf("s:/gin/routers/%02d", i),
			fmt.Sprintf("%s -> '%s'", x.Method, x.Path),
			true,
		)
	}

	if addr == "" {
		addr = fmt.Sprintf(":%d", conf.Int(8888, "i:/ms/port"))
	}
	gracefulRun(Engine, addr)
}

//
// Init
//

func init() {
	// Default = gin.New()
	Default = gin.Default()
}
