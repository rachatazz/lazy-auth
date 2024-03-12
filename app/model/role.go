package model

type QueryRole struct {
	QueryPagination
	RoleType *string `form:"role_type"`
}

type CreateRoleRequest struct {
	Name        string `json:"name"        binding:"required"`
	Description string `json:"description" binding:"required"`
}

type RoleResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type RolePageResponse struct {
	Meta MetaPagination `json:"meta"`
	Data []RoleResponse `json:"data"`
}
