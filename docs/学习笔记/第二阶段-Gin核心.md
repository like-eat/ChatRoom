# Gin 学习第二阶段：Gin 核心

> 七个核心概念，每个都有可运行的完整代码。

---

## 准备工作

```bash
# 在新建的目录中初始化
mkdir gin-learning && cd gin-learning
go mod init gin-learning
go get github.com/gin-gonic/gin
```

---

## 1. 路由与参数

```go
package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// ===== 不同 HTTP 方法 =====
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"msg": "pong"})
	})

	r.POST("/submit", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"msg": "post ok"})
	})

	// ===== 路径参数 /user/:id =====
	r.GET("/user/:id", func(c *gin.Context) {
		id := c.Param("id") // 获取路径参数
		c.JSON(http.StatusOK, gin.H{"user_id": id})
	})

	// ===== 查询参数 /search?q=xxx&page=1 =====
	r.GET("/search", func(c *gin.Context) {
		q := c.Query("q")             // 获取查询参数，没有则返回 ""
		page := c.DefaultQuery("page", "1") // 没有则用默认值
		c.JSON(http.StatusOK, gin.H{
			"query": q,
			"page":  page,
		})
	})

	r.Run(":8080")
}

// 测试：
// curl http://127.0.0.1:8080/ping
// curl -X POST http://127.0.0.1:8080/submit
// curl http://127.0.0.1:8080/user/123
// curl "http://127.0.0.1:8080/search?q=golang&page=2"
```

---

## 2. 请求绑定与验证

```go
package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// binding tag 用于验证
type LoginRequest struct {
	Telephone string `json:"telephone" binding:"required,len=11"` // 必填，11位
	Password  string `json:"password" binding:"required,min=6"`    // 必填，最短6位
	Remember  bool   `json:"remember"`                             // 可选，默认 false
}

// form tag 用于表单绑定
type RegisterRequest struct {
	Nickname  string `form:"nickname" binding:"required"`
	Telephone string `form:"telephone" binding:"required,len=11"`
	Password  string `form:"password" binding:"required,min=6"`
}

func main() {
	r := gin.Default()

	// ===== JSON 绑定 + 验证 =====
	r.POST("/login", func(c *gin.Context) {
		var req LoginRequest

		// ShouldBindJSON：绑定 JSON 并验证
		if err := c.ShouldBindJSON(&req); err != nil {
			// err.Error() 包含具体的验证失败原因
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 绑定成功后，直接使用结构体字段
		c.JSON(http.StatusOK, gin.H{
			"msg":       "登录成功",
			"telephone": req.Telephone,
		})
	})

	// ===== 表单绑定 =====
	r.POST("/register", func(c *gin.Context) {
		var req RegisterRequest

		if err := c.ShouldBind(&req); err != nil { // 自动检测 Content-Type
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"msg":       "注册成功",
			"nickname":  req.Nickname,
			"telephone": req.Telephone,
		})
	})

	r.Run(":8080")
}

// 测试：
// curl -X POST -H "Content-Type: application/json" -d '{"telephone":"13800138000","password":"123456"}' http://127.0.0.1:8080/login
// curl -X POST -d "nickname=张三&telephone=13800138000&password=123456" http://127.0.0.1:8080/register
// curl -X POST -H "Content-Type: application/json" -d '{"telephone":"123"}' http://127.0.0.1:8080/login  (验证失败)
```

---

## 3. 响应处理

```go
package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 统一返回格式（本项目就是这么做的）
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// 封装三个快捷方法
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{Code: 200, Message: "success", Data: data})
}

func BadRequest(c *gin.Context, msg string) {
	c.JSON(http.StatusBadRequest, Response{Code: 400, Message: msg})
}

func ServerError(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, Response{Code: 500, Message: "服务器内部错误"})
}

func main() {
	r := gin.Default()

	r.GET("/user/:id", func(c *gin.Context) {
		id := c.Param("id")
		if id == "0" {
			BadRequest(c, "用户ID不能为0")
			return
		}
		user := gin.H{"id": id, "name": "张三"}
		Success(c, user)
	})

	r.GET("/panic-test", func(c *gin.Context) {
		ServerError(c)
	})

	// Gin 的多种响应方式
	r.GET("/formats", func(c *gin.Context) {
		format := c.DefaultQuery("format", "json")
		switch format {
		case "json":
			c.JSON(http.StatusOK, gin.H{"msg": "JSON 格式"})
		case "xml":
			c.XML(http.StatusOK, gin.H{"msg": "XML 格式"})
		case "string":
			c.String(http.StatusOK, "纯文本格式")
		case "file":
			c.File("./main.go") // 返回文件下载
		default:
			c.JSON(http.StatusOK, gin.H{"msg": "默认 JSON"})
		}
	})

	r.Run(":8080")
}

// 测试：
// curl http://127.0.0.1:8080/user/123  → {"code":200,"message":"success","data":{"id":"123","name":"张三"}}
// curl http://127.0.0.1:8080/user/0    → {"code":400,"message":"用户ID不能为0"}
// curl "http://127.0.0.1:8080/formats?format=xml"
```

