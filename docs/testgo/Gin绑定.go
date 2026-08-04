package main

import (
	"encoding/json"
	"fmt"
)

type LoginRequest struct {
	Telephone string `json:"telephone"`
	Password  string `json:"password"`
}

type BadRequest struct {
	Telephone string `json:"telephone"`
}

type UpdateUserRequest struct {
	Nickname string  `json:"nickname"`
	Avatar   *string `json:"avatar,omitempty"`
	Age      int     `json:"age" binding:"required, gte=0,lte=120"`
}

// omitempty是没传还是传了为空
func main() {
	jsonStr := `{"telephone":"1234567890","password":"password123"}`
	var req LoginRequest
	json.Unmarshal([]byte(jsonStr), &req)
	fmt.Printf("反序列化：%+v\n", req)

	resp := LoginRequest{
		Telephone: "1234567890",
		Password:  "***",
	}
	respBytes, _ := json.Marshal(resp)
	fmt.Println("序列化：", string(respBytes))

	// ===== omitempty 演示 =====
	avatar := "https://example.com/avatar.png"
	update := UpdateUserRequest{
		Nickname: "John Doe",
		Avatar:   &avatar,
	}
	withAvatar, _ := json.Marshal(update)
	fmt.Printf("有avatar： %s\n", string(withAvatar))

	updateNoAvatar := UpdateUserRequest{
		Nickname: "Jane Doe",
	}
	noAvatar, _ := json.Marshal(updateNoAvatar)
	fmt.Printf("无avatar： %s\n", string(noAvatar))

}
