package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	v1 "kama_chat_server/api/v1"
	"kama_chat_server/config"
	"log"
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
	r.POST("/login", v1.Login)
	r.POST("/register", v1.Register)
	r.POST("/group/createGroup", v1.CreateGroup)
	r.POST("/contact/getUserList", v1.GetUserList)
	r.POST("/group/loadMyGroup", v1.LoadMyGroup)
	r.POST("/contact/loadMyJoinedGroup", v1.LoadMyJoinedGroup)
	r.POST("/contact/getContactInfo", v1.GetContactInfo)
	r.POST("/session/openSession", v1.OpenSession)
	r.POST("/session/getUserSessionList", v1.GetUserSessionList)
	r.POST("/session/getGroupSessionList", v1.GetGroupSessionList)
	r.POST("/session/deleteSession", v1.DeleteSession)
	r.POST("/contact/deleteContact", v1.DeleteContact)
	r.POST("/contact/applyContact", v1.ApplyContact)
	r.POST("/contact/getNewContactList", v1.GetNewContactList)
	conf := config.GetConfig()
	host := conf.MainConfig.Host
	port := conf.MainConfig.Port

	if err := r.Run(fmt.Sprintf("%s:%d", host, port)); err != nil {
		log.Fatal("server running fault")
		return
	}
}
