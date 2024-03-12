package middleware

import (
	"strings"

	"lazy-auth/app/errs"
	"lazy-auth/app/handler"
	"lazy-auth/app/repository"
	"lazy-auth/common"
	"lazy-auth/config"

	"github.com/gin-gonic/gin"
)

type tokenGuard struct {
	sessionRepository repository.SessionRepository
	config            config.ConfigEnv
}

type TokenGuard interface {
	ValidateToken() gin.HandlerFunc
}

func NewTokenGuard(
	sessionRepository repository.SessionRepository,
	config config.ConfigEnv,
) TokenGuard {
	return tokenGuard{sessionRepository: sessionRepository, config: config}
}

func (r tokenGuard) ValidateToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.GetHeader("Authorization")
		if !strings.HasPrefix(authorization, "Bearer") {
			handler.HandleError(c, errs.NewUnauthorizedError("invalid token"))
			return
		}

		splits := strings.Split(authorization, " ")
		claims, valid := common.ValidateToken(splits[1], r.config.JwtTokenSecret)
		if !valid {
			handler.HandleError(c, errs.NewUnauthorizedError("invalid token"))
			return
		}

		session, err := r.sessionRepository.GetById(claims.Id)
		if err != nil {
			handler.HandleError(c, errs.NewUnauthorizedError("invalid token"))
			return
		}

		c.Set("session", session)
		c.Next()
	}
}
