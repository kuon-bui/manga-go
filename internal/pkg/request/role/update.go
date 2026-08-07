package rolerequest

type UpdateRoleRequest struct {
	Name        string  `json:"name" binding:"required,min=2,max=100"`
	Description *string `json:"description" binding:"omitempty,max=1000"`
}
