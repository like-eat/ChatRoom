package main

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// 生产环境用环境变量
var jwtSecret = []byte("kama-chat-jwt-secret-key-2026")

type Claims struct {
	UserID   int64  `json:"user_id"`
	Nickname string `json:"nickname"`
	jwt.RegisteredClaims
}

// 生成JWT
func GenerateToken(userID int64, nickname string) 
(string, error) {
	claims := Claims{
		UserID:   userID,
		Nickname: nickname,
		// jwt的标准元信息：声明签发这，主题，签发时间，过期时间
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:     "gin-demo",
			Subject:    fmt.Sprintf("%d", userID),
			IssuedAt:  time.Now(),
			NotBefore:  time.Now(),
			ExpiresAt: time.Now().Add(time.Hour * 24),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// 验证JWT
func ParseToken(tokenStr string) (*Claims, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{},
		// 根据token的签名算法，返回对应的密钥
	func(t *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("不支持的签名算法：%v", t.Header["alg"])
		}

		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	// token有效,把token中的Claims类型转换
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, fmt.Errorf("无效的token")
	}

	return claims, nil
}