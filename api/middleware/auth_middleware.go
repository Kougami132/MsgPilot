package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kougami132/MsgPilot/bootstrap"
	"github.com/kougami132/MsgPilot/internal/utils"
)

// AuthMiddleware 认证中间件
func AuthMiddleware(app bootstrap.Application) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "未提供认证信息"})
			return
		}

		// 检查Bearer前缀
		if !strings.HasPrefix(authHeader, "Bearer ") {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "认证格式错误"})
			return
		}

		// 获取token
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "未提供令牌"})
			return
		}

		// 验证token
		claims, err := utils.ValidateToken(token, app.Env.AccessTokenSecret)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "无效的令牌"})
			return
		}

		// 将用户信息存储到上下文中
		ctx.Set("username", claims.Username)

		ctx.Next()
	}
}
