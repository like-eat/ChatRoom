package main

import (
	"errors"
	"fmt"
	"os"
)

// 指针 defer 错误处理
// 值传递和指针传递
func updateByValue(u User) {
	u.Name = "被修改了"
}

func updateByPointer(u *User) {
	u.Name = "被修改了"
}

type User struct {
	ID   int
	Name string
}

// defer延迟执行
func deferDemo() {
	// 模拟打开文件
	fmt.Println("打开文件")

	// defer 确保文件一定会关闭
	defer fmt.Println("关闭文件(defer1)")
	defer fmt.Println("刷新缓冲区(defer2)")

	fmt.Println("读写文件")
	fmt.Println("处理完成")
	// 函数执行结束后自动执行defer，线defer2，再defer1
}

// 错误处理 Go没有try-catch
// 模拟失败的操作
func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("除数不能为0")
	}
	return a / b, nil
}

// 模拟多步骤操作
func readConfig(path string) (string, error) {
	// 读文件
	data, err := os.ReadFile(path)
	// 读失败了
	if err != nil {
		return "", fmt.Errorf("读取配置文件失败%w", err)
	}
	// 成功处理
	return string(data), nil
}

// 错误判断 errors.Is() 和 errors.As()
func checkError(err error) {
	// errors.Is() 判断错误类型
	if errors.Is(err, os.ErrNotExist) {
		fmt.Println("-> 文件不存在")
	}
}

func main() {
	// 指针演示
	user := User{ID: 1, Name: "张三"}

	updateByValue(user)
	fmt.Println("值传递修改后:", user.Name) // 张三

	updateByPointer(&user)
	fmt.Println("指针传递修改后:", user.Name) // 被修改了

	fmt.Println("=====================================")
	// defer演示
	deferDemo()

	fmt.Println("=====================================")
	// 错误处理演示
	// 直接处理
	result, err := divide(10, 0)
	if err != nil {
		fmt.Println("除零错误:", err)
	} else {
		fmt.Println("除法结果:", result)
	}

	// 错误传递
	_, err = readConfig("不存在的文件.toml")
	if err != nil {
		fmt.Println("读取配置文件错误:", err)
		checkError(err)
	}

	// 忽略错误
	result, _ = divide(10, 2)
	fmt.Println("除法结果:", result)
}
