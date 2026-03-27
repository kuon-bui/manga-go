package tagrepo

import (
	"manga-go/internal/pkg/model"
	"manga-go/internal/pkg/redis"
	"manga-go/internal/pkg/repo/base"

	"gorm.io/gorm"
)

type TagRepo struct {
	*base.BaseRepository[model.Tag]
	rds *redis.Redis
}

func NewTagRepo(db *gorm.DB, rds *redis.Redis) *TagRepo {
	return &TagRepo{
		BaseRepository: &base.BaseRepository[model.Tag]{
			DB: db,
		},
		rds: rds,
	}
}
