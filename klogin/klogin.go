package klogin

import (
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/kamasamikon/miego/klog"
	"github.com/kamasamikon/miego/xgin"
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

	// Not login, show login page
	LoginPageName  string       // 登录页面的名字，一般是 login.html
	LoginPageParam xgin.PostMap // 给页面的参数

	// Login OK
	LoginRouter      string // default: /login
	LoginRedirectURL string // 登录成功后，跑去哪里

	//
	// Logout
	//

	// Logout OK
	LogoutRouter      string // default: /logout
	LogoutRedirectURL string // 退出登录后，跑去哪里

	//
	// Call this to verify the login parameters
	//
	LoginDataChecker func(c *gin.Context) (sessionItems xgin.PostMap, OKRedirectURL string, NGPageName string, NGPageParam xgin.PostMap, err error)
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
		klog.D("isLoggin: NG: %s", spew.Sdump(o.LoginPageParam))
		c.HTML(200, o.LoginPageName, o.LoginPageParam)
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

func (o *KLogin) doLogin(c *gin.Context) {
	if sessionItems, OKRedirectURL, NGPageName, NGPageParam, err := o.LoginDataChecker(c); err == nil {
		if OKRedirectURL == "" {
			OKRedirectURL = o.LoginRedirectURL
		}

		session := sessions.Default(c)
		for k, v := range sessionItems {
			session.Set(k, v)
		}

		session.Set("UUID", uuid.NewV4().String())
		session.Save()

		c.Redirect(302, OKRedirectURL)
	} else {
		if NGPageName == "" {
			NGPageName = o.LoginPageName
		}

		what := gin.H{}
		if o.LoginPageParam != nil {
			for k, v := range o.LoginPageParam {
				what[k] = v
			}
		}
		if NGPageParam != nil {
			for k, v := range NGPageParam {
				what[k] = v
			}
		}

		c.HTML(200, NGPageName, &what)
	}
}

func (o *KLogin) doLogout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.Redirect(302, o.LogoutRedirectURL)
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

	// Login
	if o.LoginRouter == "" {
		o.LoginRouter = "/login"
	}
	Gin.POST(o.LoginRouter, o.doLogin)

	// Logout
	if o.LogoutRouter == "" {
		o.LogoutRouter = "/logout"
	}
	Gin.GET(o.LogoutRouter, o.doLogout)
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
