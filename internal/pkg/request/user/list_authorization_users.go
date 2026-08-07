package userrequest

import (
	"errors"
	"manga-go/internal/pkg/common"

	"github.com/google/uuid"
)

type ListAuthorizationUsersRequest struct {
	common.Paging
	Search string     `form:"search" binding:"omitempty,max=255"`
	RoleID *uuid.UUID `form:"role_id"`
}

func (r *ListAuthorizationUsersRequest) Validate() error {
	r.Fulfill()
	if r.Limit > 100 {
		return errors.New("limit must not exceed 100")
	}
	return nil
}
