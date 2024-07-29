package service

import "lazy-auth/app/model"

type AuthService interface {
	GetRoles(body model.QueryRole) (*model.RolePageResponse, error)
	CreateRole(body model.CreateRoleRequest) (*model.RoleResponse, error)
	Login(body model.LoginRequest) (*model.TokenResponse, error)
	RefreshToken(refreshToken string) (*model.TokenResponse, error)
	Logout(body model.LogoutRequest) error
	ForgotPassword(body model.ForgotPasswordRequest) error
	ResetPassword(body model.ResetPasswordRequest) error
}
