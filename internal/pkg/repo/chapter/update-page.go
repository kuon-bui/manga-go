package chapterrepo

import (
	"context"
	"manga-go/internal/pkg/model"
)

func (r *ChapterRepo) UpdateChapterPages(ctx context.Context, chapterSlug string, pages []*model.Page) error {
	var chapter model.Chapter
	if err := r.DB.Model(&chapter).Where("slug = ?", chapterSlug).Association("Pages").Replace(pages); err != nil {
		return err
	}

	return nil

}
