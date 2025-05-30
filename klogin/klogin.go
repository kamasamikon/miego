package klogin

import (
	"fmt"
	"os"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"

	"miego/klog"
	"miego/pong"
	"miego/xgin"
	"miego/xmap"
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
	// 返回登录页面前调用，如果LoginPageName存在，就返回HTML，否则返回JSON
	BeforeLogin(c *gin.Context) (StatusCode int, PageName string, PageParam xmap.Map)
	// 退出（删除Session）前调用，返回跳转的URL
	BeforeLogout(c *gin.Context) (LogoutRedirectURL string)

	// ok: 成功还是失败
	// noPageMode: 返回页面还是返回JSON
	//
	// okData["noPageData"]: Map: Pong中返回 pong.OK(c, noPageData)
	// okData["redirectURL"]: Str: 跳转页面
	// okData["sessionItems"]: Map: 保存到会话
	//
	// ngData["noPageData"]: Map: Pong中返回  pong.NG(c, Error, Data)
	// ngData["noPageData"][Error"]: Int:
	// ngData["noPageData"][Data"]: Obj: 一般就是Str
	// ngData["templateName"]: Str: c.HTML 需要的模板
	// ngData["templateData"]: Map: c.HTML 需要的数据
	LoginDataChecker(c *gin.Context) (ok bool, noPageMode bool, okData xmap.Map, ngData xmap.Map)

	LoginRouter() []string
	LogoutRouter() []string

	AfterLogin(c *gin.Context, cookie string) // 登录动作成功了
	AfterLogout(c *gin.Context)               // 退出动作成功了
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

	LoginCheckerHooker func(o *LoginCenter, c *gin.Context)
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

// 检查每一个调用，看看是否已经登录了，如果没有登录，就跳转到登录接口
func (o *LoginCenter) loginChecker(h gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if o.LoginCheckerHooker != nil {
			o.LoginCheckerHooker(o, c)
		}

		//
		// Check Kooky first
		//
		if v, _ := c.Get("Kooky"); v != nil { // 来自 LoginCheckerHooker
			Kooky := v.(string)
			c.Request.Header.Set("Cookie", fmt.Sprintf("%s=%s", o.SessionName, Kooky))
		}
		if Kooky := c.Query("Kooky"); Kooky != "" {
			c.Request.Header.Set("Cookie", fmt.Sprintf("%s=%s", o.SessionName, Kooky))
		}
		if Kooky := c.GetHeader("Kooky"); Kooky != "" {
			c.Request.Header.Set("Cookie", fmt.Sprintf("%s=%s", o.SessionName, Kooky))
		}

		LoginType := o.GetLoginType(c)
		klog.D("%v", LoginType)
		c.Set("LoginType", LoginType)
		if LoginType != "" {
			session := sessions.Default(c)
			Type := session.Get(LoginType)
			if Type != nil {
				h(c)
				return
			} else {
				klog.E("Type is none for LoginType %s", LoginType)
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

func (o *LoginCenter) Rem(c *gin.Context, key string) {
	if s, ok := c.Get(sessions.DefaultKey); ok {
		if session := s.(sessions.Session); session != nil {
			session.Delete(key)
		}
	}
}
func (o *LoginCenter) Clr(c *gin.Context) {
	if s, ok := c.Get(sessions.DefaultKey); ok {
		if session := s.(sessions.Session); session != nil {
			session.Clear()
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
		ok, noPageMode, okData, ngData := l.LoginDataChecker(c)
		if ok {
			var Keys []string
			for k, v := range okData["sessionItems"].(xmap.Map) {
				session.Set(k, v)
				Keys = append(Keys, k)
			}
			session.Set(LoginType, strings.Join(Keys, ";")) // 没有什么意义，就是能看到保存了哪些字段
			session.Save()

			var cookie string
			if _, rstore := redis.GetRedisStore(redisStore); rstore != nil {
				if ses, _ := rstore.Get(c.Request, o.SessionName); ses != nil {
					cookie, _ = securecookie.EncodeMulti(ses.Name(), ses.ID, rstore.Codecs...)
					c.Header("Kooky", cookie)
				}
			}

			// 看看有什么要补充的
			l.AfterLogin(c, cookie)

			if noPageMode {
				noPageData := okData["noPageData"].(xmap.Map)
				noPageData.SafePut("CookieKey", o.SessionName, "CookieVal", cookie)
				pong.OK(c, noPageData)
			} else {
				c.Redirect(302, okData["redirectURL"].(string))
			}
		} else {
			session.Delete(LoginType)
			session.Save()

			if noPageMode {
				noPageData := ngData["noPageData"].(xmap.Map)
				pong.NG(c, 200, noPageData["Error"].(int), noPageData["Data"])
			} else {
				c.HTML(200, ngData["templateName"].(string), ngData["templateData"].(xmap.Map))
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

		// 看看有什么要补充的
		l.AfterLogout(c)

		if LogoutRedirectURL != "" {
			c.Redirect(302, LogoutRedirectURL)
		} else {
			pong.OK(c, "OK")
		}
	} else {
		pong.NG(c, 404, -1, "Not Login")
	}
}

func (o *LoginCenter) Setup(Gin *gin.Engine, SessionName string, redisAddr string) {
	if o.Gin != nil {
		return
	}

	// Gin
	o.Gin = Gin

	if SessionName != "" {
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
	}

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

func (o *LoginCenter) Route(Methods string, LoginTypes string, relativePath string, handler gin.HandlerFunc) {
	LoginTypeList := strings.Split(LoginTypes, ",")

	MethodList := strings.Split(Methods, ",")
	for _, Method := range MethodList {
		Method = strings.TrimSpace(Method)
		switch Method {
		case "POST":
			for _, LoginType := range LoginTypeList {
				LoginType = strings.TrimSpace(LoginType)
				klog.D("Method:%v, LoginType:%v", Method, LoginType)
				o.POST(LoginType, relativePath, handler)
			}
		case "GET":
			for _, LoginType := range LoginTypeList {
				LoginType = strings.TrimSpace(LoginType)
				klog.D("Method:%v, LoginType:%v", Method, LoginType)
				o.GET(LoginType, relativePath, handler)
			}
		case "HEAD":
			for _, LoginType := range LoginTypeList {
				LoginType = strings.TrimSpace(LoginType)
				o.HEAD(LoginType, relativePath, handler)
			}
		case "OPTIONS":
			for _, LoginType := range LoginTypeList {
				LoginType = strings.TrimSpace(LoginType)
				o.OPTIONS(LoginType, relativePath, handler)
			}
		case "PUT":
			for _, LoginType := range LoginTypeList {
				LoginType = strings.TrimSpace(LoginType)
				o.PUT(LoginType, relativePath, handler)
			}
		case "DELETE":
			for _, LoginType := range LoginTypeList {
				LoginType = strings.TrimSpace(LoginType)
				o.DELETE(LoginType, relativePath, handler)
			}
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
		decors = append(decors, o.loginChecker)
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
		decors = append(decors, o.loginChecker)
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
		decors = append(decors, o.loginChecker)
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
		decors = append(decors, o.loginChecker)
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
		decors = append(decors, o.loginChecker)
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
		decors = append(decors, o.loginChecker)
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
