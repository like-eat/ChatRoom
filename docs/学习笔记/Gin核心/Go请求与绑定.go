package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// binding tag 用于验证
type LoginRequest struct {
	Telephone string `form:"telephone" binding:"required,len=11"`
	Password  string `form:"password" binding:"required,min=6,max=12"`
	Remember   bool   `form:"remember"`
}

// form tag 用于表单绑定
type RegisterRequest struct {
	Nickname string `form:"nickname" binding:"required"`
	Telephone string `form:"telephone"  binding:"required,len=11"`
	Password  string `form:"password" binding:"required"`
}

func main() {
	r := gin.Default()

	// JSON绑定＋验证
	r.POST("/login", func(c *gin.Context) {
		var req LoginRequest

		// 
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return 
		}

		// 绑定成功后
		c.JSON(http.StatusOK, gin.H{
			"msg": "登陆成功",
			"telephone": req.Telephone,
		})
	})

	// 表单绑定
	r.POST("/register", func(c *gin.Context) {
		var req RegisterRequest

		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return 
		}

		c.JSON(http.StatusOK, gin.H{
			"msg": "注册成功",
			"nickname": req.Nickname,
			"telephone": req.Telephone,
		})
	})

	r.run(":8080")
}
