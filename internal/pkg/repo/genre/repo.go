package genrerepo

import (
	"manga-go/internal/pkg/model"
	"manga-go/internal/pkg/redis"
	"manga-go/internal/pkg/repo/base"

	"gorm.io/gorm"
)

type GenreRepo struct {
	*base.BaseRepository[model.Genre]
	rds *redis.Redis
}

func NewGenreRepo(db *gorm.DB, rds *redis.Redis) *GenreRepo {
	return &GenreRepo{
		BaseRepository: &base.BaseRepository[model.Genre]{
			DB: db,
		},
		rds: rds,
	}
}
