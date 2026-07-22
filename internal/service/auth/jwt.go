package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"kama_chat_server/internal/config"
)

const accessTokenType = "access"

type Claims struct {
	IsAdmin  int8   `json:"is_admin"`
	TokenUse string `json:"token_use"`
	jwt.RegisteredClaims
}

func GenerateToken(userID string, isAdmin int8) (string, error) {
	conf := config.GetConfig().JWTConfig
	if len(conf.Secret) < 32 {
		return "", errors.New("JWT secret must contain at least 32 characters")
	}

	now := time.Now()
	claims := Claims{
		IsAdmin:  isAdmin,
		TokenUse: accessTokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    conf.Issuer,
			Subject:   userID,
			Audience:  jwt.ClaimStrings{conf.Audience},
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(conf.ExpireHours) * time.Hour)),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(conf.Secret))
}

func ParseToken(tokenString string) (*Claims, error) {
	conf := config.GetConfig().JWTConfig
	if len(conf.Secret) < 32 {
		return nil, errors.New("JWT secret must contain at least 32 characters")
	}

	claims := new(Claims)
	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("unexpected JWT signing method: %s", token.Method.Alg())
			}
			return []byte(conf.Secret), nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithIssuer(conf.Issuer),
		jwt.WithAudience(conf.Audience),
		jwt.WithExpirationRequired(),
		jwt.WithIssuedAt(),
		jwt.WithLeeway(30*time.Second),
	)
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token: %w", err)
	}
	if claims.TokenUse != accessTokenType || claims.Subject == "" {
		return nil, errors.New("invalid access token claims")
	}
	return claims, nil
}
