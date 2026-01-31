package pong

// An HTTP response in JSON format.

import (
	"fmt"
	"runtime"

	"github.com/gin-gonic/gin"
)

// Body: JSON Pong
type Body struct {
	Error   int         `json:"Error"`
	Message string      `json:"Message,omitempty"`
	Data    interface{} `json:"Data,omitempty"`
}

const (
	EOK = 0

	//
	// 环境变量或者别的变量的错误
	//
	// 格式错误：比如手机号，Param里是参数的名字或者索引
	EBadFormat = -100

	// 没有找到：参数里有，但系统里没有找到
	ENotFound = -101

	// 不能为空：参数里有，但是空的
	EEmpty = -102

	// 没有设置：参数里要求有，但没有
	ENotExists = -103

	// 匹配错误：一个参数和另外一个参数有匹配的规则
	EBadMatch = -104

	// 参数错误：有，但是是错的，这个错都不确定了
	EError = -105

	//
	// 输入的参数的错误
	//

	// 格式错误：比如手机号，Param里是参数的名字或者索引
	EParaBadFormat = -1000

	// 没有找到：参数里有，但系统里没有找到
	EParaNotFound = -1001

	// 不能为空：参数里有，但是空的
	EParaEmpty = -1002

	// 没有设置：参数里要求有，但没有
	EParaNotExists = -1003

	// 匹配错误：一个参数和另外一个参数有匹配的规则
	EParaBadMatch = -1004

	// 参数错误：有，但是是错的，这个错都不确定了
	EParaError = -1005

	//
	// 权限相关
	//

	// 未登录：请重新登录
	EPermNotLogin = -2

	// 未授权：比如所在的组不对等。Role = 组，Orgn = 所在机构，Oper = 后台账户
	EPermNotAllow = -3
)

func Full(c *gin.Context, Code int, Error int, Message interface{}, Data interface{}) {
	var Text string
	if s, ok := Message.(error); ok {
		Text = s.Error()
	} else if s, ok := Message.(string); ok {
		Text = s
	} else {
		Text = ""
	}

	if Error != 0 && Data == nil {
		if _, filename, line, ok := runtime.Caller(2); ok {
			Data = fmt.Sprintf("%s:%d", filename, line)
			// Data = des.Jia(Data.(string), "HILDA")
		} else {
			Data = &gin.H{}
		}
	}

	c.JSON(Code, &Body{
		Error:   Error,
		Message: Text,
		Data:    Data,
	})
}

func OK(c *gin.Context, Data interface{}) {
	Full(c, 200, EOK, "", Data)
}

func NG(c *gin.Context, Code int, Error int, Message interface{}) {
	Full(c, Code, Error, Message, nil)
}

//
// Parameters error
//
// Param: Param里是参数的名字或者索引
//

// 格式错误：带了这个参数，但格式是错的，比如手机号，写成了12345
func NGParaBadFormat(c *gin.Context, Param string) {
	Full(c, 200, EParaBadFormat, "EParaBadFormat", Param)
}

// 没有找到：参数里有，但系统里没有找到，这个和参数错误区分不开，因为参数错误也导致目标对象找不到
func NGParaNotFound(c *gin.Context, Param string) {
	Full(c, 200, EParaNotFound, "EParaNotFound", Param)
}

// 不能为空：参数里有，但是空的
func NGParaEmpty(c *gin.Context, Param string) {
	Full(c, 200, EParaEmpty, "EParaEmpty", Param)
}

// 没有设置：参数里要求有，但没有
func NGParaNotExists(c *gin.Context, Param string) {
	Full(c, 200, EParaNotExists, "EParaNotExists", Param)
}

// 匹配错误：一个参数和另外一个参数有匹配的规则
func NGParaBadMatch(c *gin.Context, Param string) {
	Full(c, 200, EParaBadMatch, "EParaBadMatch", Param)
}

// 参数错误：有，但是是错的，这个错都不确定了
func NGParaError(c *gin.Context, Param string) {
	Full(c, 200, EParaError, "EParaError", Param)
}

// Environment/Context 上下文的错误，比如环境变量，比如会话
//
// Param: Param里是参数的名字或者索引
//
// 格式错误：带了这个参数，但格式是错的，比如手机号，写成了12345
func NGBadFormat(c *gin.Context, Param string) {
	Full(c, 200, EBadFormat, "EBadFormat", Param)
}

// 没有找到：参数里有，但系统里没有找到，这个和参数错误区分不开，因为参数错误也导致目标资源找不到
func NGNotFound(c *gin.Context, Param string) {
	Full(c, 200, ENotFound, "ENotFound", Param)
}

// 不能为空：参数里有，但是空的
func NGEmpty(c *gin.Context, Param string) {
	Full(c, 200, EEmpty, "EEmpty", Param)
}

// 没有设置：参数里要求有，但没有
func NGNotExists(c *gin.Context, Param string) {
	Full(c, 200, ENotExists, "ENotExists", Param)
}

// 匹配错误：一个参数和另外一个参数有匹配的规则
func NGBadMatch(c *gin.Context, Param string) {
	Full(c, 200, EBadMatch, "EBadMatch", Param)
}

// 参数错误：有，但是是错的，这个错都不确定了
func NGError(c *gin.Context, Param string) {
	Full(c, 200, EError, "EError", Param)
}

//
// 权限相关
//

// 未登录：请重新登录
func NGPermNotLogin(c *gin.Context) {
	Full(c, 200, EPermNotLogin, "EPermNotLogin", "")
}

// 未授权：比如所在的组不对等。Role = 组，Orgn = 所在机构，Oper = 后台账户
func NGPermNotAllow(c *gin.Context, Role string) {
	Full(c, 200, EPermNotAllow, "EPermNotAllow", Role)
}
