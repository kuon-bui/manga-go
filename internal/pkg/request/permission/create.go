package permissionrequest

type CreatePermissionRequest struct {
	Name string `json:"name" binding:"required"`
}
