package comicseeder

import (
	"errors"
	"fmt"
	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/constant"
	"manga-go/internal/pkg/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	minChapterPageCount = 15
	maxChapterPageCount = 30
)

func (s *ComicSeeder) upsertSeedPage(
	tx *gorm.DB,
	comicType constant.ComicType,
	chapter *model.Chapter,
	pageNumber int,
	baseImageURL string,
	useFake bool,
	allowUpdate bool,
) error {
	pageType, imageURL, content := s.resolveSeedPageContent(comicType, chapter.Slug, chapter.Title, pageNumber, baseImageURL)

	page, err := s.pageRepo.FindOneWithTransaction(tx, []any{
		clause.Eq{Column: "chapter_id", Value: chapter.ID},
		clause.Eq{Column: "page_number", Value: pageNumber},
	}, nil)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		page = &model.Page{ChapterID: chapter.ID, PageNumber: pageNumber}
		if useFake {
			page.Fake(s.faker)
		}
		page.PageType = pageType
		page.ImageURL = imageURL
		page.Content = content
		return s.pageRepo.CreateWithTransaction(tx, page)
	}

	if !allowUpdate {
		return nil
	}

	return s.pageRepo.UpdateWithTransaction(tx, []any{clause.Eq{Column: "id", Value: page.ID}}, map[string]any{
		"page_type": pageType,
		"image_url": imageURL,
		"content":   content,
	})
}

func (s *ComicSeeder) resolveSeedPageContent(
	comicType constant.ComicType,
	chapterSlug string,
	chapterTitle string,
	pageNumber int,
	baseImageURL string,
) (common.ContentType, string, string) {
	if comicType == constant.ComicTypeNovel {
		return common.ContentTypeText, "", s.buildNovelPageContent(chapterTitle, pageNumber)
	}

	imageURL := baseImageURL
	if comicType == constant.ComicTypeComic || imageURL == "" {
		imageURL = buildInternetImageURL(chapterSlug, pageNumber)
	}

	return common.ContentTypeImage, imageURL, ""
}

func buildInternetImageURL(chapterSlug string, pageNumber int) string {
	variants := []string{"900/1400", "960/1440", "1024/1536", "1080/1600"}
	variant := variants[(len(chapterSlug)+pageNumber)%len(variants)]
	return fmt.Sprintf("https://picsum.photos/seed/%s-%02d/%s", chapterSlug, pageNumber, variant)
}

func resolveChapterPageCount(chapterSlug string) int {
	span := maxChapterPageCount - minChapterPageCount + 1
	checksum := 0
	for _, char := range chapterSlug {
		checksum += int(char)
	}

	return minChapterPageCount + checksum%span
}

func findSeedImageURL(pages []pageSeed, pageNumber int) string {
	for _, page := range pages {
		if page.PageNumber == pageNumber {
			return page.ImageURL
		}
	}

	return ""
}

func (s *ComicSeeder) buildNovelPageContent(chapterTitle string, pageNumber int) string {
	paragraph := s.faker.Lorem().Paragraph(3)
	return fmt.Sprintf("%s\n\nPage %d\n\n%s", chapterTitle, pageNumber, paragraph)
}
