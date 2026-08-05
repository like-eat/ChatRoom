package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

// 用户结构体
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// 模拟数据库
var users = map[int]*User{
	1: {ID: 1, Name: "Alice", Age: 25},
	2: {ID: 2, Name: "Bob", Age: 30},
}

func main() {
	// 注册路由
	// 用户访问http://127.0.0.1:8080/user 自动调用handleUser函数
	http.HandleFunc("/user", handleUser)
	http.HandleFunc("/user/create", handleCreate)

	// 打印提示
	fmt.Println("服务启动在http://127.0.0.1:8080")
	// 启动服务器并阻塞
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// 处理get请求，读map返回用户
// w是写响应，r是读请求
func handleUser(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	// 如果请求方式是Get的话
	case http.MethodGet:
		// 解析请求参数 id=1
		idStr := r.URL.Query().Get("id")
		// 把"1"转换为int 1
		id, err := strconv.Atoi(idStr)

		// 有错误信息
		if err != nil {
			http.Error(w, "id 无效", http.StatusBadRequest)
			return
		}

		// 根据id查找用户
		user, ok := users[id]
		// 如果没查到就说用户不存在
		if !ok {
			http.Error(w, "用户不存在", http.StatusNotFound)
			return
		}

		// w是响应，把响应的内容写入w
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)

	// 如果请求方式是Post的话，不处理
	case http.MethodPost:
		http.Error(w, "POST 方法未实现", http.StatusNotImplemented)
	}
}

// 创建新用户写入map
func handleCreate(w http.ResponseWriter, r *http.Request) {
	// 只处理post方法
	if r.Method != http.MethodPost {
		http.Error(w, "只支持 POST 方法", http.StatusMethodNotAllowed)
		return
	}

	// 读取请求体，获得请求体
	body, _ := io.ReadAll(r.Body)

	var user User
	// json.Unmarshal把json字符串转换为结构体
	// 如果err不为nil，说明json字符串有问题
	if err := json.Unmarshal(body, &user); err != nil {
		http.Error(w, "请求体无效", http.StatusBadRequest)
		return
	}

	// 给新用户分配ID，并存入map
	user.ID = len(users) + 1
	users[user.ID] = &user

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
