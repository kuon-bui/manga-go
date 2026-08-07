package userrequest

import (
	"manga-go/internal/pkg/common"

	"github.com/google/uuid"
)

type ListAuthorizationUsersRequest struct {
	common.Paging
	Search string     `form:"search" binding:"omitempty,max=255"`
	RoleID *uuid.UUID `form:"role_id"`
}
