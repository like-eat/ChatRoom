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

## 3. 指针、defer、错误处理 — Go 语言三大基本机制

```go
package main

import (
	"errors"
	"fmt"
	"os"
)

// ========================================
// 一、指针 * 和 &
// ========================================
// &变量    = 获取变量的内存地址（取地址）
// *类型    = 声明这是一个指针类型
// *指针变量 = 解引用，获取指针指向的值
//
// 在 Gin 中你每次都会看到：c.ShouldBindJSON(&req)
// "&req" 就是把 req 的地址传给函数，让函数可以直接修改 req

// 值传递 vs 指针传递
func updateByValue(u User) {
	u.Name = "被修改了" // 改的是副本，不影响原值
}

func updateByPointer(u *User) {
	u.Name = "被修改了" // 改的是原值！
}

type User struct {
	ID   int
	Name string
}

// ========================================
// 二、defer：延迟执行（Go 的 finally）
// ========================================
// defer 语句在函数 return 之前执行，用于资源清理
// 多个 defer 按 LIFO（后进先出）顺序执行

func deferDemo() {
	// 模拟打开文件
	fmt.Println("打开文件")

	// defer 确保文件一定会被关闭，即使后面发生了 panic
	defer fmt.Println("关闭文件（defer 1）")
	defer fmt.Println("刷新缓冲区（defer 2）")

	fmt.Println("读写文件...")
	fmt.Println("处理完成")
	// 函数结束时自动执行：先 defer 2，再 defer 1
}

// ========================================
// 三、错误处理：Go 没有 try-catch
// ========================================
// Go 的错误处理模式：函数返回 (result, error)，调用方检查 error

// 模拟一个可能失败的操作
func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("除数不能为零")
	}
	return a / b, nil
}

// 模拟多步骤操作：错误逐层传递
func readConfig(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		// 用 fmt.Errorf 包装错误，添加上下文
		return "", fmt.Errorf("读取配置文件失败: %w", err)
	}
	return string(data), nil
}

// 错误判断：errors.Is 和 errors.As
func checkError(err error) {
	// errors.Is：判断错误链中是否包含某个特定错误
	if errors.Is(err, os.ErrNotExist) {
		fmt.Println("→ 文件不存在")
	}
}

func main() {
	// ===== 指针演示 =====
	user := User{ID: 1, Name: "原始名字"}

	updateByValue(user)
	fmt.Println("值传递后:", user.Name) // "原始名字" — 没变！

	updateByPointer(&user)               // &user 是取地址
	fmt.Println("指针传递后:", user.Name) // "被修改了" — 变了！

	// 这就是为什么 Gin 的绑定用 ShouldBindJSON(&req) 而不是 ShouldBindJSON(req)
	// 只有传指针，Gin 才能把 JSON 数据写入你的结构体

	// ===== defer 演示 =====
	fmt.Println("\n--- defer 演示 ---")
	deferDemo()
	// 输出顺序：
	// 打开文件 → 读写文件... → 处理完成 → 刷新缓冲区（defer 2）→ 关闭文件（defer 1）

	// ===== 错误处理演示 =====
	fmt.Println("\n--- 错误处理演示 ---")

	// 模式 1：直接处理
	result, err := divide(10, 0)
	if err != nil {
		fmt.Println("除零错误:", err)
	} else {
		fmt.Println("结果:", result)
	}

	// 模式 2：错误传递
	_, err = readConfig("不存在的文件.toml")
	if err != nil {
		fmt.Println("读取失败:", err)
		checkError(err)
	}

	// 模式 3：忽略错误（不推荐，仅演示）
	result, _ = divide(10, 2) // _ 丢弃 error（本项目不会这么做）
	fmt.Println("10/2 =", result)
}

// 输出：
// 值传递后: 原始名字
// 指针传递后: 被修改了
//
// --- defer 演示 ---
// 打开文件
// 读写文件...
// 处理完成
// 刷新缓冲区（defer 2）
// 关闭文件（defer 1）
//
// --- 错误处理演示 ---
// 除零错误: 除数不能为零
// 读取失败: 读取配置文件失败: open 不存在的文件.toml: The system cannot find the file specified.
// → 文件不存在
// 10/2 = 5
```

---

## 4. slice 和 map — 最常用的数据结构

