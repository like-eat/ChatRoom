package main

import (
	"fmt"
	"sync"
	"time"
)

// Gin的每个http请求就是一个goroutine
// 多个goroutine共享的数据需要保护

// 简单goroutine
func example1() {
	// 启动一个goroutine
	go func() {
		fmt.Println("我在另一个goroutine运行")
	}() // ()立即执行

	// 让当前程序睡100毫秒，让goroutine有时间执行
	// 因为goroutine是异步的，主程序不会等它，如果主程序立刻跑完
	// 进程就退出，goroutine就不会执行
	time.Sleep(100 * time.Millisecond) // 等待goroutine执行完
}

// 无缓冲channel 没有给chan设置大小
// 无缓冲channel相当于同步交接，发送方和接收方必须同时在场
// 有缓存channel类似于队列，发送发不需要的古代接收方
func example2() {
	// channel是goroutine之间传递数据的管道
	// 它主要解决goroutine之间的安全通信
	// channel一根管子智能传一种类型，而且是FIFO
	ch := make(chan string) // 创建一个无缓冲channel

	go func() {
		time.Sleep(1 * time.Second)
		ch <- "任务完成" // 发送
		fmt.Println("发送完成")
	}()

	fmt.Println("等待任务完成...")
	msg := <-ch // 接收
	fmt.Println("收到消息:", msg)
}

// 上锁保护共享数据
func example3() {
	type Counter struct {
		mu    sync.Mutex
		count int
	}

	counter := &Counter{}
	var wg sync.WaitGroup

	// 启动1000个goroutine 同时+1
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.mu.Lock()   // 上锁
			counter.count++     // 安全修改
			counter.mu.Unlock() // 解锁
		}()
	}

	wg.Wait() // 等待所有goroutine完成
	fmt.Println("最终计数:", counter.count)
}

// 带缓冲的channel
func example4() {
	// 缓存大小 3
	ch := make(chan string, 3)

	ch <- "任务1"
	ch <- "任务2"
	ch <- "任务3"

	fmt.Println(<-ch)
	fmt.Println(<-ch)
	fmt.Println(<-ch)
}

// select 多路复用
func example5() {
	// select是哪个通道先有数据就先执行哪个case
	login := make(chan string, 5)
	logout := make(chan string, 5)
	transmit := make(chan string, 5)

	// 模拟事件
	go func() { login <- "用户A登录" }()
	go func() { transmit <- "用户A发送消息" }()
	go func() { logout <- "用户A登出" }()

	// 循环处理三种事件
	for i := 0; i < 3; i++ {
		// switch 处理普通变量
		// select 处理channel，判断哪个通道先有数据
		select {
		case msg := <-login:
			fmt.Println("处理登录事件:", msg)
		case msg := <-logout:
			fmt.Println("处理登出事件:", msg)
		case msg := <-transmit:
			fmt.Println("处理消息传输事件:", msg)
		}
	}
}

func main() {
	fmt.Println("====实例1 goroutine ====")
	example1()

	fmt.Println("==== 实例2 无缓冲channel ====")
	example2()

	fmt.Println("==== 实例3 加锁 ====")
	example3()

	fmt.Println("==== 实例4 有缓冲channel ====")
	example4()

	fmt.Println("==== 实例5 select多路复用 ====")
	example5()
}
