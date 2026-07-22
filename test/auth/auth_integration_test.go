package auth_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"kama_chat_server/internal/dao"
	"kama_chat_server/internal/https_server"
	"kama_chat_server/internal/model"
	security "kama_chat_server/internal/service/auth"
	"kama_chat_server/internal/service/chat"
)

var startChatServer sync.Once

func TestLoginAndProtectedRoutes(t *testing.T) {
	webSocketRecorder := httptest.NewRecorder()
	https_server.GE.ServeHTTP(webSocketRecorder, httptest.NewRequest(http.MethodGet, "/wss", nil))
	if webSocketRecorder.Code != http.StatusUnauthorized {
		t.Fatalf("WebSocket handshake without token returned %d", webSocketRecorder.Code)
	}

	password := "IntegrationPass123"
	hash, err := security.HashPassword(password)
	if err != nil {
		t.Fatal(err)
	}

	unique := time.Now().UnixNano()
	user := model.UserInfo{
		Uuid:      fmt.Sprintf("UT%018d", unique%1_000_000_000_000_000_000),
		Nickname:  "jwt-test-user",
		Telephone: fmt.Sprintf("1%010d", unique%10_000_000_000),
		Avatar:    "https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png",
		Password:  hash,
		CreatedAt: time.Now(),
		IsAdmin:   0,
		Status:    0,
	}
	if err := dao.GormDB.Create(&user).Error; err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := dao.GormDB.Unscoped().Where("uuid = ?", user.Uuid).Delete(&model.UserInfo{}).Error; err != nil {
			t.Error(err)
		}
	})

	loginBody, _ := json.Marshal(map[string]string{
		"telephone": user.Telephone,
		"password":  password,
	})
	loginRecorder := httptest.NewRecorder()
	https_server.GE.ServeHTTP(loginRecorder, httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(loginBody)))
	if loginRecorder.Code != http.StatusOK {
		t.Fatalf("login returned HTTP %d: %s", loginRecorder.Code, loginRecorder.Body.String())
	}

	var loginResponse struct {
		Code int `json:"code"`
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	if err := json.Unmarshal(loginRecorder.Body.Bytes(), &loginResponse); err != nil {
		t.Fatal(err)
	}
	if loginResponse.Code != http.StatusOK || loginResponse.Data.Token == "" {
		t.Fatalf("login did not return a JWT: %s", loginRecorder.Body.String())
	}

	startChatServer.Do(func() { go chat.ChatServer.Start() })
	webSocketServer := httptest.NewServer(https_server.GE)
	defer webSocketServer.Close()
	dialer := websocket.Dialer{Subprotocols: []string{"kama-chat", loginResponse.Data.Token}}
	connection, response, err := dialer.Dial("ws"+strings.TrimPrefix(webSocketServer.URL, "http")+"/wss", nil)
	if err != nil {
		if response != nil {
			t.Fatalf("authenticated WebSocket returned %d: %v", response.StatusCode, err)
		}
		t.Fatal(err)
	}
	if connection.Subprotocol() != "kama-chat" {
		t.Fatalf("unexpected WebSocket subprotocol %q", connection.Subprotocol())
	}
	_ = connection.Close()

	userInfoBody, _ := json.Marshal(map[string]string{"uuid": user.Uuid})
	unauthorizedRecorder := httptest.NewRecorder()
	https_server.GE.ServeHTTP(unauthorizedRecorder, httptest.NewRequest(http.MethodPost, "/user/getUserInfo", bytes.NewReader(userInfoBody)))
	if unauthorizedRecorder.Code != http.StatusUnauthorized {
		t.Fatalf("protected route without token returned %d", unauthorizedRecorder.Code)
	}

	authorizedRequest := httptest.NewRequest(http.MethodPost, "/user/getUserInfo", bytes.NewReader(userInfoBody))
	authorizedRequest.Header.Set("Authorization", "Bearer "+loginResponse.Data.Token)
	authorizedRecorder := httptest.NewRecorder()
	https_server.GE.ServeHTTP(authorizedRecorder, authorizedRequest)
	if authorizedRecorder.Code != http.StatusOK {
		t.Fatalf("protected route with token returned %d: %s", authorizedRecorder.Code, authorizedRecorder.Body.String())
	}

	adminRequest := httptest.NewRequest(http.MethodPost, "/user/getUserInfoList", bytes.NewReader([]byte(`{"owner_id":"`+user.Uuid+`"}`)))
	adminRequest.Header.Set("Authorization", "Bearer "+loginResponse.Data.Token)
	adminRecorder := httptest.NewRecorder()
	https_server.GE.ServeHTTP(adminRecorder, adminRequest)
	if adminRecorder.Code != http.StatusForbidden {
		t.Fatalf("non-admin request returned %d: %s", adminRecorder.Code, adminRecorder.Body.String())
	}
}

func TestLegacyPasswordIsUpgradedAfterLogin(t *testing.T) {
	unique := time.Now().UnixNano()
	plainPassword := "legacy-password"
	user := model.UserInfo{
		Uuid:      fmt.Sprintf("UL%018d", unique%1_000_000_000_000_000_000),
		Nickname:  "legacy-test-user",
		Telephone: fmt.Sprintf("1%010d", unique%10_000_000_000),
		Avatar:    "https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png",
		Password:  plainPassword,
		CreatedAt: time.Now(),
		Status:    0,
	}
	if err := dao.GormDB.Create(&user).Error; err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := dao.GormDB.Unscoped().Where("uuid = ?", user.Uuid).Delete(&model.UserInfo{}).Error; err != nil {
			t.Error(err)
		}
	})

	body, _ := json.Marshal(map[string]string{"telephone": user.Telephone, "password": plainPassword})
	recorder := httptest.NewRecorder()
	https_server.GE.ServeHTTP(recorder, httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body)))
	if recorder.Code != http.StatusOK {
		t.Fatalf("legacy login returned %d: %s", recorder.Code, recorder.Body.String())
	}

	var updated model.UserInfo
	if err := dao.GormDB.Where("uuid = ?", user.Uuid).First(&updated).Error; err != nil {
		t.Fatal(err)
	}
	if updated.Password == plainPassword {
		t.Fatal("legacy plaintext password was not upgraded")
	}
	valid, needsUpgrade, err := security.VerifyPassword(updated.Password, plainPassword)
	if err != nil || !valid || needsUpgrade {
		t.Fatalf("upgraded password is invalid: valid=%v upgrade=%v err=%v", valid, needsUpgrade, err)
	}
}