```go
package main

import "fmt"

// ========================================
// 一、slice（切片）= 动态数组
// ========================================
// 项目中到处都是 slice：
//   []string              — 群成员列表
//   []model.Message       — 消息列表
//   []respond.GetUserListRespond — API 返回列表

func sliceDemo() {
	// ===== 创建 slice =====
	// 方式 1：字面量
	names := []string{"张三", "李四", "王五"}

	// 方式 2：make（预分配容量）
	ids := make([]int, 0, 10) // len=0, cap=10

	// ===== append：追加元素 =====
	ids = append(ids, 1, 2, 3)
	fmt.Println("ids:", ids) // [1 2 3]

	// ===== len 和 cap =====
	fmt.Printf("len=%d, cap=%d\n", len(ids), cap(ids)) // len=3, cap=10

	// ===== 遍历 =====
	for i, name := range names {
		fmt.Printf("  [%d] %s\n", i, name)
	}

	// ===== 切片操作 [start:end] =====
	sub := names[1:3] // index 1 到 2（不包含 3）
	fmt.Println("切片:", sub) // [李四 王五]

	// ===== 删除元素（中间位置）=====
	// 删除 index=1 的元素（本项目 RemoveGroupMembers 就是这样做的）
	names = append(names[:1], names[2:]...)
	fmt.Println("删除后:", names) // [张三 王五]

	// ===== nil slice vs empty slice =====
	var nilSlice []string       // nil slice — len=0，JSON 序列化为 null
	emptySlice := []string{}    // empty slice — len=0，JSON 序列化为 []
	fmt.Printf("nilSlice: %v, emptySlice: %v\n", nilSlice, emptySlice)
}

// ========================================
// 二、map = 键值对 / 字典
// ========================================
// 项目中的用法：
//   map[string]*Client   — 本项目的 Server.Clients（在线用户）
//   map[string]interface{} — Gin 的 c.Keys（请求上下文存储）
//   gin.H{}              — 其实就是 map[string]interface{}

func mapDemo() {
	// ===== 创建 map =====
	// 方式 1：make
	scores := make(map[string]int)

	// 方式 2：字面量
	users := map[string]string{
		"U001": "张三",
		"U002": "李四",
	}

	// ===== 写入 =====
	scores["张三"] = 95
	scores["李四"] = 87

	// ===== 读取 =====
	fmt.Println("张三的分数:", scores["张三"])

	// ===== 检查 key 是否存在 =====
	// 非常重要！这是 Go 的 "comma ok" 惯用法
	score, ok := scores["王五"]
	if ok {
		fmt.Println("王五的分数:", score)
	} else {
		fmt.Println("王五不存在")
	}

	// ===== 遍历 =====
	for name, score := range scores {
		fmt.Printf("  %s → %d\n", name, score)
	}

	// ===== 删除 =====
	delete(scores, "李四")
	fmt.Println("删除李四后:", scores)

	// ===== map 的零值是 nil，写入 nil map 会 panic！=====
	// var m map[string]int   // nil map
	// m["key"] = 1           // panic! 一定要先 make
	fmt.Println("用户表:", users)
}

// ========================================
// 三、项目中 slice + map 的组合用法
// ========================================

// 模拟本项目的在线用户管理
type Client struct {
	UserID   string
	Nickname string
}

type Hub struct {
	// 在线客户端：map[用户ID]*Client（本项目的 Server.Clients）
	clients map[string]*Client

	// 消息队列：有缓冲 channel（本项目的 Server.Transmit）
	broadcast chan []byte
}

func NewHub() *Hub {
	return &Hub{
		clients:   make(map[string]*Client),
		broadcast: make(chan []byte, 100),
	}
}

func (h *Hub) AddClient(c *Client) {
	h.clients[c.UserID] = c
}

// 获取所有在线用户名（返回 slice）
func (h *Hub) GetOnlineNames() []string {
	names := make([]string, 0, len(h.clients))
	for _, c := range h.clients {
		names = append(names, c.Nickname)
	}
	return names
}

func main() {
	fmt.Println("=== slice 演示 ===")
	sliceDemo()

	fmt.Println("\n=== map 演示 ===")
	mapDemo()

	fmt.Println("\n=== 项目场景演示 ===")
	hub := NewHub()
	hub.AddClient(&Client{UserID: "U001", Nickname: "张三"})
	hub.AddClient(&Client{UserID: "U002", Nickname: "李四"})
	hub.AddClient(&Client{UserID: "U003", Nickname: "王五"})

	fmt.Println("在线用户:", hub.GetOnlineNames())

	// 查找某个用户是否在线
	if client, ok := hub.clients["U002"]; ok {
		fmt.Printf("U002 在线: %s\n", client.Nickname)
	}
}
```

---

## 5. interface{} 和 context.Context

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

## 6. goroutine + channel — Gin 的并发基石

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

## 7. 方法和接口 — Go 的面向对象

