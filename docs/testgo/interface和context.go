package main

import (
	"context"
	"fmt"
	"time"
)

// interface{}(空接口) = 任意类型
type Store struct {
	data map[string]interface{} // 可以存任意类型的值
}

func (s *Store) Set(key string, value interface{}) {
	s.data[key] = value
}

func (s *Store) Get(key string) interface{} {
	return s.data[key]
}

// context.Context：传递超时、取消信号、请求级别的值
func queryDatabase(ctx context.Context, query string) (string, error) {
	select {
	// 同时监听两件事，谁先到就执行谁 
	case <-time.After(2 * time.Second): // 模拟2秒查询
		return fmt.Sprintf("查询结果: %s", query), nil
	case <-ctx.Done(): // 如果超时或取消，返回错误
		return "", ctx.Err()
	}
}

func main() {
	// === interface{} 示例 ===
	store := &Store{data: make(map[string]interface{})}
	store.Set("User", "张三")
	store.Set("Age", 30)
	store.Set("active", true)

	// 类型断言
	user := store.Get("User").(string)
	age := store.Get("Age").(int)
	fmt.Printf("User: %s, Age: %d\n", user, age)

	// 安全类型断言
	if v, ok := store.Get("active").(bool); ok {
		fmt.Printf("Active: %t\n", v)
	}

	// === context.Context 示例 ===
	// 设置超时时间为1秒
	// 给操作设置一个最晚截至时间，到点没结束就强制停止
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel() // 确保在函数退出时取消上下文

	result, err := queryDatabase(ctx, "SELECT * FROM users")
	if err != nil {
		fmt.Println("查询失败:", err)
	} else {
		fmt.Println(result)
	}
}
