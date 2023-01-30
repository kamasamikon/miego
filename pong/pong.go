package pong

import (
	"fmt"
	"runtime"

	"github.com/gin-gonic/gin"
)

// JPong : JSON Pong
type Body struct {
	Error   int         `json:"Error"`
	Message string      `json:"Message"`
	Data    interface{} `json:"Data"`
}

func Full(c *gin.Context, Code int, Error int, Message interface{}, Data interface{}) {
	if Error != 0 && Data == nil {
		if _, filename, line, ok := runtime.Caller(2); ok {
			Data = fmt.Sprintf("%s:%d", filename, line)
		} else {
			Data = &gin.H{}
		}
	}

	var Text string
	if s, ok := Message.(error); ok {
		Text = s.Error()
	} else if s, ok := Message.(string); ok {
		Text = s
	} else {
		Text = ""
	}

	c.JSON(Code, &Body{
		Error:   Error,
		Message: Text,
		Data:    Data,
	})
}

func OK(c *gin.Context, Data interface{}) {
	Full(c, 200, 0, "", Data)
}

func NG(c *gin.Context, Code int, Error int, Message interface{}) {
	Full(c, Code, Error, Message, nil)
}

//
// Parameters error
//

// 格式错误：比如手机号，Param里是参数的名字或者索引
func NG_Para_BadFormat(c *gin.Context, Param string) {
	Full(c, 200, -1000, Param, nil)
}

// 没有找到：参数里有，但系统里没有找到
func NG_Para_NotFound(c *gin.Context, Param string) {
	Full(c, 200, -1001, Param, nil)
}

// 不能为空：参数里有，但是空的
func NG_Para_Empty(c *gin.Context, Param string) {
	Full(c, 200, -1002, Param, nil)
}

// 没有设置：参数里要求有，但没有
func NG_Para_NotExists(c *gin.Context, Param string) {
	Full(c, 200, -1003, Param, nil)
}

// 匹配错误：一个参数和另外一个参数有匹配的规则
func NG_Para_BadMatch(c *gin.Context, Param string) {
	Full(c, 200, -1004, Param, nil)
}

// 参数错误：有，但是是错的，这个错都不确定了
func NG_Para_Error(c *gin.Context, Param string) {
	Full(c, 200, -1005, Param, nil)
}

//
// 权限相关
//

// 未登录：请重新登录
func NG_Perm_NotLogin(c *gin.Context) {
	Full(c, 200, -2, "", nil)
}

// 未授权：比如所在的组不对等。Role = 组，Orgn = 所在机构，Oper = 后台账户
func NG_Perm_NotAllow(c *gin.Context, Role string) {
	Full(c, 200, -3, Role, nil)
}