```go
package main

import (
	"fmt"
	"strings"
)

// ========================================
// 一、方法 = 带接收者的函数
// ========================================
// func (接收者) 方法名(参数) 返回值 { ... }
// 这是 Go 的 OOP 方式——没有 class，只有 struct + method
//
// 本项目就是这种模式：
//   type userInfoService struct{}
//   func (u *userInfoService) Login(req LoginRequest) (string, *LoginRespond, int) { ... }

type UserService struct {
	db map[string]*User // 模拟数据库连接
}

type User struct {
	UUID     string
	Nickname string
	Email    string
}

// NewUserService 构造函数（Go 没有构造函数，用普通函数模拟）
func NewUserService() *UserService {
	return &UserService{
		db: make(map[string]*User),
	}
}

// ===== 指针接收者 vs 值接收者 =====

// 指针接收者 (*UserService)：可以修改接收者本身
// 本项目 99% 的方法都用指针接收者
func (s *UserService) CreateUser(uuid, nickname string) *User {
	user := &User{UUID: uuid, Nickname: nickname}
	s.db[uuid] = user    // 修改 s.db（需要指针）
	return user
}

func (s *UserService) UpdateNickname(uuid, newName string) error {
	user, ok := s.db[uuid]
	if !ok {
		return fmt.Errorf("用户 %s 不存在", uuid)
	}
	user.Nickname = newName // 修改 user（需要指针）
	return nil
}

// 值接收者 (UserService)：不能修改接收者，只读操作
func (s UserService) FindByUUID(uuid string) (*User, error) {
	user, ok := s.db[uuid]
	if !ok {
		return nil, fmt.Errorf("用户 %s 不存在", uuid)
	}
	return user, nil
}

// 普通方法（没有接收者）= 就是普通函数
func GenerateUUID(prefix string) string {
	return fmt.Sprintf("%s-1234-5678", prefix)
}

// ========================================
// 二、接口：定义行为，不定义数据
// ========================================
// Go 的接口是隐式实现的——不需要 "implements" 关键字
// 只要你的类型实现了接口的所有方法，就自动满足该接口

// 定义一个 "可持久化的" 接口
type UserRepository interface {
	CreateUser(uuid, nickname string) *User
	FindByUUID(uuid string) (*User, error)
	UpdateNickname(uuid, newName string) error
}

// UserService 实现了 UserRepository 的所有方法
// 不需要写 "implements UserRepository"，编译器自动检查
var _ UserRepository = (*UserService)(nil) // 编译时接口检查（本项目常用技巧）

// ===== 多态：用接口类型作为参数 =====
func PrintUserInfo(repo UserRepository, uuid string) {
	user, err := repo.FindByUUID(uuid)
	if err != nil {
		fmt.Println("找不到用户:", err)
		return
	}
	// 转大写做演示
	fmt.Printf("用户信息: %s (%s)\n", strings.ToUpper(user.Nickname), user.UUID)
}

// ===== 接口组合 =====
type Reader interface {
	Read(p []byte) (n int, err error)
}

type Writer interface {
	Write(p []byte) (n int, err error)
}

// 组合 Reader 和 Writer（就像 io.ReadWriter）
type ReadWriter interface {
	Reader
	Writer
}

// ========================================
// 三、项目中常见的接口用法
// ========================================

// Gin 的 HandlerFunc 本质上就是实现了 http.Handler 接口
// type Handler interface {
//     ServeHTTP(ResponseWriter, *Request)
// }

// 模拟：任何有 "Handle" 方法的类型就是一个简易路由器
type Handler interface {
	Handle(path string) string
}

type UserHandler struct{}
func (h UserHandler) Handle(path string) string {
	return "处理用户请求: " + path
}

type AdminHandler struct{}
func (h AdminHandler) Handle(path string) string {
	return "处理管理员请求: " + path
}

// 用接口实现多态
func route(handler Handler, path string) {
	fmt.Println(handler.Handle(path))
}

func main() {
	// ===== 方法演示 =====
	svc := NewUserService()

	svc.CreateUser("U001", "张三")
	svc.CreateUser("U002", "李四")

	svc.UpdateNickname("U001", "张三三")

	user, _ := svc.FindByUUID("U001")
	fmt.Printf("用户: %+v\n", user)

	fmt.Println("生成 UUID:", GenerateUUID("USER"))

	// ===== 接口演示 =====
	fmt.Println("\n--- 接口多态 ---")
	PrintUserInfo(svc, "U001")

	// 不同的 Handler 类型，同样的 Handle 方法
	route(UserHandler{}, "/login")
	route(AdminHandler{}, "/dashboard")

	// ===== 空接口 = any = interface{}（Go 1.18+）=====
	var anything any
	anything = 42
	fmt.Printf("anything 是 int: %d\n", anything.(int))
	anything = "hello"
	fmt.Printf("anything 是 string: %s\n", anything.(string))
}

// 输出：
// 用户: &{UUID:U001 Nickname:张三三 Email:}
// 生成 UUID: USER-1234-5678
//
// --- 接口多态 ---
// 用户信息: 张三三 (U001)
// 处理用户请求: /login
// 处理管理员请求: /dashboard
// anything 是 int: 42
// anything 是 string: hello
```

---

## 配套练习

按顺序完成，每个都能独立跑通后再进入下一阶段：

1. 用 `net/http` 写一个 `/hello?name=xxx` 接口，返回 `{"greeting": "hello, xxx"}`
2. 写一个 `LoginRequest` 结构体，用 `json.Marshal` 和 `json.Unmarshal` 各写一个例子
3. 写一个函数，接收 `*User` 指针参数，修改 `Name` 字段，验证值传递和指针传递的区别
4. 用 `make` 创建 slice 和 map，练习 `append`、`delete`、`range` 和 "comma ok" 模式
5. 给练习 3 的 `User` 添加一个方法 `Greet() string`，返回 `"你好，我是 XXX"`
6. 创建一个带 `sync.Mutex` 的 map，启动 10 个 goroutine 同时读写，观察结果
7. 写一个 select 同时监听 3 个 channel，模拟本项目消息服务器的三个事件流
