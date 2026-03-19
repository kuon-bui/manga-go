package userrepo

import (
	"manga-go/internal/pkg/model"
	"manga-go/internal/pkg/redis"
	"manga-go/internal/pkg/repo/base"

	"gorm.io/gorm"
)

type UserRepository struct {
	*base.BaseRepository[model.User]
	redis *redis.Redis
}

func NewUserRepository(db *gorm.DB, redis *redis.Redis) *UserRepository {
	return &UserRepository{
		BaseRepository: &base.BaseRepository[model.User]{
			DB: db,
		},
		redis: redis,
	}
}
