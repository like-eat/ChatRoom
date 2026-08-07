package main 

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// 请求DTO
type RegisterReq struct {
	Nickname string `json:"nickname"`
	Telephone string `json:"telephone"`
	Password string `json:"password"`
	Gender int8 `json:"gender"`
}

type LoginReq struct {
	Telephone string `json:"telephone"`
	Password string `json:"password"`
}

type UpdateUserReq struct {
	Nickname string `json:"nickname"`
	Avatar string `json:"avatar"`
	Gender *int8 `json:"gender"`
}

// 响应DTO
type LoginResp struct {
	Token string `json:"token"`
	UserID int64 `json:"user_id"`
	Nickname string `json:"nickname"`
	Telephone string `json:"telephone"`
	Avatar string `json:"avatar"`
}

type UserListResp struct {
	Total int64 `json:"total"`
	List []*User `json:"list"`
}