package klogin

import (
	"github.com/gin-gonic/gin"
)

var loginMap map[string]Login

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
			Default.POST(rp.LoginType, rp.RelativePath, rp.Handler)
		case "GET":
			Default.GET(rp.LoginType, rp.RelativePath, rp.Handler)
		case "PUT":
			Default.PUT(rp.LoginType, rp.RelativePath, rp.Handler)
		case "DELETE":
			Default.DELETE(rp.LoginType, rp.RelativePath, rp.Handler)
		}
	}
}

func Get(c *gin.Context, key string) interface{} {
	return Default.Get(c, key)
}

func Set(c *gin.Context, key string, val interface{}) {
	Default.Set(c, key, val)
}

func Save(c *gin.Context) {
	Default.Save(c)
}

func Register(Type string, login Login) {
	loginMap[Type] = login
}

func Setup(Gin *gin.Engine, SessionName string) {
	for k, v := range loginMap {
		Default.Register(k, v)
	}
	Default.Setup(Gin, SessionName)
}

// sessionItem
func SessionGet(c *gin.Context, name string) string {
	return Default.Get(c, name)
}

func LoginType(c *gin.Context) string {
	return Default.GetLoginType(c)
}

func init() {
	loginMap = make(map[string]Login)
}
