package gorm

import (
	"apylee_chat_server/internal/dao"
	"apylee_chat_server/internal/dto/request"
	"apylee_chat_server/internal/model"
	"apylee_chat_server/internal/service/chat"
	"apylee_chat_server/internal/service/sms"
	"apylee_chat_server/pkg/util/random"
	"apylee_chat_server/pkg/zlog"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"regexp"
	"time"
)

type userInfoService struct {
}

var UserInfoService = new(userInfoService)

// dao层加不了校验，在service层加
// checkTelephoneValid 检验电话是否有效
func (u *userInfoService) checkTelephoneValid(telephone string) bool {
	pattern := `^1([38][0-9]|14[579]|5[^4]|16[6]|7[1-35-8]|9[189])\d{8}$`
	match, err := regexp.MatchString(pattern, telephone)
	if err != nil {
		zlog.Fatal(err.Error())
	}
	return match
}

// checkEmailValid 校验邮箱是否有效
func (u *userInfoService) checkEmailValid(email string) bool {
	pattern := `^[\w\.-]+@[\w-]+\.[\w{2,3}(\.\w{2})?]$`
	match, err := regexp.MatchString(pattern, email)
	if err != nil {
		zlog.Fatal(err.Error())
	}
	return match
}

// TODO
func (u *userInfoService) checkUserIsAdminOrNot(user model.UserInfo) bool {
	return false
}

// Login 登录
func (u *userInfoService) Login(c *gin.Context, loginReq request.LoginRequest) (string, error) {
	password := loginReq.Password
	var user model.UserInfo
	res := dao.GormDB.First(&user, "telephone = ?", loginReq.Telephone)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			message := "user not existed"
			zlog.Info(message)
			return message, nil
		}
		zlog.Error(res.Error.Error())
		return "", res.Error
	}
	if user.Password != password {
		message := "password not correct"
		zlog.Info(message)
		return message, nil
	}
	// 手机号验证，最后一步才调用api，省钱hhh
	if err := sms.VerificationCode(loginReq.Telephone); err != nil {
		zlog.Error(err.Error())
		return "", err
	}
	// 登录成功，chat client建立
	if err := chat.NewClientInit(c, user.Uuid); err != nil {
		return "", err
	}
	return "", nil
}

// Register 注册
func (u *userInfoService) Register(c *gin.Context, registerReq request.RegisterRequest) (string, error) {
	// 校验手机号是否有效
	if !u.checkTelephoneValid(registerReq.Telephone) {
		message := "telephone invalid"
		zlog.Info(message)
		return message, nil
	}
	var newUser model.UserInfo
	res := dao.GormDB.First(&newUser, "telephone = ?", registerReq.Telephone)
	if res.Error == nil {
		// 用户已经存在，注册失败
		message := "user already existed"
		zlog.Info(message)
		return message, nil
	} else {
		// 其他报错
		if !errors.Is(res.Error, gorm.ErrRecordNotFound) {
			zlog.Error(res.Error.Error())
			return "", res.Error
		}
		// 可以继续注册
	}
	newUser.Uuid = "U" + random.GetNowAndLenRandomString(11)
	newUser.TelePhone = registerReq.Telephone
	newUser.Password = registerReq.Password
	newUser.NickName = registerReq.NickName
	newUser.CreatedAt = time.Now()
	newUser.IsAdmin = u.checkUserIsAdminOrNot(newUser)
	// 手机号验证，最后一步才调用api，省钱hhh
	err := sms.VerificationCode(registerReq.Telephone)
	if err != nil {
		zlog.Error(err.Error())
		return "", err
	}
	res = dao.GormDB.Create(&newUser)
	if res.Error != nil {
		zlog.Error(res.Error.Error())
		return "", res.Error
	}
	// 注册成功，chat client建立
	if err := chat.NewClientInit(c, newUser.Uuid); err != nil {
		return "", err
	}
	return "", nil
}
