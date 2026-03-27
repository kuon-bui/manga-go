package translationgrouprepo

import (
	"manga-go/internal/pkg/model"
	"manga-go/internal/pkg/redis"
	"manga-go/internal/pkg/repo/base"

	"gorm.io/gorm"
)

type TranslationGroupRepo struct {
	*base.BaseRepository[model.TranslationGroup]
	rds *redis.Redis
}

func NewTranslationGroupRepo(db *gorm.DB, rds *redis.Redis) *TranslationGroupRepo {
	return &TranslationGroupRepo{
		BaseRepository: &base.BaseRepository[model.TranslationGroup]{
			DB: db,
		},
		rds: rds,
	}
}