---

## 4. 中间件

```go
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// ===== 自定义中间件 1：请求日志 =====
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Next() 之前 = 请求处理前
		fmt.Printf("[请求开始] %s %s\n", c.Request.Method, c.Request.URL.Path)

		c.Next() // 执行下一个中间件 / 实际 handler

		// Next() 之后 = 请求处理后
		cost := time.Since(start)
		fmt.Printf("[请求结束] %s %s 耗时: %v 状态码: %d\n",
			c.Request.Method, c.Request.URL.Path, cost, c.Writer.Status())
	}
}

// ===== 自定义中间件 2：简易鉴权 =====
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token != "Bearer secret-token" {
			// Abort()：停止，不执行后续 handler
			c.JSON(401, gin.H{"error": "未授权"})
			c.Abort()
			return
		}
		// 把用户信息存入上下文，后续 handler 可以取
		c.Set("user_id", "U001")
		c.Next()
	}
}

// ===== 自定义中间件 3：计算请求耗时 =====
func CostMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		// c.Get 取出前面中间件存入的值
		userID, _ := c.Get("user_id")
		log.Printf("用户 %v 的请求处理完毕", userID)
	}
}

func main() {
	r := gin.Default()

	// ===== 全局中间件（所有路由都生效）=====
	r.Use(LoggerMiddleware())

	// ===== 公开路由（不需要鉴权）=====
	r.GET("/public", func(c *gin.Context) {
		c.JSON(200, gin.H{"msg": "这是公开接口"})
	})

	// ===== 鉴权路由组 =====
	auth := r.Group("/api")
	auth.Use(AuthMiddleware(), CostMiddleware())
	{
		auth.GET("/profile", func(c *gin.Context) {
			userID, _ := c.Get("user_id")
			c.JSON(200, gin.H{"user_id": userID, "name": "张三"})
		})

		auth.GET("/orders", func(c *gin.Context) {
			c.JSON(200, gin.H{"orders": []string{"订单1", "订单2"}})
		})
	}

	r.Run(":8080")
}

// 测试：
// curl http://127.0.0.1:8080/public                    → 成功
// curl http://127.0.0.1:8080/api/profile               → 401
// curl -H "Authorization: Bearer secret-token" http://127.0.0.1:8080/api/profile  → 成功
```

---

## 5. 路由分组

```go
package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 模拟的鉴权中间件
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("Authorization") == "" {
			c.JSON(401, gin.H{"error": "需要登录"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// 模拟的管理员中间件
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetHeader("X-Role")
		if role != "admin" {
			c.JSON(403, gin.H{"error": "需要管理员权限"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func main() {
	r := gin.Default()

	// ===== v1 版本 =====
	v1 := r.Group("/api/v1")
	{
		// 公开接口
		v1.POST("/login", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"msg": "登录成功"})
		})
		v1.POST("/register", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"msg": "注册成功"})
		})

		// 需要登录的接口
		auth := v1.Group("")
		auth.Use(Auth())
		{
			auth.GET("/user/info", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"name": "张三"})
			})
			auth.PUT("/user/update", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"msg": "更新成功"})
			})

			// 管理员专享
			admin := auth.Group("/admin")
			admin.Use(AdminOnly())
			{
				admin.GET("/users", func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{"users": []string{"张三", "李四"}})
				})
				admin.POST("/user/disable", func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{"msg": "用户已禁用"})
				})
			}
		}
	}

	// ===== v2 版本 =====
	v2 := r.Group("/api/v2")
	{
		v2.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})
	}

	r.Run(":8080")
}

// 整个路由结构：
// /api/v1/login             公开
// /api/v1/register          公开
// /api/v1/user/info         需登录
// /api/v1/user/update       需登录
// /api/v1/admin/users       需登录 + 管理员
// /api/v1/admin/user/disable  需登录 + 管理员
// /api/v2/health            公开
```

