package klogin

import (
	"fmt"
	"os"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"

	"github.com/kamasamikon/miego/klog"
	"github.com/kamasamikon/miego/pong"
	"github.com/kamasamikon/miego/xgin"
	"github.com/kamasamikon/miego/xmap"
)

var redisStore redis.Store

// HTTP Response header: yes, no, auto, ...
const (
	// Header: Need Login
	// $ echo -n NEED-LOGIN | md5sum
	// 26bb6bc34e2e91420e7d0a8a522d26f8  -
	H_NEED_LOGIN = "K-26bb6bc34e2e91420e7d0a8a522d26f8"

	// Header: Login URL
	// $ echo -n LOGIN-URL | md5sum
	// 76cf351b86b71ce0e5c514fc520e26f2  -
	H_LOGIN_URL = "K-76cf351b86b71ce0e5c514fc520e26f2"
)

type Login interface {
	BeforeLogout(c *gin.Context) (LogoutRedirectURL string)
	BeforeLogin(c *gin.Context) (StatusCode int, PageName string, PageParam xmap.Map)
	LoginDataChecker(c *gin.Context) (sessionItems xmap.Map, OKRedirectURL string, NGPageName string, NGPageParam xmap.Map, noPageMode bool, err error)

	LoginRouter() []string
	LogoutRouter() []string
}

type LoginCenter struct {
	Gin *gin.Engine

	//
	// Session
	//
	SessionName string
	Session     gin.HandlerFunc

	//
	// Before and After check.
	//
	BCheckerList []func(h gin.HandlerFunc) gin.HandlerFunc
	ACheckerList []func(h gin.HandlerFunc) gin.HandlerFunc

	//
	// Router VS LoginType
	//
	// "GET@/wx/xxx" => "wx"
	// "POST@/user/add" => "ht"
	//
	MapRouterVsLogin map[string]string // Method@Path => LoginType
	MapLogin         map[string]Login  // LoginType => LoginConfigure
}

func (o *LoginCenter) Register(Type string, login Login) {
	o.MapLogin[Type] = login
}

func (o *LoginCenter) SetLoginType(LoginType string, Method string, fullPath string) {
	Key := fmt.Sprintf("%s@%s", Method, fullPath)
	o.MapRouterVsLogin[Key] = LoginType
}

func (o *LoginCenter) GetLoginType(c *gin.Context) string {
	Method := c.Request.Method
	fullPath := c.FullPath()

	Key := fmt.Sprintf("%s@%s", Method, fullPath)
	LoginType, _ := o.MapRouterVsLogin[Key]
	return LoginType
}

func (o *LoginCenter) isLoggin(h gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		//
		// Check Kooky first
		//
		if Kooky := c.Query("Kooky"); Kooky != "" {
			c.Header("Cookie", fmt.Sprintf("%s=%s", o.SessionName, Kooky))
		}
		if Kooky := c.GetHeader("Kooky"); Kooky != "" {
			c.Header("Cookie", fmt.Sprintf("%s=%s", o.SessionName, Kooky))
		}

		session := sessions.Default(c)

		LoginType := o.GetLoginType(c)
		c.Set("LoginType", LoginType)
		if LoginType != "" {
			Type := session.Get(LoginType)
			if Type != nil {
				h(c)
				return
			}
		}

		l := o.MapLogin[LoginType]
		if l != nil {
			c.Set("LoginType", LoginType)
			c.Header(H_LOGIN_URL, l.LoginRouter()[0])

			// Return Status or Login page
			StatusCode, LoginPageName, LoginPageParam := l.BeforeLogin(c)
			if LoginPageName == "" {
				c.JSON(StatusCode, LoginPageParam)
			} else {
				c.HTML(StatusCode, LoginPageName, LoginPageParam)
			}
		}
	}
}

func (o *LoginCenter) Get(c *gin.Context, key string) string {
	if s, ok := c.Get(sessions.DefaultKey); ok {
		if session := s.(sessions.Session); session != nil {
			if val := session.Get(key); val != nil {
				return val.(string)
			}
		}
	}
	return ""
}

