package xgin

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/kamasamikon/miego/conf"
)

// Default :Only and default Engine
var _Default *gin.Engine

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
	fmt.Println("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		fmt.Println("Server Shutdown Error: ", err)
	}
}

func Go(Engine *gin.Engine, addr string) {
	if conf.Bool(true, "b:/gin/releaseMode") == true {
		gin.SetMode(gin.ReleaseMode)
	}

	if Engine == nil {
		Engine = Default()
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

func Default() *gin.Engine {
	if _Default == nil {
		if conf.Bool(true, "b:/gin/releaseMode") == true {
			gin.SetMode(gin.ReleaseMode)
		}

		_Default = gin.New()

		if conf.Bool(true, "b:/gin/Logger/enable") == true {
			_Default.Use(gin.Logger())
		}
		if conf.Bool(true, "b:/gin/Recovery/enable") == true {
			_Default.Use(gin.RecoveryWithWriter(nil, HandleRecovery))
		}
	}
	return _Default
}
