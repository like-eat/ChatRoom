package main

import "fmt"

// slice切片 = 动态数组
func sliceDemo() {
	// 创建
	names := []string{"张三", "李四", "王五"}

	ids := make([]int, 0, 10)
	// ids追加内容
	ids = append(ids, 1, 2, 3)
	fmt.Println(ids)
	// len是现在有几个，cap是最多能装几个，slice特有的
	fmt.Printf("len=%d, cap=%d\n", len(ids), cap(ids))

	// 遍历 i 是索引，name是值
	for i, name := range names {
		fmt.Printf("names[%d]=%s\n", i, name)
	}

	// 切片操作
	sub := names[1:3] // 左闭右开
	fmt.Println("切片sub=", sub)

	// 删除元素
	names = append(names[:1], names[2:]...) // 删除索引为1
	fmt.Println("删除后names=", names)

	// nil slice 和 empty slice
	var nilSlice []string
	emptySlice := []string{}

	// %d是整数 %s是字符串 %v是万能
	fmt.Printf("nilSlice: %v, emptySlice: %v\n", nilSlice, emptySlice)
}

// map 哈希表
func mapDemo() {
	scores := make(map[string]int)

	users := map[string]string{
		"U001": "张三",
		"U002": "李四",
		"U003": "王五",
	}

	// 写入
	scores["张三"] = 90
	scores["李四"] = 80
	scores["王五"] = 70

	// 读取
	fmt.Println("张三的成绩=", scores["张三"])

	// 检查key是否存在
	score, ok := scores["王五"]
	if ok {
		fmt.Println("王五的成绩=", score)
	} else {
		fmt.Println("王五的成绩不存在")
	}

	// 遍历
	for name, score := range scores {
		fmt.Printf("%s的成绩=%d\n", name, score)
	}

	// 删除
	delete(scores, "李四")
	fmt.Println("删除李四后:", scores)

	fmt.Println("用户表:", users)
}

// slice+map的组合用法
type Client struct {
	UserID   string
	Nickname string
}

type Hub struct {
	// 在线客户端
	clients map[string]*Client
	// 消息队列
	broadcast chan string
}

func NewHub() *Hub {
	return &Hub{
		clients:   make(map[string]*Client),
		broadcast: make(chan string, 100),
	}
}

func (h *Hub) AddClient(client *Client) {
	h.clients[client.UserID] = client
}

// 获取所有在线用户名
// (h *Hub) 是把函数绑定到Hub类型上，类似于面向对象的类方法
func (h *Hub) GetOnlineNames() []string {
	names := make([]string, 0, len(h.clients))
	for _, c := range h.clients {
		names = append(names, c.Nickname)
	}
	return names
}

func main() {
	fmt.Println("slice演示")
	sliceDemo()

	fmt.Println("map演示")
	mapDemo()

	fmt.Println("slice+map组合演示")

	hub := NewHub()
	// 增加用户
	hub.AddClient(&Client{UserID: "U001", Nickname: "张三"})
	hub.AddClient(&Client{UserID: "U002", Nickname: "李四"})
	hub.AddClient(&Client{UserID: "U003", Nickname: "王五"})

	fmt.Println("在线用户:", hub.GetOnlineNames())

	// 查找某个用户是否在线
	if client, ok := hub.clients["U002"]; ok {
		fmt.Printf("用户%s在线，昵称=%s\n", client.UserID, client.Nickname)
	}
}
