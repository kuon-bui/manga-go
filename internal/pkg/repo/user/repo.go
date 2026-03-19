package userrepo

import (
	"base-go/internal/pkg/model"
	"base-go/internal/pkg/redis"
	"base-go/internal/pkg/repo/base"

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
