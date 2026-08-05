package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// 自定义中间件 1请求日志
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		fmt.Printf("{请求开始前} %s %s \n", c.Request.Method, c.Request.URL.Path)

		c.Next() // 执行下一个中间件

		// next()之后 = 请求处理后
		cost := time.Since(start)
		fmt.Printf("{请求结束后} %s %s %v 状态码：%d \n", c.Request.Method, c.Request.URL.Path,
		cost, c.Writer.Status())
	}
}

// 自定义中间件 2简易鉴权
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token != "Beader secret-token" {
			// 停止不执行后续handler
			c.JSON(401, gin.H{"error": "未授权"})
			c.Abort()
			return
		}

		// 把用户信息存入上下文
		c.Set("user", "user123")
		c.Next()
	}
}

// 自定义中间件 3计算请求耗时
func CostMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		// c.Get 取出前面存入的值
		userID, _ := c.Get("user_id")
		log.Printf("用户 %v 的请求处理完毕", userID)
	}
}

func main() {
	r := gin.Default()

	// 全局中间体
	r.Use(LoggerMiddleware())

	// 公开路由
	r.GET("/public", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "欢迎访问公共路由"})
	})

	// 鉴权路由
	auth := r.Group("/api")
	auth.Use(AuthMiddleware(), CostMiddleware())
	{
		auth.GET("/profile", func(c *gin.Context) {
			userID, _ := c.Get("user_id")
			c.JSON(200, gin.H{"user_id": userID, "message": "张三"})
		})

		auth.GET("/orders", func(c *gin.Contenxt) {
			c.JSON(200, gin.H{"orders": []string{"订单1", "订单2", "订单3"}, "message": "张三的订单"})
		})
	}

	r.Run(":8080")
}