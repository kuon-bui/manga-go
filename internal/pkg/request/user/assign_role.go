package userrequest

import "github.com/google/uuid"

type AssignRoleRequest struct {
	RoleIDs []uuid.UUID `json:"role_ids" binding:"required,min=1"`
}
