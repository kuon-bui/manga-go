package rolerequest

type UpdateRoleRequest struct {
	Name string `json:"name" binding:"required"`
}
