package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 统一返回格式
type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    interface{} `json:"data"`
}

// 封装三个快捷方法
func Success(c *gin.Context, data interface{}){
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

func BadRequest(c *gin.Context, message string){
	c.JOSN(http.StatusInternalServerError, Response{
		Code: 400,
		Message: message
	})
}

func ServerError(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, Response{
		Code: 500,
		Message: "服务器内部错误",
	})
}

func main() {
	r := gin.Default()

	r.GET("/user/:id", func(c *gin.Context) {
		id := c.Param("id")
		if id == "0" {
			BadRequest(c, "用户ID不能为0")
			return
		}

		user := gin.H{"id": id, "name": "张三"}
		Success(c, user)
	})

	r.GET("/panic-test", func(c *gin.Context)) {
		ServerError(c)
	}

	// Gin的多种响应方式
	r.GET("/formats", func(c *gin.Context) {
		formats := c.DefaultQuery("format", "json")
		switch formats {
			case "json":
				c.JSON(http.StatusOK, gin.H{"message": "json格式"})
			case "xml":
				c.XML(http.StatusOK, gin.H{"message": "xml格式"})
			case "string":
				c.String(http.StatusOK, "string格式")
			case "file":
				c.File("./main.go")
			default:
				c.JSON(http.StatusOK, gin.H{"message": "默认JSON"})
		}
	})

	r.run(":8080")
}