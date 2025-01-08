package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	v1 "kama_chat_server/api/v1"
	"kama_chat_server/config"
	"kama_chat_server/internal/service/chat"
	"kama_chat_server/pkg/zlog"
)

func main() {
	r := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Content-Type", "Authorization"}
	corsConfig.ExposeHeaders = []string{"X-Custom-Header"}
	corsConfig.AllowCredentials = true
	r.Use(cors.New(corsConfig))
	r.Static("/static/avatars", config.GetConfig().StaticAvatarPath)
	r.Static("/static/files", config.GetConfig().StaticFilePath)
	r.POST("/login", v1.Login)
	r.POST("/register", v1.Register)
	r.POST("/user/updateUserInfo", v1.UpdateUserInfo)
	r.POST("/group/createGroup", v1.CreateGroup)
	r.POST("/group/loadMyGroup", v1.LoadMyGroup)
	r.POST("/group/checkGroupAddMode", v1.CheckGroupAddMode)
	r.POST("/group/enterGroupDirectly", v1.EnterGroupDirectly)
	r.POST("/group/leaveGroup", v1.LeaveGroup)
	r.POST("/group/dismissGroup", v1.DismissGroup)
	r.POST("/session/openSession", v1.OpenSession)
	r.POST("/session/getUserSessionList", v1.GetUserSessionList)
	r.POST("/session/getGroupSessionList", v1.GetGroupSessionList)
	r.POST("/session/deleteSession", v1.DeleteSession)
	r.POST("/session/checkOpenSessionAllowed", v1.CheckOpenSessionAllowed)
	r.POST("/contact/getUserList", v1.GetUserList)
	r.POST("/contact/loadMyJoinedGroup", v1.LoadMyJoinedGroup)
	r.POST("/contact/getContactInfo", v1.GetContactInfo)
	r.POST("/contact/deleteContact", v1.DeleteContact)
	r.POST("/contact/applyContact", v1.ApplyContact)
	r.POST("/contact/getNewContactList", v1.GetNewContactList)
	r.POST("/contact/passContactApply", v1.PassContactApply)
	r.POST("/contact/blackContact", v1.BlackContact)
	r.POST("/contact/cancelBlackContact", v1.CancelBlackContact)
	r.POST("/contact/getAddGroupList", v1.GetAddGroupList)
	r.POST("/contact/refuseContactApply", v1.RefuseContactApply)
	r.POST("/contact/blackApply", v1.BlackApply)
	r.POST("/message/getMessageList", v1.GetMessageList)
	r.POST("/message/getGroupMessageList", v1.GetGroupMessageList)
	r.POST("/message/uploadAvatar", v1.UploadAvatar)
	r.GET("/ws", v1.WsHandler)

	conf := config.GetConfig()
	host := conf.MainConfig.Host
	port := conf.MainConfig.Port

	go chat.ChatServer.Start()

	if err := r.Run(fmt.Sprintf("%s:%d", host, port)); err != nil {
		zlog.Fatal("server running fault")
		return
	}

}
