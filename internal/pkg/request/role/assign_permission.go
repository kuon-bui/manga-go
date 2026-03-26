package rolerequest

import "github.com/google/uuid"

type AssignPermissionRequest struct {
	PermissionIDs []uuid.UUID `json:"permission_ids" binding:"required,min=1"`
}
