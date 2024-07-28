package middleware

import (
	"lazy-auth/app/errs"
	"lazy-auth/app/handler"
	"lazy-auth/app/repository"

	"github.com/gin-gonic/gin"
)

type roleGuard struct {
	userRepository repository.UserRepository
}

type RoleGuard interface {
	ValidateRole(role ...string) gin.HandlerFunc
}

func NewRoleGuard(
	userRepository repository.UserRepository,
) RoleGuard {
	return roleGuard{userRepository: userRepository}
}

func (r roleGuard) ValidateRole(role ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		session, _ := c.Get("session")
		user, _ := r.userRepository.GetById(session.(*repository.Session).UserID)
		for _, v := range role {
			if user.Role.Name == v {
				c.Next()
				return
			}
		}
		handler.HandleError(c, errs.NewForbiddenError("forbidden"))
	}
}
