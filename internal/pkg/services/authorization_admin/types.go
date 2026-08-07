package authorizationadmin

import (
	"manga-go/internal/pkg/model"

	"github.com/google/uuid"
)

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

type AdminUserSummary struct {
	ID                   uuid.UUID     `json:"id"`
	Name                 string        `json:"name"`
	Email                string        `json:"email"`
	Roles                []RoleSummary `json:"roles"`
	AuthorizationVersion string        `json:"authorizationVersion"`
}

type PagedUsers struct {
	Data  []AdminUserSummary `json:"data"`
	Total int64              `json:"total"`
	Page  int                `json:"page"`
	Limit int                `json:"limit"`
}

type RoleAccessSummary struct {
	ID                   uuid.UUID `json:"id"`
	Name                 string    `json:"name"`
	Description          *string   `json:"description"`
	Permissions          []string  `json:"permissions"`
	AssignedUserCount    int       `json:"assignedUserCount"`
	AuthorizationVersion string    `json:"authorizationVersion"`
}

type PagedAuditLogs struct {
	Data  []*model.AuthorizationAuditLog `json:"data"`
	Total int64                          `json:"total"`
	Page  int                            `json:"page"`
	Limit int                            `json:"limit"`
}
