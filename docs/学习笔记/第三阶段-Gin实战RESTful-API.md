# Gin 学习第三阶段：完整 RESTful API

> 一个完整的用户管理系统：注册、登录、JWT 鉴权、CRUD、分页查询。
> 新建一个文件夹，把所有代码写进去，`go run main.go` 就能跑。

---

## 项目结构

```text
gin-demo/
├── main.go          # 入口 + 路由
├── handler.go       # 请求处理（Controller 层）
├── service.go       # 业务逻辑（Service 层）
├── model.go         # 数据模型
├── store.go         # 模拟数据库（内存存储）
├── middleware.go     # JWT 鉴权中间件
├── jwt.go           # JWT 签发与验证
└── response.go      # 统一返回格式
```

---

## 1. response.go — 统一返回格式

```go
package main

import "net/http"
import "github.com/gin-gonic/gin"

// 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{Code: 200, Message: "success", Data: data})
}

func Error(c *gin.Context, code int, msg string) {
	c.JSON(code, Response{Code: code, Message: msg})
}
```

---

## 2. model.go — 数据模型

```go
package main

import "time"

// 用户模型
type User struct {
	ID        int64     `json:"id"`
	Nickname  string    `json:"nickname"`
	Telephone string    `json:"telephone"`
	Password  string    `json:"-"` // json:"-" 表示序列化时隐藏
	Avatar    string    `json:"avatar"`
	Gender    int8      `json:"gender"` // 0.男 1.女
	CreatedAt time.Time `json:"created_at"`
}

// ===== 请求 DTO =====

type RegisterReq struct {
	Nickname  string `json:"nickname" binding:"required,min=2,max=20"`
	Telephone string `json:"telephone" binding:"required,len=11"`
	Password  string `json:"password" binding:"required,min=6,max=20"`
	Gender    int8   `json:"gender" binding:"oneof=0 1"`
}

type LoginReq struct {
	Telephone string `json:"telephone" binding:"required,len=11"`
	Password  string `json:"password" binding:"required"`
}

type UpdateUserReq struct {
	Nickname string `json:"nickname" binding:"omitempty,min=2,max=20"`
	Avatar   string `json:"avatar"`
	Gender   *int8  `json:"gender" binding:"omitempty,oneof=0 1"` // 指针，区分 0 和未传
}

// ===== 响应 DTO =====

type LoginResp struct {
	Token      string `json:"token"`
	UserID     int64  `json:"user_id"`
	Nickname   string `json:"nickname"`
	Telephone  string `json:"telephone"`
	Avatar     string `json:"avatar"`
}

type UserListResp struct {
	Total int64   `json:"total"`
	List  []*User `json:"list"`
}
```

---

## 3. store.go — 模拟数据库（内存存储 + 自增 ID）

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

// Store 用内存 map 模拟数据库，用 Mutex 保护并发访问
type Store struct {
	mu        sync.RWMutex
	users     map[int64]*User     // ID → User
	byPhone   map[string]int64    // 手机号 → ID
	autoID    int64               // 自增 ID
}

func NewStore() *Store {
	return &Store{
		users:   make(map[int64]*User),
		byPhone: make(map[string]int64),
		autoID:  0,
	}
}

// 创建用户
func (s *Store) CreateUser(req RegisterReq) (*User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 检查手机号是否已注册
	if _, exists := s.byPhone[req.Telephone]; exists {
		return nil, fmt.Errorf("手机号已注册")
	}

	s.autoID++
	user := &User{
		ID:        s.autoID,
		Nickname:  req.Nickname,
		Telephone: req.Telephone,
		Password:  req.Password, // 生产环境要 bcrypt
		Avatar:    "https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png",
		Gender:    req.Gender,
		CreatedAt: time.Now(),
	}

	s.users[user.ID] = user
	s.byPhone[user.Telephone] = user.ID
	return user, nil
}

// 根据手机号查找用户
func (s *Store) FindByPhone(telephone string) (*User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	id, exists := s.byPhone[telephone]
	if !exists {
		return nil, fmt.Errorf("用户不存在")
	}
	user := *s.users[id] // 返回副本
	return &user, nil
}

// 根据 ID 查找用户
func (s *Store) FindByID(id int64) (*User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.users[id]
	if !exists {
		return nil, fmt.Errorf("用户不存在")
	}
	result := *user // 返回副本
	return &result, nil
}

