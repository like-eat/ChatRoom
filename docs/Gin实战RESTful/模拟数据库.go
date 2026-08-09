package main

import (
	"fmt"
	"sync"
	"time"
)

type Store struct {
	mu      sync.RWMutex
	users   map[string]*User
	byPhone map[string]int64
	autoID  int64
}

func NewStore() *Store {
	return &Store{
		users:   make(map[string]*User),
		byPhone: make(map[string]int64),
		autoID:  0,
	}
}

// 创建用户
func (s *Store) CreateUser(req *RegisterReq) (*User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 检查手机号是否注册
	if _, exists := s.byPhone[req.Telephone]; exists {
		return nil, fmt.Errorf("手机号已经注册")
	}

	s.autoID++
	user := &User{
		ID:        s.autoID,
		Nickname:  req.Nickname,
		Telephone: req.Telephone,
		Password:  req.Password,
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
	s.mu.Lock()
	defer s.mu.Unlock() // 在函数结束的时候执行解锁

	id, exists := s.byPhone[telephone]
	if !exists {
		return nil, fmt.Errorf("用户不存在")
	}

	// 查到了就返回用户
	users := *s.users[id]

	return &users, nil
}

// 根据ID 查找用户
func (s *Store) FindByID(id int64) (*User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	user, exists := s.users[id]
	if !exists {
		return nil, fmt.Errorf("用户不存在")
	}

	// 查到了就返回用户
	result := *user
	return &result, nil
}

// 更新用户信息
func (s *Store) updateUser(id int64, req UpdateUserReq)
(*User, error) {
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
	if req.Gender != "" {
		user.Gender = req.Gender
	}

	return user, nil
}

// 分页查询用户列表
func (s *Store) ListUsers(page, pageSize int)
([]*User, int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 计算用户总数
	total := int64(len(s.users))
	start := (page - 1) * pageSize
	// 如果开始位置大于用户总数，返回空列表
	if start >= int(total) {
		return []*User{}, total
	}

	var result []*User
	i := 0
	for _, user := range s.users {
		if i >= start && len(result) < pageSize {
			u := *user
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

	_, exists := s.users[id]
	if !exists {
		return fmt.Errorf("用户不存在")
	}

	delete(s.users, id)
	delete(s.byPhone, id)

	return nil
}
