package translationgrouprequest

import "github.com/google/uuid"

type TransferOwnershipRequest struct {
	NewOwnerID uuid.UUID `json:"newOwnerId" binding:"required"`
}
