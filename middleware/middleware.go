package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"muxi-thief/api/response"
	"net/http"
)

var ProviderSet = wire.NewSet(NewMiddleware, NewJWTClient)

type ParTokener interface {
	ParseToken(tokenString string) (string, error)
}

type Middleware struct {
	jwt ParTokener
}

func NewMiddleware(jwt ParTokener) *Middleware {
	return &Middleware{jwt}
}

// AuthMiddleware 从请求头中获取认证信息并解析出 user_id
func (m *Middleware) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取 Authorization 请求头
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, response.Err{Err: errors.New("验证失败!：Authorization header is empty").Error()})
			c.Abort()
			return
		}

		//解析jwt
		code, err := m.jwt.ParseToken(authHeader)
		if err != nil || code == "" {
			c.JSON(http.StatusUnauthorized, response.Err{Err: err.Error()})
			c.Abort()
			return
		}

		// 将 code 存储到上下文中
		c.Set("code", code)

		// 继续处理请求
		c.Next()
	}
}
