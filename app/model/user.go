package model

type QueryUser struct {
	QueryPagination
	ID      *string `form:"id"`
	Keyword *string `form:"keyword"`
}

type CreateUserRequest struct {
	Email       string `json:"email"        binding:"required,email"`
	Username    string `json:"username"     binding:"required,min=1"`
	Password    string `json:"password"     binding:"required,password"`
	DisplayName string `json:"display_name" binding:"required"`
	FirstName   string `json:"first_name"   binding:"required"`
	LastName    string `json:"last_name"    binding:"required"`
}

type UpdateUserRequest struct {
	DisplayName *string `json:"display_name"`
	FirstName   *string `json:"first_name"`
	LastName    *string `json:"last_name"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required,password"`
	NewPassword string `json:"new_password" binding:"required,password"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
	Ticket   string `json:"ticket"`
	Password string `json:"password" binding:"required,password"`
}

type UserResponse struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	VerifyFlag  bool   `json:"verify_flag"`
}

type UserPageResponse struct {
	Meta MetaPagination `json:"meta"`
	Data []UserResponse `json:"data"`
}
