package pong

import (
	"fmt"
	"runtime"

	"github.com/gin-gonic/gin"
)

// JPong : JSON Pong
type Body struct {
	Error   int         `json:"Error"`
	Message string      `json:"Message,omitempty"`
	Data    interface{} `json:"Data,omitempty"`
}

const (
	E_OK = 0

	//
	// 环境变量或者别的变量的错误
	//
	// 格式错误：比如手机号，Param里是参数的名字或者索引
	E_BadFormat = -100

	// 没有找到：参数里有，但系统里没有找到
	E_NotFound = -101

	// 不能为空：参数里有，但是空的
	E_Empty = -102

	// 没有设置：参数里要求有，但没有
	E_NotExists = -103

	// 匹配错误：一个参数和另外一个参数有匹配的规则
	E_BadMatch = -104

	// 参数错误：有，但是是错的，这个错都不确定了
	E_Error = -105

	//
	// 输入的参数的错误
	//

	// 格式错误：比如手机号，Param里是参数的名字或者索引
	E_Para_BadFormat = -1000

	// 没有找到：参数里有，但系统里没有找到
	E_Para_NotFound = -1001

	// 不能为空：参数里有，但是空的
	E_Para_Empty = -1002

	// 没有设置：参数里要求有，但没有
	E_Para_NotExists = -1003

	// 匹配错误：一个参数和另外一个参数有匹配的规则
	E_Para_BadMatch = -1004

	// 参数错误：有，但是是错的，这个错都不确定了
	E_Para_Error = -1005

	//
	// 权限相关
	//

	// 未登录：请重新登录
	E_Perm_NotLogin = -2

	// 未授权：比如所在的组不对等。Role = 组，Orgn = 所在机构，Oper = 后台账户
	E_Perm_NotAllow = -3
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
	Full(c, 200, E_OK, "", Data)
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
func NG_Para_BadFormat(c *gin.Context, Param string) {
	Full(c, 200, E_Para_BadFormat, "E_Para_BadFormat", Param)
}

// 没有找到：参数里有，但系统里没有找到，这个和参数错误区分不开，因为参数错误也导致目标对象找不到
func NG_Para_NotFound(c *gin.Context, Param string) {
	Full(c, 200, E_Para_NotFound, "E_Para_NotFound", Param)
}

// 不能为空：参数里有，但是空的
func NG_Para_Empty(c *gin.Context, Param string) {
	Full(c, 200, E_Para_Empty, "E_Para_Empty", Param)
}

// 没有设置：参数里要求有，但没有
func NG_Para_NotExists(c *gin.Context, Param string) {
	Full(c, 200, E_Para_NotExists, "E_Para_NotExists", Param)
}

// 匹配错误：一个参数和另外一个参数有匹配的规则
func NG_Para_BadMatch(c *gin.Context, Param string) {
	Full(c, 200, E_Para_BadMatch, "E_Para_BadMatch", Param)
}

// 参数错误：有，但是是错的，这个错都不确定了
func NG_Para_Error(c *gin.Context, Param string) {
	Full(c, 200, E_Para_Error, "E_Para_Error", Param)
}

// Environment/Context 上下文的错误，比如环境变量，比如会话
//
// Param: Param里是参数的名字或者索引
//
// 格式错误：带了这个参数，但格式是错的，比如手机号，写成了12345
func NG_BadFormat(c *gin.Context, Param string) {
	Full(c, 200, E_BadFormat, "E_BadFormat", Param)
}

// 没有找到：参数里有，但系统里没有找到，这个和参数错误区分不开，因为参数错误也导致目标资源找不到
func NG_NotFound(c *gin.Context, Param string) {
	Full(c, 200, E_NotFound, "E_NotFound", Param)
}

// 不能为空：参数里有，但是空的
func NG_Empty(c *gin.Context, Param string) {
	Full(c, 200, E_Empty, "E_Empty", Param)
}

// 没有设置：参数里要求有，但没有
func NG_NotExists(c *gin.Context, Param string) {
	Full(c, 200, E_NotExists, "E_NotExists", Param)
}

// 匹配错误：一个参数和另外一个参数有匹配的规则
func NG_BadMatch(c *gin.Context, Param string) {
	Full(c, 200, E_BadMatch, "E_BadMatch", Param)
}

// 参数错误：有，但是是错的，这个错都不确定了
func NG_Error(c *gin.Context, Param string) {
	Full(c, 200, E_Error, "E_Error", Param)
}

//
// 权限相关
//

// 未登录：请重新登录
func NG_Perm_NotLogin(c *gin.Context) {
	Full(c, 200, E_Perm_NotLogin, "E_Perm_NotLogin", "")
}

// 未授权：比如所在的组不对等。Role = 组，Orgn = 所在机构，Oper = 后台账户
func NG_Perm_NotAllow(c *gin.Context, Role string) {
	Full(c, 200, E_Perm_NotAllow, "E_Perm_NotAllow", Role)
}
