package main

import (
	"fmt"
	"kama_chat_server/config"
	"kama_chat_server/internal/https_server"
	"kama_chat_server/internal/service/chat"
	"kama_chat_server/pkg/zlog"
)

func main() {
	conf := config.GetConfig()
	host := conf.MainConfig.Host
	port := conf.MainConfig.Port

	go chat.ChatServer.Start()

	if err := https_server.GE.RunTLS(fmt.Sprintf("%s:%d", host, port), "pkg/ssl/cert.crt", "pkg/ssl/cert.key"); err != nil {
		zlog.Fatal("server running fault")
		return
	}

}
