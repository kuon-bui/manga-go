package comicrepo

import (
	"manga-go/internal/pkg/model"
	"manga-go/internal/pkg/redis"
	"manga-go/internal/pkg/repo/base"

	"gorm.io/gorm"
)

type ComicRepo struct {
	*base.BaseRepository[model.Comic]
	rds *redis.Redis
}

func NewComicRepo(db *gorm.DB, rds *redis.Redis) *ComicRepo {
	return &ComicRepo{
		BaseRepository: &base.BaseRepository[model.Comic]{
			DB: db,
		},
		rds: rds,
	}
}
