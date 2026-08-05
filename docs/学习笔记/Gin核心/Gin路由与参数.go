package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func main(){
	r := gin.Default()

	// Http方法
	r.Get("/ping", func(c *gin.Context)){
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	}

	r.POST("/submit", func(c *gin.Context)){
		c.JSON(http.StatusOK, gin.H{
			"message": "post ok",
		})
	}

	// 路径参数 user/:id
	r.GET("/user/:id", func(c *gin.Context)){
		id := c.Param("id") // 获取路径上的参数
		c.JSON(http.StatusOK, gin.H{
			"user_id": id
		})
	}

	// 查询参数 /search?q=xxx&page=1
	r.GET("/search", func(c *gin.Context)){
		q := c.Query("q") // 获取查询参数
		page := c.Query("page")
		c.JSON(http.StatusOK, gin.H{
			"q": q,
			"page": page,
		})
	}

	r.run(":8080") // 启动服务
}