func (o *LoginCenter) Set(c *gin.Context, key string, val interface{}) {
	if s, ok := c.Get(sessions.DefaultKey); ok {
		if session := s.(sessions.Session); session != nil {
			session.Set(key, val)
		}
	}
}

func (o *LoginCenter) Save(c *gin.Context) {
	if s, ok := c.Get(sessions.DefaultKey); ok {
		if session := s.(sessions.Session); session != nil {
			session.Save()
		}
	}
}

func (o *LoginCenter) doLogin(c *gin.Context) {
	session := sessions.Default(c)

	LoginType := o.GetLoginType(c)

	l := o.MapLogin[LoginType]
	if l != nil {
		sessionItems, OKRedirectURL, NGPageName, NGPageParam, noPageMode, err := l.LoginDataChecker(c)
		if err == nil {
			var Keys []string
			for k, v := range sessionItems {
				session.Set(k, v)
				Keys = append(Keys, k)
			}
			session.Set(LoginType, strings.Join(Keys, ";"))
			session.Save()

			var cookie string
			if _, rstore := redis.GetRedisStore(redisStore); rstore != nil {
				if ses, _ := rstore.Get(c.Request, o.SessionName); ses != nil {
					cookie, _ = securecookie.EncodeMulti(ses.Name(), ses.ID, rstore.Codecs...)
					c.Header("Kooky", cookie)
				}
			}

			if noPageMode {
				pong.OK(c, xmap.Make("CookieKey", o.SessionName, "CookieVal", cookie))
			} else {
				c.Redirect(302, OKRedirectURL)
			}
		} else {
			session.Delete(LoginType)
			session.Save()

			if noPageMode {
				pong.NG(c, 200, -1, NGPageParam)
			} else {
				c.HTML(200, NGPageName, NGPageParam)
			}
		}
	}
}

func (o *LoginCenter) doLogout(c *gin.Context) {
	session := sessions.Default(c)

	LoginType := o.GetLoginType(c)
	l := o.MapLogin[LoginType]
	if l != nil {
		LogoutRedirectURL := l.BeforeLogout(c)

		session.Delete(LoginType)
		session.Save()

		c.Redirect(302, LogoutRedirectURL)
	}
}

func (o *LoginCenter) Setup(Gin *gin.Engine, SessionName string, redisAddr string) {
	if o.Gin != nil {
		return
	}

	// Gin
	o.Gin = Gin
	o.SessionName = SessionName

	// Redis/Session
	if redisAddr == "" {
		redisHost := os.Getenv("DOCKER_GATEWAY")
		redisAddr = redisHost + ":6379"
	}
	store, err := redis.NewStore(10, "tcp", redisAddr, "", []byte("secret"))
	if err != nil {
		klog.E(err.Error())
		return
	}
	redisStore = store

	if _, rstore := redis.GetRedisStore(redisStore); rstore != nil {
		// ten years
		rstore.SetMaxAge(3600 * 24 * 365 * 10)
	}

	o.Session = sessions.Sessions(o.SessionName, store)
	Gin.Use(o.Session)

	var Key string
	for LoginType, l := range o.MapLogin {
		for _, URL := range l.LoginRouter() {
			Key = fmt.Sprintf("%s@%s", "GET", URL)
			o.MapRouterVsLogin[Key] = LoginType
			Key = fmt.Sprintf("%s@%s", "POST", URL)
			o.MapRouterVsLogin[Key] = LoginType

			POST("", URL, o.doLogin)
			GET("", URL, o.doLogin)
		}
		for _, URL := range l.LogoutRouter() {
			Key = fmt.Sprintf("%s@%s", "GET", URL)
			o.MapRouterVsLogin[Key] = LoginType
			Key = fmt.Sprintf("%s@%s", "POST", URL)
			o.MapRouterVsLogin[Key] = LoginType

			POST("", URL, o.doLogout)
			GET("", URL, o.doLogout)
		}
	}
}