// 更新用户
func (s *Store) UpdateUser(id int64, req UpdateUserReq) (*User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	user, exists := s.users[id]
	if !exists {
		return nil, fmt.Errorf("用户不存在")
	}

	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}
	if req.Gender != nil {
		user.Gender = *req.Gender
	}

	return user, nil
}

// 分页查询用户列表
func (s *Store) ListUsers(page, pageSize int) ([]*User, int64) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	total := int64(len(s.users))
	start := (page - 1) * pageSize
	if start >= int(total) {
		return []*User{}, total
	}

	var result []*User
	i := 0
	for _, user := range s.users {
		if i >= start && len(result) < pageSize {
			u := *user // 副本
			result = append(result, &u)
		}
		i++
	}
	return result, total
}

// 删除用户
func (s *Store) DeleteUser(id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	user, exists := s.users[id]
	if !exists {
		return fmt.Errorf("用户不存在")
	}

	delete(s.users, id)
	delete(s.byPhone, user.Telephone)
	return nil
}
```

---

## 4. jwt.go — JWT 签发与验证

```go
package main

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// 生产环境用环境变量，这里写死做演示
var jwtSecret = []byte("kama-chat-jwt-secret-key-2026")

type Claims struct {
	UserID   int64  `json:"user_id"`
	Nickname string `json:"nickname"`
	jwt.RegisteredClaims
}

// 生成 JWT
func GenerateToken(userID int64, nickname string) (string, error) {
	claims := Claims{
		UserID:   userID,
		Nickname: nickname,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "gin-demo",
			Subject:   fmt.Sprintf("%d", userID),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// 验证 JWT
func ParseToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("不支持的签名算法: %v", t.Header["alg"])
		}
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("无效的 token")
	}
	return claims, nil
}
```

---

## 5. middleware.go — JWT 鉴权中间件

```go
package main

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// JWT 鉴权中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Authorization header 获取 token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			Error(c, 401, "缺少 Authorization 头")
			c.Abort()
			return
		}

		// Bearer <token>
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			Error(c, 401, "Authorization 格式错误")
			c.Abort()
			return
		}

		// 验证 token
		claims, err := ParseToken(parts[1])
		if err != nil {
			Error(c, 401, "Token 无效或已过期")
			c.Abort()
			return
		}

		// 把用户信息存入上下文，后续 handler 可用
		c.Set("user_id", claims.UserID)
		c.Set("nickname", claims.Nickname)
		c.Next()
	}
}
```

---

## 6. handler.go — 请求处理层

```go
package main

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// 全局 store 实例
var store = NewStore()

// ===== 注册 =====
func Register(c *gin.Context) {
	var req RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, 400, "参数错误: "+err.Error())
		return
	}

	user, err := store.CreateUser(req)
	if err != nil {
		Error(c, 400, err.Error())
		return
	}

	// 注册成功，直接返回 token
	token, _ := GenerateToken(user.ID, user.Nickname)

	Success(c, LoginResp{
		Token:     token,
		UserID:    user.ID,
		Nickname:  user.Nickname,
		Telephone: user.Telephone,
		Avatar:    user.Avatar,
	})
}

// ===== 登录 =====
func Login(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, 400, "参数错误: "+err.Error())
		return
	}

	user, err := store.FindByPhone(req.Telephone)
	if err != nil {
		Error(c, 400, "手机号或密码错误")
		return
	}

	// 生产环境用 bcrypt.CompareHashAndPassword
	if user.Password != req.Password {
		Error(c, 400, "手机号或密码错误")
		return
	}

	token, _ := GenerateToken(user.ID, user.Nickname)

	Success(c, LoginResp{
		Token:     token,
		UserID:    user.ID,
		Nickname:  user.Nickname,
		Telephone: user.Telephone,
		Avatar:    user.Avatar,
	})
}

// ===== 获取当前用户信息 =====
func GetProfile(c *gin.Context) {
	userID := c.GetInt64("user_id") // 从中间件注入的上下文取值

	user, err := store.FindByID(userID)
	if err != nil {
		Error(c, 404, "用户不存在")
		return
	}

	Success(c, user)
}

// ===== 更新用户信息 =====
func UpdateProfile(c *gin.Context) {
	userID := c.GetInt64("user_id")

	var req UpdateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, 400, "参数错误: "+err.Error())
		return
	}

	user, err := store.UpdateUser(userID, req)
	if err != nil {
		Error(c, 400, err.Error())
		return
	}

	Success(c, user)
}

