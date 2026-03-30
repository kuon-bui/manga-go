package translationgrouprequest

import "github.com/google/uuid"

type KickMemberRequest struct {
	MemberID uuid.UUID `json:"memberId" binding:"required"`
}
