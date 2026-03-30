package translationgrouprequest

import "github.com/google/uuid"

type GrantPermissionRequest struct {
	MemberID uuid.UUID `json:"memberId" binding:"required"`
}
