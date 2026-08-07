package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 统一响应结构
type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    interface{} `json:"data"`
}

func Sucess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{Code: 200, Message: "success", Data: data})
}

func Error(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, Response{Code: code, Message: msg})
}