func (o *LoginCenter) POST(LoginType string, relativePath string, handler gin.HandlerFunc) {
	var decors []func(h gin.HandlerFunc) gin.HandlerFunc

	if LoginType != "" {
		o.SetLoginType(LoginType, "POST", relativePath)

		if o.BCheckerList != nil {
			decors = append(decors, o.BCheckerList...)
		}
		decors = append(decors, o.isLoggin)
		if o.ACheckerList != nil {
			decors = append(decors, o.ACheckerList...)
		}
	}
	o.Gin.POST(relativePath, xgin.Decorator(handler, decors...))
}
func (o *LoginCenter) GET(LoginType string, relativePath string, handler gin.HandlerFunc) {
	var decors []func(h gin.HandlerFunc) gin.HandlerFunc

	if LoginType != "" {
		o.SetLoginType(LoginType, "GET", relativePath)

		if o.BCheckerList != nil {
			decors = append(decors, o.BCheckerList...)
		}
		decors = append(decors, o.isLoggin)
		if o.ACheckerList != nil {
			decors = append(decors, o.ACheckerList...)
		}
	}

	o.Gin.GET(relativePath, xgin.Decorator(handler, decors...))
}
func (o *LoginCenter) HEAD(LoginType string, relativePath string, handler gin.HandlerFunc) {
	var decors []func(h gin.HandlerFunc) gin.HandlerFunc

	if LoginType != "" {
		o.SetLoginType(LoginType, "HEAD", relativePath)

		if o.BCheckerList != nil {
			decors = append(decors, o.BCheckerList...)
		}
		decors = append(decors, o.isLoggin)
		if o.ACheckerList != nil {
			decors = append(decors, o.ACheckerList...)
		}
	}
	o.Gin.HEAD(relativePath, xgin.Decorator(handler, decors...))
}
func (o *LoginCenter) OPTIONS(LoginType string, relativePath string, handler gin.HandlerFunc) {
	var decors []func(h gin.HandlerFunc) gin.HandlerFunc

	if LoginType != "" {
		o.SetLoginType(LoginType, "OPTIONS", relativePath)

		if o.BCheckerList != nil {
			decors = append(decors, o.BCheckerList...)
		}
		decors = append(decors, o.isLoggin)
		if o.ACheckerList != nil {
			decors = append(decors, o.ACheckerList...)
		}
	}
	o.Gin.OPTIONS(relativePath, xgin.Decorator(handler, decors...))
}

func (o *LoginCenter) PUT(LoginType string, relativePath string, handler gin.HandlerFunc) {
	var decors []func(h gin.HandlerFunc) gin.HandlerFunc

	if LoginType != "" {
		o.SetLoginType(LoginType, "PUT", relativePath)

		if o.BCheckerList != nil {
			decors = append(decors, o.BCheckerList...)
		}
		decors = append(decors, o.isLoggin)
		if o.ACheckerList != nil {
			decors = append(decors, o.ACheckerList...)
		}
	}

	o.Gin.PUT(relativePath, xgin.Decorator(handler, decors...))

}
func (o *LoginCenter) DELETE(LoginType string, relativePath string, handler gin.HandlerFunc) {
	var decors []func(h gin.HandlerFunc) gin.HandlerFunc

	if LoginType != "" {
		o.SetLoginType(LoginType, "DELETE", relativePath)

		if o.BCheckerList != nil {
			decors = append(decors, o.BCheckerList...)
		}
		decors = append(decors, o.isLoggin)
		if o.ACheckerList != nil {
			decors = append(decors, o.ACheckerList...)
		}
	}

	o.Gin.DELETE(relativePath, xgin.Decorator(handler, decors...))
}

var Default *LoginCenter

func init() {
	Default = &LoginCenter{}
	Default.MapRouterVsLogin = make(map[string]string)
	Default.MapLogin = make(map[string]Login)
}
