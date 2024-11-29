package main

import (
	v1 "apylee_chat_server/api/v1"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	r := gin.Default()
	r.POST("/login", v1.Login)
	r.POST("/register", v1.Register)
	if err := r.Run("127.0.0.1:8080"); err != nil {
		log.Fatal("server running fault")
		return
	}
}
