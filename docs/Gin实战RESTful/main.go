package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// 初始化数据，添加两条演示用户
	store.CreateUser(RegisterReq{
		Nickname: "张三", Telephone: "13800001111",
		Password: "123456", Gender: 0,
	})
	store.CreateUser(RegisterReq{
		Nickname: "李四", Telephone: "13800001112",
		Password: "123456", Gender: 1,
	})

	// 公开路由
	r.POST("/register", Register)
	r.POST("/login", Login)

	// 需要登录的路由
	auth := r.Group("/api").Use(JWTAuth())
	{
		// 获取自己的信息
		auth.GET("/profile", GetProfile)
		// 更新自己的信息
		auth.PUT("/profile", UpdateProfile)
		// 删除自己的账号
		auth.DELETE("/user/:id", DeleteUser)
	}

	// 公开的查询路由
	r.GET("/users", ListUsers)       // 用户列表分页
	r.GET("/users/:id", GetUserByID) // 获取单个用户

	fmt.Println("🚀 服务启动在 http://127.0.0.1:8080")
	fmt.Println("📋 API 列表：")
	fmt.Println("  POST  /register      注册")
	fmt.Println("  POST  /login         登录")
	fmt.Println("  GET   /api/profile   获取个人信息（需 Token）")
	fmt.Println("  PUT   /api/profile   更新个人信息（需 Token）")
	fmt.Println("  DELETE /api/user/:id 删除用户（需 Token，只能删自己）")
	fmt.Println("  GET   /users?page=1&page_size=10  用户列表")
	fmt.Println("  GET   /user/:id      查看用户详情")

	r.Run(":8080")
}
