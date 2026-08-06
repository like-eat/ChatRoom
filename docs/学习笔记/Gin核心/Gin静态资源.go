package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// 映射单个目录
	r.Static("/static", "./public")

	// 映射到文件系统
	r.StaticFS("/assets", gin.Dir("./assets", true))

	// 单个文件
	r.StaticFile("/favicon.ico", "./public/favicon.ico")

	// 自定义文件服务器
	r.GET("/download/:filename", func(c *gin.Context) {
		filename := c.Param("filename")
		
		c.File("./public/" + filename)
	})

	// 完整实例
	r.POST("/upload/avatar", func(c *gin.Context){
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(400, gin.H{"error": "没有文件"})
			return
		}
		c.SaveUploadedFile(file, "./public/avatar" + file.Filename)
		c.JSON(200, gin.H{
			"url": "/static/avatar" + file.Filename,
		})
	})

	r.Run(":8080")
}