package main

import (
	"fmt"
	"strings"
)

// 方法=带接收者的函数
type UserService struct {
	db map[string]*User // 模拟数据库连接
}

type User struct {
	UUID     string
	Nickname string
	Email    string
}

// NewUserService构造函数，初始化一个实例
func NewUserService() *UserService {
	return &UserService{
		db: make(map[string]*User),
	}
}

// 创建UserService的方法 CreateUser 和 UpdateNickname 方法
func (s *UserService) CreateUser(uuid, nickname string) *User {
	user := &User{UUID: uuid, Nickname: nickname}
	s.db[uuid] = user // 修改 s.db（需要指针）
	return user
}

func (s *UserService) UpdateNickname(uuid, nickname string) error {
	// 检查用户是否存在
	user, ok := s.db[uuid]
	if !ok {
		return fmt.Errorf("用户 %s 不存在", uuid)
	}
	// 更新用户昵称
	user.Nickname = nickname
	return nil
}

// UserService的根据UUID查找用户的方法 FindByUUID
func (s UserService) FindByUUID(uuid string) (*User, error) {
	// 检查用户是否存在
	user, ok := s.db[uuid]
	if !ok {
		return nil, fmt.Errorf("用户 %s 不存在", uuid)
	}
	return user, nil
}

// 普通函数 生成UUID
func GenerateUUID(prefix string) string {
	return fmt.Sprintf("%s-1234-5678", prefix)
}

// 定义一个接口，必须要有CreateUser，FindByUUID，UpdateNickname 方法
type UserRepository interface {
	CreateUser(uuid, nickname string) *User
	FindByUUID(uuid string) (*User, error)
	UpdateNickname(uuid, nickname string) error
}

// UserService 实现 UserRepository 接口
var _ UserRepository = (*UserService)(nil)

// 多态 用接口类型作为参数
func PrintUserInfo(repo UserRepository, uuid string) {
	// 多态就在参数用了接口类型
	// 同一个函数，传不同的能力就表现不同
	user, err := repo.FindByUUID(uuid)
	// repo不是一个具体的类型，只要传进来的接口实现了FindByUUID的方法，函数就能正常处理
	if err != nil {
		fmt.Println("用户不存在", err)
		return
	}
	// 用户存在输出信息
	fmt.Printf("用户信息：%s (%s)\n", strings.ToUpper(user.Nickname), user.UUID)
}

// 接口组合
type Reader interface {
	Read(p []byte) (n int, err error)
}

type Writer interface {
	Write(p []byte) (n int, err error)
}

type ReaderWriter interface {
	Reader
	Writer
}

// 项目当中常见的接口用法

type Handler interface {
	Handle(path string) string
}

type UserHandler struct{}

func (h UserHandler) Handle(path string) string {
	return "处理用户请求：" + path
}

type AdminHandler struct{}

func (h AdminHandler) Handle(path string) string {
	return "处理管理员请求：" + path
}

// 用接口实现多态
func route(handler Handler, path string) {
	fmt.Println(handler.Handle(path))
}

func main() {
	svc := NewUserService()

	svc.CreateUser("U001", "张三")
	svc.CreateUser("U002", "李四")

	svc.UpdateNickname("U001", "张三三")

	user, _ := svc.FindByUUID("U001")
	fmt.Printf("用户：%v\n", user)

	fmt.Println("生成UUID：", GenerateUUID("user"))

	// 接口演示
	fmt.Println("多态：")
	PrintUserInfo(svc, "U001")
	PrintUserInfo(svc, "U002")

	route(UserHandler{}, "/login")
	route(AdminHandler{}, "/dashboard")

	// 空接口
	var anything any
	anything = 42
	fmt.Printf("anything 是 int 类型：%d\n", anything.(int))
	anything = "hello"
	fmt.Printf("anything 是 string 类型：%s\n", anything.(string))

}
