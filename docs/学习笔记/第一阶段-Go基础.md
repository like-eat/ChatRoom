# Gin 学习第一阶段：Go 基础

> 这些是学 Gin 之前必须掌握的 Go 基础。每个示例都可以单独运行。

---

## 1. net/http 标准库 — 理解 Gin 封装了什么

```go
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// 用户结构体
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
}

// 模拟数据库
var users = map[int]*User{
	1: {ID: 1, Name: "张三", Age: 25},
	2: {ID: 2, Name: "李四", Age: 30},
}

func main() {
	// 注册路由 —— Gin 的 GET/POST 就是对这个的封装
	http.HandleFunc("/user", handleUser)        // 处理 /user 和 /user?id=1
	http.HandleFunc("/user/create", handleCreate) // POST 创建用户

	fmt.Println("服务启动在 http://127.0.0.1:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// 一个 Handler 处理 GET 和 POST 两种请求
func handleUser(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// 解析查询参数 /user?id=1 —— Gin 的 c.Query("id")
		idStr := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "id 无效", http.StatusBadRequest)
			return
		}
		user, ok := users[id]
		if !ok {
			http.Error(w, "用户不存在", http.StatusNotFound)
			return
		}
		// 返回 JSON —— Gin 的 c.JSON()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)

	case http.MethodPost:
		http.Error(w, "不支持的方法", http.StatusMethodNotAllowed)
	}
}

func handleCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "只支持 POST", http.StatusMethodNotAllowed)
		return
	}

	// 读取请求体 —— Gin 的 ShouldBindJSON()
	body, _ := io.ReadAll(r.Body)
	var user User
	if err := json.Unmarshal(body, &user); err != nil {
		http.Error(w, "JSON 格式错误", http.StatusBadRequest)
		return
	}

	user.ID = len(users) + 1
	users[user.ID] = &user

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// 测试：
// curl http://127.0.0.1:8080/user?id=1
// curl -X POST -d '{"name":"王五","age":28}' http://127.0.0.1:8080/user/create
```

---

## 2. struct + json tag — Gin 绑定全靠它

```go
package main

import (
	"encoding/json"
	"fmt"
)

// json tag 决定了字段在 JSON 里的名字
// Gin 的 ShouldBindJSON 就是靠 json tag 来匹配字段的
type LoginRequest struct {
	Telephone string `json:"telephone"`  // JSON 中的 "telephone" → 结构体字段 Telephone
	Password  string `json:"password"`
}

// 没有 json tag，JSON 字段名会和结构体字段名做模糊匹配
type BadRequest struct {
	Telephone string // JSON 中的 "telephone" 无法匹配 Telephone → 零值 ""
}

// 嵌套结构体 + 指针类型（可选字段）
type UpdateUserRequest struct {
	Nickname string  `json:"nickname"`          // 必填
	Avatar   *string `json:"avatar,omitempty"`  // 可选字段，用指针 + omitempty
	Age      int     `json:"age" binding:"required,gte=0,lte=150"` // Gin 的 binding tag
}

func main() {
	// ===== JSON → struct（反序列化，这是 Gin 绑定做的事） =====
	jsonStr := `{"telephone":"13800138000","password":"123456"}`
	var req LoginRequest
	json.Unmarshal([]byte(jsonStr), &req)
	fmt.Printf("反序列化: %+v\n", req)

	// ===== struct → JSON（序列化，这是 Gin c.JSON() 做的事） =====
	resp := LoginRequest{Telephone: "13800138000", Password: "***"}
	respBytes, _ := json.Marshal(resp)
	fmt.Printf("序列化: %s\n", string(respBytes))

	// ===== omitempty 演示 =====
	avatar := "https://example.com/avatar.png"
	update := UpdateUserRequest{
		Nickname: "新昵称",
		Avatar:   &avatar, // 传了值就序列化
	}
	withAvatar, _ := json.Marshal(update)
	fmt.Printf("有 avatar: %s\n", string(withAvatar))

	updateNoAvatar := UpdateUserRequest{Nickname: "新昵称"} // Avatar 是 nil
	noAvatar, _ := json.Marshal(updateNoAvatar)
	fmt.Printf("无 avatar: %s\n", string(noAvatar)) // avatar 字段不会出现
}

// 输出：
// 反序列化: {Telephone:13800138000 Password:123456}
// 序列化: {"telephone":"13800138000","password":"***"}
// 有 avatar: {"nickname":"新昵称","avatar":"https://example.com/avatar.png","age":0}
// 无 avatar: {"nickname":"新昵称","age":0}
```

---

## 3. interface{} 和 context.Context

