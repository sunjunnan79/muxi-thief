package middleware

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type JWTClient struct {
	secretKey string
	timeout   time.Duration
}

func NewJWTClient() *JWTClient {
	return &JWTClient{
		secretKey: "muxi-thief-2024",
		timeout:   300 * time.Minute,
	}
}

// GenerateToken 生成 ParTokener token
func (c *JWTClient) GenerateToken(code string) (string, error) {
	if code == "" {
		return "", errors.New("没有code")
	}
	// 设置过期时间
	expirationTime := time.Now().Add(c.timeout)

	// 创建 token
	claims := &jwt.StandardClaims{
		Subject:   code,
		ExpiresAt: expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签署 token
	return token.SignedString([]byte(c.secretKey))

}

// ParseToken 解析 ParTokener token 并返回 userID
func (c *JWTClient) ParseToken(tokenString string) (string, error) {

	claims := &jwt.StandardClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.NewValidationError("unexpected signing method", jwt.ValidationErrorMalformed)
		}
		return []byte(c.secretKey), nil
	})

	if err != nil || !token.Valid {
		return "", err
	}

	return claims.Subject, nil // 返回 code
}
