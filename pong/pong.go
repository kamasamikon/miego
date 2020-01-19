package pong

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"runtime"
)

// JPong : JSON Pong
type Body struct {
	Error   int         `json:"Error"`
	Message string      `json:"Message"`
	Data    interface{} `json:"Data"`
}

func Full(c *gin.Context, Code int, Error int, Message interface{}, Data interface{}) {
	if Data == nil {
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
		"Error":   Error,
		"Message": Text,
		"Data":    Data,
	})
}

func OK(c *gin.Context, Data interface{}) {
	Full(c, 200, 0, "", Data)
}

func NG(c *gin.Context, Code int, Error int, Message interface{}) {
	Full(c, Code, Error, Message, nil)
}
