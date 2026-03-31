package permissionrepo

import (
	"manga-go/internal/pkg/model"
	"manga-go/internal/pkg/repo/base"

	"gorm.io/gorm"
)

type PermissionRepo struct {
	*base.BaseRepository[model.Permission]
}

func NewPermissionRepo(db *gorm.DB) *PermissionRepo {
	return &PermissionRepo{
		BaseRepository: &base.BaseRepository[model.Permission]{
			DB: db,
		},
	}
}
