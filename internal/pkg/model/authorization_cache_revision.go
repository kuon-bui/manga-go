package model

import "time"

type AuthorizationCacheRevision struct {
	Scope     string    `json:"scope" gorm:"column:scope;primaryKey"`
	Version   uint64    `json:"version" gorm:"column:version"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at"`
}

func (AuthorizationCacheRevision) TableName() string {
	return "authorization_cache_revisions"
}
