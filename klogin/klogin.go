package klogin

import (
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/kamasamikon/miego/klog"
	"github.com/kamasamikon/miego/xgin"
	"github.com/kamasamikon/miego/xmap"
	"github.com/twinj/uuid"
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

	BeforeLogin  func(c *gin.Context) (LoginPageName string, LoginPageParam xmap.Map)
	BeforeLogout func(c *gin.Context) (LogoutRedirectURL string)
}

func (o *KLogin) isLoggin(h gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		UUID := session.Get("UUID")
		if UUID != nil {
			klog.D("isLoggin: OK: UUID: %s", UUID.(string))
			h(c)
			return
		}

		LoginPageName, LoginPageParam := o.BeforeLogin(c)
		klog.D("isLoggin: NG: %s", spew.Sdump(LoginPageParam))
		c.HTML(200, LoginPageName, LoginPageParam)
	}
}

func (o *KLogin) Get(c *gin.Context, key string) (string, bool) {
	session := sessions.Default(c)

	if val := session.Get(key); val == nil {
		return "", false
	} else {
		klog.D(key)
		klog.Dump(val)
		return val.(string), true
	}
}
func (o *KLogin) Set(c *gin.Context, key string, val interface{}) {
	session := sessions.Default(c)
	session.Set(key, val)
	session.Save()
}

func (o *KLogin) doLogin(c *gin.Context) {
	session := sessions.Default(c)

	sessionItems, OKRedirectURL, NGPageName, NGPageParam, err := o.LoginDataChecker(c)
	klog.D("")
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
	o.Gin.POST(relativePath, xgin.Decorator(handler, o.isLoggin))
}
func (o *KLogin) GET(relativePath string, handler gin.HandlerFunc) {
	o.Gin.GET(relativePath, xgin.Decorator(handler, o.isLoggin))
}
func (o *KLogin) PUT(relativePath string, handler gin.HandlerFunc) {
	o.Gin.PUT(relativePath, xgin.Decorator(handler, o.isLoggin))
}
func (o *KLogin) DELETE(relativePath string, handler gin.HandlerFunc) {
	o.Gin.DELETE(relativePath, xgin.Decorator(handler, o.isLoggin))
}
