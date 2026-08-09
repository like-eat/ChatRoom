package main

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// 全局store
var store = NewStore()

// 注册
func Register(c *gin.Context) {
	var req RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, 400, "参数错误:"+err.Error())
		return
	}

	user, err := store.CreateUser(req)
	if err != nil {
		Error(c, 400, err.Error("创建用户失败"))
		return
	}

	token, _ := GenerateToken(user.ID, user.Nickname)

	Success(c, LoginResp{
		Token:     token,
		UserID:    user.ID,
		Nickname:  user.Nickname,
		Telephone: user.Telephone,
		Avatar:    user.Avatar,
	})
}

// 登录
func Login(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, 400, "参数错误:"+err.Error())
		return
	}

	if user.Password != req.Password {
		Error(c, 400, "密码错误")
		return
	}

	token, _ := GenerateToken(user.ID, user.Nickname)

	Success(c, LoginResp{
		Token:     token,
		UserID:    user.ID,
		Nickname:  user.Nickname,
		Telephone: user.Telephone,
		Avatar:    user.Avatar,
	})
}

// 获取当前用户信息
func GetProfile(c *gin.Context) {
	// 从中间件注入的上下文
	userID := c.GetInt64("user_id")

	user, err := store.FindByID(userID)
	if err != nil {
		Error(c, 400, "用户不存在")
		return
	}

	Success(c, user)
}

// 更新用户信息
func UpdateProfile(c *gin.Context) {
	userID := c.GetInt64("user_id")

	var req UpdateProfileReq
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, 400, "参数错误:"+err.Error())
		return
	}
	// 更新用户信息
	user, err := store.UpdateUser(userID, req)
	if err != nil {
		Error(c, 400, err.Error())
		return
	}

	Success(c, user)
}

// 获取用户列表（分页）
func ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	user, total := store.ListUsers(page, pageSize)

	Success(c, UserListResp{Total: total, List: users})
}

// 获取单个用户
func GetUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		Error(c, 400, "用户ID格式错误")
		return
	}

	user, err := store.FindByID(id)
	if err != nil {
		Error(c, 400, "用户不存在")
		return
	}

	Success(c, user)
}

// 删除用户
func DeleteUser(c *gin.Context) {
	userID := c.GetInt64("user_id")
	idStr := c.Param("id")
	// strconv.ParseInt是把字符串转换成整数
	// Atoi也是把字符串转换成整数
	targetID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		Error(c, 400, "用户ID格式错误")
		return
	}

	// 只能删自己
	if userID != targetID {
		Error(c, 400, "只能删自己")
		return
	}

	if err := store.DeleteUser(targetID); err != nil {
		Error(c, 400, err.Error())
		return
	}

	Success(c, gin.H)
}
