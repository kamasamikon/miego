package xgin

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"miego/conf"
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

func RoutersToConf(Engine *gin.Engine) {
	Routes := Engine.Routes()

	var strfmt string
	cnt := len(Routes)
	if cnt < 10 {
		strfmt = "s:/gin/routers/%01d"
	} else if cnt < 100 {
		strfmt = "s:/gin/routers/%02d"
	} else if cnt < 1000 {
		strfmt = "s:/gin/routers/%03d"
	} else {
		strfmt = "s:/gin/routers/%04d"
	}

	for i, x := range Routes {
		conf.Set(
			fmt.Sprintf(strfmt, i),
			fmt.Sprintf("%s -> '%s'", x.Method, x.Path),
			true,
		)
	}
}

func Go(Engine *gin.Engine, addr string) {
	if conf.Bool(true, "b:/gin/releaseMode") == true {
		gin.SetMode(gin.ReleaseMode)
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

		if conf.Bool(true, "b:/gin/cors/enable") == true {
			_Default.Use(cors.Default())
		}
		if conf.Bool(true, "b:/gin/Logger/enable") == true {
			_Default.Use(gin.Logger())
		}
		if conf.Bool(true, "b:/gin/Recovery/enable") == true {
			_Default.Use(gin.RecoveryWithWriter(nil, HandleRecovery))
		}
	}
	return _Default
}
