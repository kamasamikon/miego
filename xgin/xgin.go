package xgin

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"miego/conf"
)

// Default :Only and default Engine
var _Default *gin.Engine

func IsPortAvailable(port int) bool {
	addr := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return false
	}
	listener.Close()
	return true
}

func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

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
	quit := make(chan os.Signal, 1)
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
		strfmt = "gin/routers/%01d"
	} else if cnt < 100 {
		strfmt = "gin/routers/%02d"
	} else if cnt < 1000 {
		strfmt = "gin/routers/%03d"
	} else {
		strfmt = "gin/routers/%04d"
	}

	for i, x := range Routes {
		conf.SSetf(
			fmt.Sprintf(strfmt, i),
			fmt.Sprintf("%s -> '%s'", x.Method, x.Path),
		)
	}
}

// ParseAddress 解析 "IP:Port?" 格式的字符串
// 支持：IP可选、端口可选、问号可选
func ParseAddress(s string) (ip string, port int, q bool, e error) {
	if strings.HasSuffix(s, "?") {
		q = true
		s = strings.TrimSuffix(s, "?")
	}

	segs := strings.Split(s, ":")
	if len(segs) != 2 {
		e = fmt.Errorf("Bad address: '%s'", s)
		return
	}

	ipPart := segs[0]
	portPart := segs[1]

	// 解析IP部分（可能为空）
	if ipPart != "" {
		if tmp := net.ParseIP(ipPart); tmp != nil {
			ip = ipPart
		} else {
			e = fmt.Errorf("无效的IP地址: %s", ipPart)
			return
		}
	}

	// 解析端口部分（可能为空）
	if portPart != "" {
		xport, err := strconv.Atoi(portPart)
		if err != nil {
			e = fmt.Errorf("无效的端口号: %s", portPart)
			return
		}
		if xport < 0 || xport > 65535 {
			e = fmt.Errorf("端口号超出范围: %d (应为 0-65535)", xport)
			return
		}
		port = xport
		return
	}

	return
}

func Go(
	Engine *gin.Engine,
	addr string,
	blockMode bool,
	cb func(Engine *gin.Engine),
) error {
	if conf.BTrue("gin/releaseMode") {
		gin.SetMode(gin.ReleaseMode)
	}

	// addr == ":" => default
	// addr == ":?" => random
	// addr == ":NNN?" => NNN or random
	if ip, port, q, err := ParseAddress(addr); err == nil {
		if q {
			if port == 0 {
				// ...:? => GetFreePort
				port, _ = GetFreePort()
			} else {
				// ...:NNN? => NNN ? GetFreePort
				if !IsPortAvailable(port) {
					port, _ = GetFreePort()
				}
			}
		} else {
			if port == 0 {
				// ...: => Load from conf
				port = int(conf.I("ms/port", 8888))
			}
		}
		addr = fmt.Sprintf("%s:%d", ip, port)
	} else {
		return err
	}

	conf.SSetf("gin/addr", addr)
	if a, err := net.ResolveTCPAddr("tcp", addr); err == nil {
		conf.SSetf("gin/addr/ip", a.IP.String())
		conf.ISetf("gin/addr/port", a.Port)
	}
	RoutersToConf(Engine)

	if cb != nil {
		cb(Engine)
	}

	if blockMode {
		gracefulRun(Engine, addr)
	} else {
		go gracefulRun(Engine, addr)
	}
	return nil
}

func Default() *gin.Engine {
	if _Default == nil {
		if conf.BTrue("gin/releaseMode") {
			gin.SetMode(gin.ReleaseMode)
		}

		_Default = gin.New()

		if conf.BTrue("gin/cors/enable") {
			_Default.Use(cors.Default())
		}
		if conf.BTrue("gin/Logger/enable") {
			_Default.Use(gin.Logger())
		}
		if conf.BTrue("gin/Recovery/enable") {
			_Default.Use(gin.RecoveryWithWriter(nil, HandleRecovery))
		}
	}
	return _Default
}
