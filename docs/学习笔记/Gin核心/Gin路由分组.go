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
// gin.HandlerFunc 是中间件的类型
// gin.Context 携带了一个请求的完整上下文，承载http的全部信息
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetHeader("X-Role")
		if role != "admin": {
			c.JSON(403, gin.H{"error": "需要管理员权限"})
			c.Abort()
			return
		}
		// c.Next() 中间件链的接力棒，执行完这个中间件继续执行下一个
		c.Next()
	}
}

func main() {
	r := gin.Default()

    // v1 版本
	v1 := r.Group("/api/v1")
	{
		// 公开接口
		v1.POST("/login", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "登录成功"})
		})
		v1.POST("/register", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "注册成功"})
		})

		// 需要登录的接口
		auth := v1.Group("") // 创建路由分组
		auth.Use(Auth()) // 并给分组挂载中间件
		{
			auth.GET("/user/info", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "张三"})
			})
			auth.PUT("/user/update", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
			})

			// 管理员接口
			admin := auth.Group("admin")
			admin.Use(AdminOnly()) // 并给分组挂载中间件
			{
				admin.GET("/users", func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{"users": []string{"张三", "李四", "王五"}})
				})
				admin.POST("/users/disable", func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{"message": "用户禁用成功"})
				})
			}
	 	}
	}

	// v2 版本
	v2 := r.Group("/api/v2") 
	{
		v2.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})
	}

	r.Run(":8080")
}