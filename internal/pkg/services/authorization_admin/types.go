package authorizationadmin

import "github.com/google/uuid"

type RoleSummary struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
}

type AuthorizationProfile struct {
	UserID      uuid.UUID     `json:"userId"`
	Roles       []RoleSummary `json:"roles"`
	Permissions []string      `json:"permissions"`
	Version     string        `json:"version"`
}
