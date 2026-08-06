package comicstatrepo

import (
	"manga-go/internal/pkg/model"
	"manga-go/internal/pkg/repo/base"

	"gorm.io/gorm"
)

type ComicStatRepo struct {
	*base.BaseRepository[model.ComicStat]
}

func NewComicStatRepo(db *gorm.DB) *ComicStatRepo {
	return &ComicStatRepo{BaseRepository: &base.BaseRepository[model.ComicStat]{DB: db}}
}
