package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// 限制文件上传大小
	r.MaxMultipartMemory = 8 << 20 // 8MB

	// 单文件上传
	r.POST("/upload/avatar", func(c *gin.Context) {
		// 获取上传的文件
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "没有找到文件"})
			return
		}

		// 文件信息
		fmt.Printf("文件名 %s, 大小 %s bytes\n", file.Filename, file.Size)

		// 保存文件到指定目录
		savePath := filepath.Join("./uploads", file.Filename)
		if err := c.SaveUploadedFile(file, savePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件失败"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "文件上传成功",
			"filename":    file.Filename,
			"size":        file.Size,
			"url":       "/uploads/" + file.Filename,
		})
	})

	// 多文件上传
	r.POST("/upload/files", func(c *gin.Context) {
		form, err := c.MultipartForm()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "没有找到文件"})
			return
		}

		files := form.File["files"]
		var uploaded []string

		for _, file := range files {
			savePath := filepath.Join("./uploads", file.Filename)
			if err := c.SaveUploadedFile(file, savePath); err != nil {
				continue // 跳过失败的文件, 不影响其他文件的上传
			}
			uploaded = append(uploaded, savePath)
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "文件上传成功",
			"files":    uploaded,
		})
	})

	// 手动处理文件内容
	r.POST("/upload/manual", func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "没有找到文件"})
			return
		}

		// 打开文件
		src, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "打开文件失败"})
			return
		}
		defer src.Close()

		// 手动创建目标文件
		savePath := filepath.Join("./uploads", file.Filename)
		dst, err := os.Create(savePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "创建文件失败"})
			return
		}
		defer dst.Close()

		// 手动复制
		buf := make([]byte, 1024*1024)
		for {
			n, _ := src.Read(buf)
			if n == 0 {
				break
			}
			dst.Write(buf[:n])
		}

		c.JSON(http.StatusOK, gin.H{"msg": "文件上传成功"})
	})

	r.Run(":8080")
}