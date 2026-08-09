package main

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// JWT鉴权中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Auth header中获取token
		authHeader := c.GetHeader("Authorization")
		// 取token
		if authHeader == "" {
			Error(c, 401, "缺少 Authorization 头")
			c.Abort()
			return
		}
		// 检查token格式是否正确
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			Error(c, 401, "Authorization 格式错误")
			c.Abort()
			return
		}

		// 验证token
		claims, err := ParseToken(parts[1])
		if err != nil {
			Error(c, 401, "无效的 token")
			c.Abort()
			return
		}

		// 把用户信息存入上下文，后续 handler可用
		c.Set("user_id", claims.UserID)
		c.Set("nickname", claims.Nickname)
		c.Next()
	}
}