```go
package main

import (
	"context"
	"fmt"
	"time"
)

// ========================================
// interface{}（空接口）= 任意类型
// 在 Gin 中：c.Set("key", value) → value 是 interface{}
//           c.Get("key") → 返回 interface{}，需要类型断言
// ========================================

// Gin 的 c.Set / c.Get 就是这样的
type Store struct {
	data map[string]interface{} // 可以存任意类型的值
}

func (s *Store) Set(key string, value interface{}) {
	s.data[key] = value
}

func (s *Store) Get(key string) interface{} {
	return s.data[key]
}

// ========================================
// context.Context：传递超时、取消信号、请求级别的值
// Gin 的 c.Request.Context() 就返回这个
// ========================================

// 模拟一个耗时的数据库查询
func queryDatabase(ctx context.Context, query string) (string, error) {
	// 用 select 同时监听两个 channel
	select {
	case <-time.After(2 * time.Second): // 模拟 2 秒查询
		return fmt.Sprintf("查询结果: %s -> 100 条记录", query), nil
	case <-ctx.Done(): // 超时或取消
		return "", ctx.Err()
	}
}

func main() {
	// ===== interface{} 演示 =====
	store := &Store{data: make(map[string]interface{})}
	store.Set("user", "张三")        // 存 string
	store.Set("count", 42)          // 存 int
	store.Set("active", true)       // 存 bool

	// 类型断言：从 interface{} 取出值
	user := store.Get("user").(string) // "张三"
	count := store.Get("count").(int)  // 42
	fmt.Println(user, count)

	// 安全类型断言（推荐）
	if v, ok := store.Get("active").(bool); ok {
		fmt.Println("active =", v)
	}

	// ===== context.Context 演示 =====
	// 设置 1 秒超时
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel() // 用完记得取消

	result, err := queryDatabase(ctx, "SELECT * FROM users")
	if err != nil {
		fmt.Println("查询失败:", err) // 1 秒后会超时
	} else {
		fmt.Println(result)
	}
}

// 输出：
// 张三 42
// active = true
// 查询失败: context deadline exceeded
```

---

## 4. goroutine + channel — Gin 的并发基石

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

// ========================================
// Gin 的每个 HTTP 请求就是一个 goroutine
// 但多个 goroutine 共享的数据需要保护
// ========================================

// ===== 示例 1：简单 goroutine =====
func example1() {
	// go 关键字启动一个 goroutine
	go func() {
		fmt.Println("我在另一个 goroutine 里运行")
	}()

	time.Sleep(100 * time.Millisecond) // 等等 goroutine 执行完
}

// ===== 示例 2：无缓冲 channel（同步）=====
func example2() {
	ch := make(chan string) // 无缓冲：发送方必须等接收方

	go func() {
		time.Sleep(1 * time.Second)
		ch <- "任务完成" // 发送
		fmt.Println("发送成功")
	}()

	fmt.Println("等待接收...")
	msg := <-ch // 接收（阻塞，直到有数据）
	fmt.Println("收到:", msg)
}

// ===== 示例 3：sync.Mutex 保护共享数据 =====
// 这就是本项目中 Server.mutex 的作用
func example3() {
	type Counter struct {
		mu    sync.Mutex
		count int
	}

	counter := &Counter{}
	var wg sync.WaitGroup

	// 启动 1000 个 goroutine 同时 +1
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.mu.Lock()   // 加锁
			counter.count++     // 安全地修改
			counter.mu.Unlock() // 解锁
		}()
	}

	wg.Wait()
	fmt.Println("count:", counter.count) // 一定是 1000
}

// ===== 示例 4：有缓冲 channel（本项目 Transmit 的做法）=====
func example4() {
	// 缓冲大小 3：最多存 3 个，满了就阻塞发送方
	ch := make(chan int, 3)

	ch <- 1 // 不阻塞
	ch <- 2 // 不阻塞
	ch <- 3 // 不阻塞
	// ch <- 4 // 阻塞！缓冲区满了

	fmt.Println(<-ch) // 1
	fmt.Println(<-ch) // 2
	fmt.Println(<-ch) // 3
}

// ===== 示例 5：select 多路复用（本项目的 Login/Logout/Transmit）=====
func example5() {
	login := make(chan string, 5)
	logout := make(chan string, 5)
	transmit := make(chan string, 5)

	// 模拟事件
	go func() { login <- "用户A登录" }()
	go func() { transmit <- "你好" }()
	go func() { logout <- "用户A退出" }()

	// 一个循环处理三种事件
	for i := 0; i < 3; i++ {
		select {
		case msg := <-login:
			fmt.Println("【登录事件】", msg)
		case msg := <-logout:
			fmt.Println("【登出事件】", msg)
		case msg := <-transmit:
			fmt.Println("【消息转发】", msg)
		}
	}
}

func main() {
	fmt.Println("=== 示例 1：goroutine ===")
	example1()

	fmt.Println("\n=== 示例 2：无缓冲 channel ===")
	example2()

	fmt.Println("\n=== 示例 3：Mutex ===")
	example3()

	fmt.Println("\n=== 示例 4：有缓冲 channel ===")
	example4()

	fmt.Println("\n=== 示例 5：select 多路复用 ===")
	example5()
}
```

---

## 配套练习

按顺序完成，每个都能独立跑通后再进入下一阶段：

1. 用 `net/http` 写一个 `/hello?name=xxx` 接口，返回 `{"greeting": "hello, xxx"}`
2. 写一个 `LoginRequest` 结构体，用 `json.Marshal` 和 `json.Unmarshal` 各写一个例子
3. 创建一个带 `sync.Mutex` 的 map，启动 10 个 goroutine 同时读写，观察结果
4. 写一个 select 同时监听 3 个 channel，模拟本项目消息服务器的三个事件流
