package permissionrequest

type UpdatePermissionRequest struct {
	Name string `json:"name" binding:"required"`
}
