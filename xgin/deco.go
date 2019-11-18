package xgin

import (
	"github.com/gin-gonic/gin"
)

// http://xion.io/post/code/go-decorated-functions.html
//
// r.POST("/v1/login", Decorator(CheckAuth, CheckToken, Login))
//
// func CheckToken(h gin.HandlerFunc) gin.HandlerFunc {
//     return func(c *gin.Context) {
//         header := c.Request.Header.Get("token")
//         if header == "" {
//             c.JSON(200, gin.H{
//                 "code":   3,
//                 "result": "failed",
//                 "msg":    ". Missing token",
//             })
//             return
//         }
//         h(c)
//     }
// }
//
// func Login(c *gin.Context) {
//     c.JSON(200, gin.H{
//         "code":   0,
//         "result": "success",
//         "msg":    "验证成功",
//     })
// }
func Decorator(h gin.HandlerFunc, decors ...func(gin.HandlerFunc) gin.HandlerFunc) gin.HandlerFunc {
	for i := range decors {
		d := decors[len(decors)-1-i] // iterate in reverse
		h = d(h)
	}
	return h
}
