package main

import (
	"apylee_chat_server/internal/service"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	r := gin.Default()
	r.POST("/login", service.Login)
	r.POST("/register", service.Register)
	if err := r.Run("127.0.0.1:8080"); err != nil {
		log.Fatal("server running fault")
		return
	}
}
