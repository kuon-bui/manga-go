package rolerepo

import (
	"manga-go/internal/pkg/model"
	"manga-go/internal/pkg/repo/base"

	"gorm.io/gorm"
)

type RoleRepo struct {
	*base.BaseRepository[model.Role]
}

func NewRoleRepo(db *gorm.DB) *RoleRepo {
	return &RoleRepo{
		BaseRepository: &base.BaseRepository[model.Role]{
			DB: db,
		},
	}
}
