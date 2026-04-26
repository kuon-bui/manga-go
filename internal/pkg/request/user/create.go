package userrequest

type CreateUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Avatar   string `json:"avatar"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}
