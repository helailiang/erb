package middleware

import (
	"erb/internal/config"
	"erb/pkg/jwt"
	"erb/pkg/response"
	"strings"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware 认证中间件
func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "未提供认证token")
			c.Abort()
			return
		}

		// 检查token格式 Bearer <token>
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(c, "token格式错误")
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 解析token
		claims, err := jwt.ParseToken(tokenString, cfg.JWT.Secret)
		if err != nil {
			if err == jwt.ErrTokenExpired {
				response.Unauthorized(c, "token已过期")
			} else {
				response.Unauthorized(c, "token无效")
			}
			c.Abort()
			return
		}

		// 将用户信息存储到上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)

		c.Next()
	}
}

