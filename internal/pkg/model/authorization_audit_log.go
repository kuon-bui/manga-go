package model

import (
	"time"

	"manga-go/internal/pkg/common"

	"github.com/google/uuid"
)

type AuthorizationAuditLog struct {
	ID                 uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	ActorUserID        *uuid.UUID     `json:"actorUserId" gorm:"column:actor_user_id"`
	ActorNameSnapshot  string         `json:"actorName" gorm:"column:actor_name_snapshot"`
	ActorEmailSnapshot string         `json:"actorEmail" gorm:"column:actor_email_snapshot"`
	Action             string         `json:"action" gorm:"column:action"`
	TargetType         string         `json:"targetType" gorm:"column:target_type"`
	TargetID           uuid.UUID      `json:"targetId" gorm:"column:target_id"`
	TargetNameSnapshot string         `json:"targetName" gorm:"column:target_name_snapshot"`
	Before             common.JSONMap `json:"before" gorm:"column:before;type:jsonb"`
	After              common.JSONMap `json:"after" gorm:"column:after;type:jsonb"`
	CreatedAt          time.Time      `json:"createdAt" gorm:"column:created_at"`
}

func (AuthorizationAuditLog) TableName() string {
	return "authorization_audit_logs"
}
