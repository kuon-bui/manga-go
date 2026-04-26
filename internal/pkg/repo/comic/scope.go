package comicrepo

import (
	comicrequest "manga-go/internal/pkg/request/comic"

	"gorm.io/gorm"
)

func applyComicFilters(filters *comicrequest.ListComicsRequest) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if filters.TranslationGroupSlug != "" {
			db = db.
				Joins("JOIN translation_groups ON translation_groups.id = comics.translation_group_id").
				Where("translation_groups.slug = ?", filters.TranslationGroupSlug)
		}

		if filters.Status != "" {
			db = db.Where("comics.status = ?", filters.Status)
		}

		if filters.Search != "" {
			searchPattern := "%" + filters.Search + "%"
			db = db.Where(
				`(
				comics.title ILIKE ?
				OR EXISTS (
					SELECT 1
					FROM jsonb_array_elements_text(COALESCE(comics.alternative_titles, '[]'::jsonb)) AS alt(title)
					WHERE alt.title ILIKE ?
				)
			)`,
				searchPattern,
				searchPattern,
			)
		}

		if len(filters.GenreSlugs) > 0 {
			db = db.
				Joins("JOIN comic_genres ON comic_genres.comic_id = comics.id").
				Joins("JOIN genres ON genres.id = comic_genres.genre_id").
				Where("genres.slug IN ?", filters.GenreSlugs)
		}

		if len(filters.TagSlugs) > 0 {
			db = db.
				Joins("JOIN comic_tags ON comic_tags.comic_id = comics.id").
				Joins("JOIN tags ON tags.id = comic_tags.tag_id").
				Where("tags.slug IN ?", filters.TagSlugs)
		}

		return db
	}
}
