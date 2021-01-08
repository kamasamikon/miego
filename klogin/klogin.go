package klogin

import (
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/twinj/uuid"

	"github.com/kamasamikon/miego/klog"
	"github.com/kamasamikon/miego/xgin"
	"github.com/kamasamikon/miego/xmap"
)

type KLogin struct {
	Gin *gin.Engine

	//
	// Session
	//
	SessionName string

	//
	// Login
	//

	// Login OK
	LoginRouter []string // default: /login

	//
	// Logout
	//

	// Logout OK
	LogoutRouter []string // default: /logout

	//
	// Call this to verify the login parameters
	//
	LoginDataChecker func(c *gin.Context) (sessionItems xmap.Map, OKRedirectURL string, NGPageName string, NGPageParam xmap.Map, err error)

	BeforeLogin  func(c *gin.Context) (StatusCode int, LoginPageName string, LoginPageParam xmap.Map)
	BeforeLogout func(c *gin.Context) (LogoutRedirectURL string)

	//
	// Before and After check.
	//
	BCheckerList []func(h gin.HandlerFunc) gin.HandlerFunc
	ACheckerList []func(h gin.HandlerFunc) gin.HandlerFunc
}

func (o *KLogin) isLoggin(h gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		UUID := session.Get("UUID")
		if UUID != nil {
			klog.D("IsLoggin: OK: UUID: %s", UUID.(string))
			h(c)
			return
		}

		// Return Status or Login page
		StatusCode, LoginPageName, LoginPageParam := o.BeforeLogin(c)
		if LoginPageName == "" {
			klog.D("IsLoggin: NG: JSON: %s", spew.Sdump(LoginPageParam))
			c.JSON(StatusCode, LoginPageParam)
		} else {
			klog.D("IsLoggin: NG: HTML: %s", spew.Sdump(LoginPageParam))
			c.HTML(StatusCode, LoginPageName, LoginPageParam)
		}
	}
}

func (o *KLogin) Get(c *gin.Context, key string) (string, bool) {
	session := sessions.Default(c)

	if val := session.Get(key); val == nil {
		return "", false
	} else {
		return val.(string), true
	}
}

func (o *KLogin) Set(c *gin.Context, key string, val interface{}) {
	session := sessions.Default(c)
	session.Set(key, val)
}

func (o *KLogin) Save(c *gin.Context) {
	session := sessions.Default(c)
	session.Save()
}

func (o *KLogin) doLogin(c *gin.Context) {
	session := sessions.Default(c)

	sessionItems, OKRedirectURL, NGPageName, NGPageParam, err := o.LoginDataChecker(c)
	if err == nil {
		for k, v := range sessionItems {
			session.Set(k, v)
		}

		session.Set("UUID", uuid.NewV4().String())
		session.Save()

		klog.D(OKRedirectURL)
		c.Redirect(302, OKRedirectURL)
	} else {
		klog.Dump(NGPageName)
		klog.Dump(NGPageParam)
		session.Clear()
		session.Save()

		c.HTML(200, NGPageName, NGPageParam)
	}
}

func (o *KLogin) doLogout(c *gin.Context) {
	session := sessions.Default(c)

	LogoutRedirectURL := o.BeforeLogout(c)
	session.Clear()
	session.Save()
	c.Redirect(302, LogoutRedirectURL)
}

func (o *KLogin) Setup(Gin *gin.Engine) {
	// Gin
	o.Gin = Gin

	// Redis/Session
	redisHost := os.Getenv("DOCKER_GATEWAY")
	redisAddr := redisHost + ":6379"
	store, err := redis.NewStore(10, "tcp", redisAddr, "", []byte("secret"))
	if err != nil {
		klog.E(err.Error())
		return
	}
	Gin.Use(sessions.Sessions(o.SessionName, store))

	// Routers
	if o.LoginRouter == nil {
		Gin.POST("/login", o.doLogin)
	} else {
		for _, r := range o.LoginRouter {
			Gin.POST(r, o.doLogin)
		}
	}
	if o.LogoutRouter == nil {
		Gin.GET("/logout", o.doLogout)
	} else {
		for _, r := range o.LogoutRouter {
			Gin.GET(r, o.doLogout)
		}
	}
}

func (o *KLogin) POST(relativePath string, handler gin.HandlerFunc) {
	var decors []func(h gin.HandlerFunc) gin.HandlerFunc
	if o.BCheckerList != nil {
		decors = append(decors, o.BCheckerList...)
	}
	decors = append(decors, o.isLoggin)
	if o.ACheckerList != nil {
		decors = append(decors, o.ACheckerList...)
	}
	o.Gin.POST(relativePath, xgin.Decorator(handler, decors...))
}
func (o *KLogin) GET(relativePath string, handler gin.HandlerFunc) {
	var decors []func(h gin.HandlerFunc) gin.HandlerFunc
	if o.BCheckerList != nil {
		decors = append(decors, o.BCheckerList...)
	}
	decors = append(decors, o.isLoggin)
	if o.ACheckerList != nil {
		decors = append(decors, o.ACheckerList...)
	}

	o.Gin.GET(relativePath, xgin.Decorator(handler, decors...))
}
func (o *KLogin) PUT(relativePath string, handler gin.HandlerFunc) {
	var decors []func(h gin.HandlerFunc) gin.HandlerFunc
	if o.BCheckerList != nil {
		decors = append(decors, o.BCheckerList...)
	}
	decors = append(decors, o.isLoggin)
	if o.ACheckerList != nil {
		decors = append(decors, o.ACheckerList...)
	}

	o.Gin.PUT(relativePath, xgin.Decorator(handler, decors...))
}
func (o *KLogin) DELETE(relativePath string, handler gin.HandlerFunc) {
	var decors []func(h gin.HandlerFunc) gin.HandlerFunc
	if o.BCheckerList != nil {
		decors = append(decors, o.BCheckerList...)
	}
	decors = append(decors, o.isLoggin)
	if o.ACheckerList != nil {
		decors = append(decors, o.ACheckerList...)
	}

	o.Gin.DELETE(relativePath, xgin.Decorator(handler, decors...))
}
