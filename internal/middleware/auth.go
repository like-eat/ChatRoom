package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"kama_chat_server/internal/dao"
	"kama_chat_server/internal/model"
	"kama_chat_server/internal/service/auth"
)

const (
	ContextUserID  = "authenticated_user_id"
	ContextIsAdmin = "authenticated_is_admin"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := strings.TrimSpace(c.GetHeader("Authorization"))
		parts := strings.SplitN(authorization, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || strings.TrimSpace(parts[1]) == "" {
			abortUnauthorized(c, "缺少有效的登录令牌")
			return
		}

		claims, err := auth.ParseToken(strings.TrimSpace(parts[1]))
		if err != nil {
			abortUnauthorized(c, "登录令牌无效或已过期")
			return
		}

		var user model.UserInfo
		if err := dao.GormDB.Select("uuid", "is_admin", "status").Where("uuid = ? AND status = 0", claims.Subject).First(&user).Error; err != nil {
			abortUnauthorized(c, "用户不存在或账号已被禁用")
			return
		}

		c.Set(ContextUserID, user.Uuid)
		c.Set(ContextIsAdmin, user.IsAdmin)
		c.Next()
	}
}

func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		value, exists := c.Get(ContextIsAdmin)
		isAdmin, ok := value.(int8)
		if !exists || !ok || isAdmin != 1 {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code":    http.StatusForbidden,
				"message": "需要管理员权限",
			})
			return
		}
		c.Next()
	}
}

func abortUnauthorized(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"code":    http.StatusUnauthorized,
		"message": message,
	})
}
