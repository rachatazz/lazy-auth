package middleware

import (
	"lazy-auth/app/errs"
	"lazy-auth/app/handler"
	"lazy-auth/config"

	"github.com/gin-gonic/gin"
)

type secretGuard struct {
	config config.ConfigEnv
}

type SecretGuard interface {
	ValidateSecret() gin.HandlerFunc
}

func NewSecretGuard(config config.ConfigEnv) SecretGuard {
	return secretGuard{config: config}
}

func (r secretGuard) ValidateSecret() gin.HandlerFunc {
	return func(c *gin.Context) {
		secretKey := c.GetHeader("x-secret-key")
		if r.config.AdminSecret != secretKey {
			handler.HandleError(c, errs.NewUnauthorizedError("secret is invalid"))
			return
		}
		c.Next()
	}
}
