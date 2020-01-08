package xgin

import (
	"github.com/gin-gonic/gin"
)

// JPong : JSON Pong
type JPong struct {
	Error   int         `json:"Error"`
	Message string      `json:"Message"`
	Data    interface{} `json:"Data"`
}

func J(c *gin.Context, Code int, Error int, Message string, Data interface{}) {
	if Data == nil {
		Data = &gin.H{}
	}

	c.JSON(Code, &gin.H{
		"Error":   Error,
		"Message": Message,
		"Data":    Data,
	})
}
