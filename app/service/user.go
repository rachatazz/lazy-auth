package service

import "lazy-auth/app/model"

type UserService interface {
	CreateUser(body model.CreateUserRequest, role string) (*model.UserResponse, error)
	GetUsers(query model.QueryUser) (*model.UserPageResponse, error)
	GetUserById(id string) (*model.UserResponse, error)
	UpdateUserById(id string, body model.UpdateUserRequest) (*model.UserResponse, error)
}