// ===== 获取用户列表（分页）=====
func ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	users, total := store.ListUsers(page, pageSize)

	Success(c, UserListResp{Total: total, List: users})
}

// ===== 获取单个用户 =====
func GetUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		Error(c, 400, "用户ID格式错误")
		return
	}

	user, err := store.FindByID(id)
	if err != nil {
		Error(c, 404, "用户不存在")
		return
	}

	Success(c, user)
}

// ===== 删除用户（自己或管理员）=====
func DeleteUser(c *gin.Context) {
	userID := c.GetInt64("user_id")
	idStr := c.Param("id")
	targetID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		Error(c, 400, "用户ID格式错误")
		return
	}

	// 只能删自己（简化版，实际应有管理员权限）
	if userID != targetID {
		Error(c, 403, "只能删除自己的账号")
		return
	}

	if err := store.DeleteUser(targetID); err != nil {
		Error(c, 400, err.Error())
		return
	}

	Success(c, gin.H{"msg": "删除成功"})
}
```

---

## 7. main.go — 入口 + 路由注册

```go
package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// ===== 初始化数据：添加两条演示用户 =====
	store.CreateUser(RegisterReq{
		Nickname: "张三", Telephone: "13800001111", Password: "123456", Gender: 0,
	})
	store.CreateUser(RegisterReq{
		Nickname: "李四", Telephone: "13800002222", Password: "123456", Gender: 1,
	})

	// ===== 公开路由 =====
	r.POST("/register", Register)
	r.POST("/login", Login)

	// ===== 需要登录的路由 =====
	auth := r.Group("/api").Use(JWTAuth())
	{
		auth.GET("/profile", GetProfile)        // 获取自己的信息
		auth.PUT("/profile", UpdateProfile)      // 更新自己的信息
		auth.DELETE("/user/:id", DeleteUser)     // 删除账号
	}

	// ===== 公开的查询路由（实际项目中可能需要限制）=====
	r.GET("/users", ListUsers)                  // 用户列表（分页）
	r.GET("/user/:id", GetUserByID)             // 查看某个用户信息

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
```

---

## 初始化项目

```bash
mkdir gin-demo && cd gin-demo
go mod init gin-demo
go get github.com/gin-gonic/gin
go get github.com/golang-jwt/jwt/v5

# 把上面 7 个 .go 文件全部创建在当前目录下

go run .
```

---

## 测试接口

```bash
# 1. 注册
curl -X POST -H "Content-Type: application/json" \
  -d '{"nickname":"王五","telephone":"13800003333","password":"123456","gender":0}' \
  http://127.0.0.1:8080/register

# 2. 登录（把返回的 token 记下来）
curl -X POST -H "Content-Type: application/json" \
  -d '{"telephone":"13800003333","password":"123456"}' \
  http://127.0.0.1:8080/login

# 3. 获取个人信息（替换 YOUR_TOKEN）
curl -H "Authorization: Bearer YOUR_TOKEN" \
  http://127.0.0.1:8080/api/profile

# 4. 更新个人信息
curl -X PUT -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{"nickname":"王五五","gender":1}' \
  http://127.0.0.1:8080/api/profile

# 5. 用户列表（分页）
curl "http://127.0.0.1:8080/users?page=1&page_size=5"

# 6. 查看某个用户
curl http://127.0.0.1:8080/user/1

# 7. 删除用户（只能删自己）
curl -X DELETE -H "Authorization: Bearer YOUR_TOKEN" \
  http://127.0.0.1:8080/api/user/YOUR_USER_ID
```

---

## 完成后对照本项目

做完这个小项目后，再看本项目的代码，你会发现结构完全对应：

| 小项目 | 本项目 |
|---|---|
| `handler.go` | `api/v1/*_controller.go` |
| `service.go` / `store.go` | `internal/service/gorm/*_service.go` |
| `model.go` | `internal/model/*.go` |
| `middleware.go` | `internal/middleware/auth.go` |
| `jwt.go` | `internal/service/auth/jwt.go` |
| `response.go` | `api/v1/controller.go` 的 `JsonBack` |
| `main.go` 路由注册 | `internal/https_server/https_server.go` |

区别只是：本项目的 store 是真正 MySQL + Redis，你的练习项目是内存 map，但**分层思想完全一样**。
