package handler

import (
	"lazy-auth/app/model"
	"lazy-auth/app/repository"
	"lazy-auth/app/service"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) userHandler {
	return userHandler{userService: userService}
}

func (h userHandler) CreateUserAdmin(c *gin.Context) {
	var body model.CreateUserRequest
	err := ValidationPipe(c, &body, ValidateBody)
	if err != nil {
		HandleError(c, err)
		return
	}

	user, err := h.userService.CreateUser(body, "admin")
	if err != nil {
		HandleError(c, err)
		return
	}

	HandleOk(c, user, nil)
}

func (h userHandler) CreateUser(c *gin.Context) {
	var body model.CreateUserRequest
	err := ValidationPipe(c, &body, ValidateBody)
	if err != nil {
		HandleError(c, err)
		return
	}

	user, err := h.userService.CreateUser(body, "user")
	if err != nil {
		HandleError(c, err)
		return
	}

	HandleOk(c, user, nil)
}

func (h userHandler) GetUsers(c *gin.Context) {
	var query model.QueryUser
	err := ValidationPipe(c, &query, ValidateQuery)
	if err != nil {
		HandleError(c, err)
		return
	}

	userResponse, err := h.userService.GetUsers(query)
	if err != nil {
		HandleError(c, err)
		return
	}

	HandleOk(c, userResponse.Data, userResponse.Meta)
}

func (h userHandler) GetMe(c *gin.Context) {
	session, _ := c.Get("session")
	user, err := h.userService.GetUserById(session.(*repository.Session).UserID)
	if err != nil {
		HandleError(c, err)
		return
	}

	HandleOk(c, user, nil)
}

func (h userHandler) UpdateMe(c *gin.Context) {
	session, _ := c.Get("session")
	var body model.UpdateUserRequest
	err := ValidationPipe(c, &body, ValidateBody)
	if err != nil {
		HandleError(c, err)
		return
	}
	h.userService.UpdateUserById(session.(*repository.Session).UserID, body)
	HandleOk(c, nil, nil)
}

func (h userHandler) ChangePassword(c *gin.Context) {
	session, _ := c.Get("session")

	var body model.ChangePasswordRequest
	err := ValidationPipe(c, &body, ValidateBody)
	if err != nil {
		HandleError(c, err)
		return
	}

	_, err = h.userService.ChangePassword(session.(*repository.Session).UserID, body)
	if err != nil {
		HandleError(c, err)
		return
	}

	HandleOk(c, nil, nil)
}