---

## 6. 文件上传

```go
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

	// 限制上传文件最大 10MB
	r.MaxMultipartMemory = 10 << 20 // 10 MB

	// ===== 单文件上传 =====
	r.POST("/upload/avatar", func(c *gin.Context) {
		// 获取上传文件
		file, err := c.FormFile("file") // "file" 是表单字段名
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "没有找到文件"})
			return
		}

		// 文件信息
		fmt.Printf("文件名: %s, 大小: %d bytes\n", file.Filename, file.Size)

		// 保存到磁盘（Gin 的方法：SaveUploadedFile）
		savePath := filepath.Join("./uploads", file.Filename)
		if err := c.SaveUploadedFile(file, savePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件失败"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"msg":      "上传成功",
			"filename": file.Filename,
			"size":     file.Size,
			"url":      "/uploads/" + file.Filename,
		})
	})

	// ===== 多文件上传 =====
	r.POST("/upload/files", func(c *gin.Context) {
		form, err := c.MultipartForm()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "不是 multipart 表单"})
			return
		}

		files := form.File["files"] // "files" 是表单字段名
		var uploaded []string

		for _, file := range files {
			savePath := filepath.Join("./uploads", file.Filename)
			if err := c.SaveUploadedFile(file, savePath); err != nil {
				continue
			}
			uploaded = append(uploaded, file.Filename)
		}

		c.JSON(http.StatusOK, gin.H{
			"msg":   "上传成功",
			"files": uploaded,
		})
	})

	// ===== 手动处理文件内容（不使用 SaveUploadedFile）=====
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

		// 手动复制（可以在这里做校验、压缩等处理）
		buf := make([]byte, 1024*1024) // 1MB 缓冲区
		for {
			n, _ := src.Read(buf)
			if n == 0 {
				break
			}
			dst.Write(buf[:n])
		}

		c.JSON(http.StatusOK, gin.H{"msg": "上传成功"})
	})

	r.Run(":8080")
}

// 测试：
// curl -F "file=@/path/to/avatar.png" http://127.0.0.1:8080/upload/avatar
// curl -F "files=@a.png" -F "files=@b.png" http://127.0.0.1:8080/upload/files
```

---

## 7. 静态资源

```go
package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// ===== 方式 1：映射单个目录 =====
	// 访问 /static/avatar.png → 实际读取 ./public/avatar.png
	r.Static("/static", "./public")

	// ===== 方式 2：映射到文件系统（更灵活）=====
	r.StaticFS("/assets", gin.Dir("./public", true)) // true = 显示目录列表

	// ===== 方式 3：单个文件 =====
	r.StaticFile("/favicon.ico", "./public/favicon.ico")

	// ===== 方式 4：自定义文件服务器 =====
	r.GET("/download/:filename", func(c *gin.Context) {
		filename := c.Param("filename")
		// File 会设置正确的 Content-Type
		c.File("./public/" + filename)
	})

	// ===== 完整示例：头像上传 + 静态访问 =====
	r.POST("/upload/avatar", func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(400, gin.H{"error": "没有文件"})
			return
		}
		c.SaveUploadedFile(file, "./public/avatars/"+file.Filename)
		c.JSON(200, gin.H{
			"url": "/static/avatars/" + file.Filename,
		})
	})

	r.Run(":8080")
}

// 使用前创建目录结构：
// mkdir -p public/avatars
// echo "hello" > public/test.txt

// 测试：
// curl http://127.0.0.1:8080/static/test.txt
// curl http://127.0.0.1:8080/assets/test.txt
```

---

## 七个概念在项目中的对应

| 概念 | 本项目代码位置 |
|---|---|
| 路由 | [https_server.go](internal/https_server/https_server.go) — `GE.GET` / `GE.POST` |
| 绑定 | [api/v1/user_info_controller.go](api/v1/user_info_controller.go) — `c.BindJSON` |
| 响应 | [api/v1/controller.go](api/v1/controller.go) — `JsonBack()` 统一返回 |
| 中间件 | [middleware/auth.go](internal/middleware/auth.go) — JWT 鉴权 |
| 路由分组 | [https_server.go](internal/https_server/https_server.go) — 公开/登录/管理员三组 |
| 文件上传 | [api/v1/message_controller.go](api/v1/message_controller.go) — `UploadFile` |
| 静态资源 | [https_server.go](internal/https_server/https_server.go) — `Static` / `StaticFS` |
