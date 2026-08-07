package main

import (
	"fmt"
	"sync"
	"time"
)

type Store struct {
	mu sync.RWMutex
	users map[string]*User
	byPhone map[string]int64
	autoID int64
}

func NewStore() *Store {
	return &Store{
		users:  make(map[string]*User),
		byPhone: make(map[string]int64),
		autoID:  0,
	}
}

// 创建用户
func (s *Store) CreateUser(user *User) (*User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
}