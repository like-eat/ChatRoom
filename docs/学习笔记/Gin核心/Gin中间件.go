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
		


	}
}