package auth

import (
	"github.com/gin-gonic/gin"

	"github.com/kamasamikon/miego/klogin"
)

var loginMap map[string]klogin.Login

type RouterParam struct {
	LoginType    string
	Method       string
	RelativePath string
	Handler      gin.HandlerFunc
}

var RouterParamList []*RouterParam

func POST(LoginType string, relativePath string, handler gin.HandlerFunc) {
	RouterParamList = append(
		RouterParamList,
		&RouterParam{
			LoginType:    LoginType,
			Method:       "POST",
			RelativePath: relativePath,
			Handler:      handler,
		},
	)
}
func GET(LoginType string, relativePath string, handler gin.HandlerFunc) {
	RouterParamList = append(
		RouterParamList,
		&RouterParam{
			LoginType:    LoginType,
			Method:       "GET",
			RelativePath: relativePath,
			Handler:      handler,
		},
	)
}
func PUT(LoginType string, relativePath string, handler gin.HandlerFunc) {
	RouterParamList = append(
		RouterParamList,
		&RouterParam{
			LoginType:    LoginType,
			Method:       "PUT",
			RelativePath: relativePath,
			Handler:      handler,
		},
	)
}
func DELETE(LoginType string, relativePath string, handler gin.HandlerFunc) {
	RouterParamList = append(
		RouterParamList,
		&RouterParam{
			LoginType:    LoginType,
			Method:       "DELETE",
			RelativePath: relativePath,
			Handler:      handler,
		},
	)
}

func Start() {
	for _, rp := range RouterParamList {
		switch rp.Method {
		case "POST":
			klogin.Default.POST(rp.LoginType, rp.RelativePath, rp.Handler)
		case "GET":
			klogin.Default.GET(rp.LoginType, rp.RelativePath, rp.Handler)
		case "PUT":
			klogin.Default.PUT(rp.LoginType, rp.RelativePath, rp.Handler)
		case "DELETE":
			klogin.Default.DELETE(rp.LoginType, rp.RelativePath, rp.Handler)
		}
	}
}

func Get(c *gin.Context, key string) interface{} {
	return klogin.Default.Get(c, key)
}

func Set(c *gin.Context, key string, val interface{}) {
	klogin.Default.Set(c, key, val)
}

func Save(c *gin.Context) {
	klogin.Default.Save(c)
}

func Register(Type string, login klogin.Login) {
	loginMap[Type] = login
}

func Setup(Gin *gin.Engine, SessionName string) {
	for k, v := range loginMap {
		klogin.Default.Register(k, v)
	}
	klogin.Default.Setup(Gin, SessionName)
}

// sessionItem
func SessionGet(c *gin.Context, name string) string {
	return klogin.Default.Get(c, name)
}

func LoginType(c *gin.Context) string {
	return klogin.Default.GetLoginType(c)
}

func init() {
	loginMap = make(map[string]klogin.Login)
}
