package rolerequest

import "github.com/google/uuid"

type AssignPermissionRequest struct {
	PermissionIDs []uuid.UUID `json:"permissionIds" binding:"required,min=1,dive,required"`
}
