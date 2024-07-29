package handler

import (
	"lazy-auth/app/model"
	"lazy-auth/app/repository"
	"lazy-auth/app/service"

	"github.com/gin-gonic/gin"
)

type authHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) authHandler {
	return authHandler{authService: authService}
}

func (h authHandler) CreateRole(c *gin.Context) {
	var body model.CreateRoleRequest
	err := ValidationPipe(c, &body, ValidateBody)
	if err != nil {
		HandleError(c, err)
		return
	}

	role, err := h.authService.CreateRole(body)
	if err != nil {
		HandleError(c, err)
		return
	}

	HandleOk(c, role, nil)
}

func (h authHandler) GetRoles(c *gin.Context) {
	var query model.QueryRole
	err := ValidationPipe(c, &query, ValidateQuery)
	if err != nil {
		HandleError(c, err)
		return
	}

	roleResponse, err := h.authService.GetRoles(query)
	if err != nil {
		HandleError(c, err)
		return
	}

	HandleOk(c, roleResponse.Data, roleResponse.Meta)
}

func (h authHandler) Login(c *gin.Context) {
	var body model.LoginRequest
	err := ValidationPipe(c, &body, ValidateBody)
	if err != nil {
		HandleError(c, err)
		return
	}
	token, err := h.authService.Login(body)
	if err != nil {
		HandleError(c, err)
		return
	}

	HandleOk(c, token, nil)
}

func (h authHandler) RefreshToken(c *gin.Context) {
	var body model.RefreshTokenRequest
	err := ValidationPipe(c, &body, ValidateBody)
	if err != nil {
		HandleError(c, err)
		return
	}

	token, err := h.authService.RefreshToken(body.RefreshToken)
	if err != nil {
		HandleError(c, err)
		return
	}

	HandleOk(c, token, nil)
}

func (h authHandler) Logout(c *gin.Context) {
	var body model.LogoutRequest
	err := ValidationPipe(c, &body, ValidateBody)
	if err != nil {
		HandleError(c, err)
		return
	}

	err = h.authService.Logout(body)
	if err != nil {
		HandleError(c, err)
		return
	}

	HandleOk(c, nil, nil)
}

func (h authHandler) ChangePassword(c *gin.Context) {
	session, _ := c.Get("session")

	var body model.ChangePasswordRequest
	err := ValidationPipe(c, &body, ValidateBody)
	if err != nil {
		HandleError(c, err)
		return
	}

	err = h.authService.ChangePassword(session.(*repository.Session).UserID, body)
	if err != nil {
		HandleError(c, err)
		return
	}

	HandleOk(c, nil, nil)
}

func (h authHandler) ForgotPassword(c *gin.Context) {
	var body model.ForgotPasswordRequest
	err := ValidationPipe(c, &body, ValidateBody)
	if err != nil {
		HandleError(c, err)
		return
	}

	err = h.authService.ForgotPassword(body)
	if err != nil {
		HandleError(c, err)
		return
	}

	HandleOk(c, nil, nil)
}

func (h authHandler) ResetPassword(c *gin.Context) {
	var body model.ResetPasswordRequest
	err := ValidationPipe(c, &body, ValidateBody)
	if err != nil {
		HandleError(c, err)
		return
	}

	err = h.authService.ResetPassword(body)
	if err != nil {
		HandleError(c, err)
		return
	}

	HandleOk(c, nil, nil)
}
