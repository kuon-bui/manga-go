package comicseeder

import (
	"errors"
	"fmt"
	"manga-go/internal/pkg/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *ComicSeeder) seedFakeComics(tx *gorm.DB, users []*model.User, translationGroups []*model.TranslationGroup) error {
	users = filterSeededUsers(users)
	translationGroups = filterSeededTranslationGroups(translationGroups)

	authors, err := s.authorRepo.FindAllWithTx(tx, []any{func(db *gorm.DB) *gorm.DB {
		return db.Order("name ASC")
	}}, nil)
	if err != nil {
		return err
	}
	genres, err := s.genreRepo.FindAllWithTx(tx, []any{func(db *gorm.DB) *gorm.DB {
		return db.Order("slug ASC")
	}}, nil)
	if err != nil {
		return err
	}
	tags, err := s.tagRepo.FindAllWithTx(tx, []any{func(db *gorm.DB) *gorm.DB {
		return db.Order("slug ASC")
	}}, nil)
	if err != nil {
		return err
	}

	for index := 1; index <= fakeComicCount; index++ {
		comic := &model.Comic{}
		comic.Fake(s.faker)
		comic.Title = fmt.Sprintf("%s Seed %02d", comic.Title, index)
		comic.Slug = fmt.Sprintf("seed-comic-%02d-%s", index, comic.Slug)
		comic.IsPublished = true
		comic.IsHot = index%2 == 0
		comic.IsFeatured = index%3 == 0

		if len(users) > 0 {
			uploadedByID := users[(index-1)%len(users)].ID
			comic.UploadedByID = &uploadedByID
		}
		if len(translationGroups) > 0 {
			translationGroupID := translationGroups[(index-1)%len(translationGroups)].ID
			comic.TranslationGroupID = &translationGroupID
		}

		payload := map[string]any{
			"title":                comic.Title,
			"description":          comic.Description,
			"type":                 comic.Type,
			"status":               comic.Status,
			"age_rating":           comic.AgeRating,
			"is_published":         comic.IsPublished,
			"is_hot":               comic.IsHot,
			"is_featured":          comic.IsFeatured,
			"published_year":       comic.PublishedYear,
			"translation_group_id": comic.TranslationGroupID,
			"uploaded_by_id":       comic.UploadedByID,
		}

		existing, err := s.comicRepo.FindOneWithTransaction(tx, []any{clause.Eq{Column: "slug", Value: comic.Slug}}, nil)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := s.comicRepo.CreateWithTransaction(tx, comic); err != nil {
				return err
			}
		} else {
			comic = existing
			if err := s.comicRepo.UpdateWithTransaction(tx, []any{clause.Eq{Column: "id", Value: comic.ID}}, payload); err != nil {
				return err
			}
		}

		if err := s.assignFakeAssociations(tx, comic, authors, genres, tags, index); err != nil {
			return err
		}
		if err := s.seedFakeChapters(tx, comic, users, index); err != nil {
			return err
		}
	}

	return nil
}

func (s *ComicSeeder) assignFakeAssociations(tx *gorm.DB, comic *model.Comic, authors []*model.Author, genres []*model.Genre, tags []*model.Tag, index int) error {
	selectedAuthors := takeAuthors(authors, index, 2)
	selectedArtists := takeAuthors(authors, index+1, 2)
	selectedGenres := takeGenres(genres, index, 3)
	selectedTags := takeTags(tags, index, 3)

	return s.comicRepo.UpdateComicWithTransaction(tx, comic.ID, map[string]any{}, map[string]any{
		"Authors": selectedAuthors,
		"Artists": selectedArtists,
		"Genres":  selectedGenres,
		"Tags":    selectedTags,
	})
}

func (s *ComicSeeder) seedFakeChapters(tx *gorm.DB, comic *model.Comic, users []*model.User, index int) error {
	chapterCount := 2 + index%3
	for chapterIndex := 1; chapterIndex <= chapterCount; chapterIndex++ {
		slug := fmt.Sprintf("%s-seed-ch-%02d", comic.Slug, chapterIndex)
		publishedAt := resolveChapterPublishedAt(comic.Slug, chapterIndex)
		chapter, err := s.chapterRepo.FindOneWithTransaction(tx, []any{clause.Eq{Column: "slug", Value: slug}}, nil)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			chapter = &model.Chapter{ComicID: comic.ID}
			chapter.Fake(s.faker)
			chapter.Number = fmt.Sprintf("%d", chapterIndex)
			chapter.Slug = slug
			chapter.IsPublished = true
			chapter.PublishedAt = publishedAt
			if len(users) > 0 {
				uploadedByID := users[(index+chapterIndex-1)%len(users)].ID
				chapter.UploadedByID = &uploadedByID
			}
			if err := s.chapterRepo.CreateWithTransaction(tx, chapter); err != nil {
				return err
			}
		} else if err := s.chapterRepo.UpdateWithTransaction(tx, []any{clause.Eq{Column: "id", Value: chapter.ID}}, map[string]any{
			"published_at": publishedAt,
		}); err != nil {
			return err
		}

		pageCount := resolveChapterPageCount(chapter.Slug)
		for pageIndex := 1; pageIndex <= pageCount; pageIndex++ {
			if err := s.upsertSeedPage(tx, comic.Type, chapter, pageIndex, "", true, true); err != nil {
				return err
			}
		}
	}

	return nil
}

func takeAuthors(authors []*model.Author, offset, count int) []*model.Author {
	if len(authors) == 0 {
		return []*model.Author{}
	}
	result := make([]*model.Author, 0, count)
	for index := 0; index < count; index++ {
		result = append(result, authors[(offset+index)%len(authors)])
	}
	return result
}

func takeGenres(genres []*model.Genre, offset, count int) []*model.Genre {
	if len(genres) == 0 {
		return []*model.Genre{}
	}
	result := make([]*model.Genre, 0, count)
	for index := range count {
		result = append(result, genres[(offset+index)%len(genres)])
	}
	return result
}

func takeTags(tags []*model.Tag, offset, count int) []*model.Tag {
	if len(tags) == 0 {
		return []*model.Tag{}
	}
	result := make([]*model.Tag, 0, count)
	for index := range count {
		result = append(result, tags[(offset+index)%len(tags)])
	}
	return result
}
