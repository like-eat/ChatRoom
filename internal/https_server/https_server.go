package https_server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	v1 "kama_chat_server/api/v1"
	"kama_chat_server/internal/config"
	"kama_chat_server/internal/middleware"
	"os"
)

var GE *gin.Engine

func init() {
	GE = gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	GE.Use(cors.New(corsConfig))
	if err := os.MkdirAll(config.GetConfig().StaticAvatarPath, 0755); err != nil {
		panic(err)
	}
	GE.Static("/static/avatars", config.GetConfig().StaticAvatarPath)

	// Public authentication endpoints.
	GE.POST("/login", v1.Login)
	GE.POST("/register", v1.Register)

	authorized := GE.Group("/")
	authorized.Use(middleware.AuthRequired())
	{
		authorized.POST("user/updateUserInfo", v1.UpdateUserInfo)
		authorized.POST("user/getUserInfo", v1.GetUserInfo)
		authorized.POST("user/wsLogout", v1.WsLogout)
		authorized.POST("group/createGroup", v1.CreateGroup)
		authorized.POST("group/loadMyGroup", v1.LoadMyGroup)
		authorized.POST("group/leaveGroup", v1.LeaveGroup)
		authorized.POST("group/dismissGroup", v1.DismissGroup)
		authorized.POST("group/getGroupInfo", v1.GetGroupInfo)
		authorized.POST("group/updateGroupInfo", v1.UpdateGroupInfo)
		authorized.POST("group/getGroupMemberList", v1.GetGroupMemberList)
		authorized.POST("group/removeGroupMembers", v1.RemoveGroupMembers)
		authorized.POST("session/openSession", v1.OpenSession)
		authorized.POST("session/getUserSessionList", v1.GetUserSessionList)
		authorized.POST("session/getGroupSessionList", v1.GetGroupSessionList)
		authorized.POST("session/deleteSession", v1.DeleteSession)
		authorized.POST("session/checkOpenSessionAllowed", v1.CheckOpenSessionAllowed)
		authorized.POST("contact/getUserList", v1.GetUserList)
		authorized.POST("contact/loadMyJoinedGroup", v1.LoadMyJoinedGroup)
		authorized.POST("contact/getContactInfo", v1.GetContactInfo)
		authorized.POST("contact/deleteContact", v1.DeleteContact)
		authorized.POST("contact/applyContact", v1.ApplyContact)
		authorized.POST("contact/getNewContactList", v1.GetNewContactList)
		authorized.POST("contact/passContactApply", v1.PassContactApply)
		authorized.POST("contact/getAddGroupList", v1.GetAddGroupList)
		authorized.POST("contact/refuseContactApply", v1.RefuseContactApply)
		authorized.POST("message/getMessageList", v1.GetMessageList)
		authorized.POST("message/getGroupMessageList", v1.GetGroupMessageList)
		authorized.POST("message/uploadAvatar", v1.UploadAvatar)
	}

	admin := GE.Group("/")
	admin.Use(middleware.AuthRequired(), middleware.AdminRequired())
	{
		admin.POST("user/getUserInfoList", v1.GetUserInfoList)
		admin.POST("user/ableUsers", v1.AbleUsers)
		admin.POST("user/disableUsers", v1.DisableUsers)
		admin.POST("user/deleteUsers", v1.DeleteUsers)
		admin.POST("user/setAdmin", v1.SetAdmin)
		admin.POST("group/getGroupInfoList", v1.GetGroupInfoList)
		admin.POST("group/deleteGroups", v1.DeleteGroups)
		admin.POST("group/setGroupsStatus", v1.SetGroupsStatus)
	}

	// Browser WebSocket clients send the JWT in Sec-WebSocket-Protocol;
	// the server derives the client UUID from the signed subject claim.
	GE.GET("/wss", v1.WsLogin)

}
