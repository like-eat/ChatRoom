package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 模拟鉴权中间件
func Auth() gin.HandlerFunc {
	return func(c *gin.Context){
		if c.GetHeader("Authorization") == "" {
			c.JSON(401, gin.H{"error": "需要登录"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// 模拟管理员中间件
func AdminOnly() gin.HandlerFunc {
	
}