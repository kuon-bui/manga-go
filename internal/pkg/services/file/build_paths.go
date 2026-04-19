package fileservice

import (
	"context"
	"errors"
	"gorm.io/gorm/clause"

	"github.com/google/uuid"
)

// BuildChapterImagePath resolves IDs to slugs and builds the S3 path
// Path format: comics/{comicSlug}/chapters/{chapterSlug}/pages/{uniqueFilename}
func (s *FileService) BuildChapterImagePath(ctx context.Context, comicIdStr, chapterIdStr, uniqueFilename string) (string, error) {
	comicId, err := uuid.Parse(comicIdStr)
	if err != nil {
		return "", errors.New("invalid comicId format")
	}

	chapterId, err := uuid.Parse(chapterIdStr)
	if err != nil {
		return "", errors.New("invalid chapterId format")
	}

	// Verify comic exists and get slug
	comic, err := s.comicRepo.FindOne(ctx, []any{
		clause.Eq{Column: "id", Value: comicId},
	}, nil)
	if err != nil {
		return "", errors.New("comic not found")
	}

	// Verify chapter exists and get slug
	chapter, err := s.chapterRepo.FindOne(ctx, []any{
		clause.Eq{Column: "id", Value: chapterId},
		clause.Eq{Column: "comic_id", Value: comicId},
	}, nil)
	if err != nil {
		return "", errors.New("chapter not found or doesn't belong to this comic")
	}

	// Build path
	path := "comics/" + comic.Slug + "/chapters/" + chapter.Slug + "/pages/" + uniqueFilename
	return path, nil
}

// BuildTempChapterImagePath builds path for chapter images without chapterId
// Used when creating new chapter - images uploaded before chapter is created
// Path format: comics/{comicSlug}/temp-uploads/{uniqueFilename}
func (s *FileService) BuildTempChapterImagePath(ctx context.Context, comicIdStr, uniqueFilename string) (string, error) {
	comicId, err := uuid.Parse(comicIdStr)
	if err != nil {
		return "", errors.New("invalid comicId format")
	}

	// Verify comic exists and get slug
	comic, err := s.comicRepo.FindOne(ctx, []any{
		clause.Eq{Column: "id", Value: comicId},
	}, nil)
	if err != nil {
		return "", errors.New("comic not found")
	}

	// Build path to temp folder
	path := "comics/" + comic.Slug + "/temp-uploads/" + uniqueFilename
	return path, nil
}

// BuildCoverImagePath resolves comic ID to slug and builds the S3 path
// Path format: comics/{comicSlug}/cover/{uniqueFilename}
func (s *FileService) BuildCoverImagePath(ctx context.Context, comicIdStr, uniqueFilename string) (string, error) {
	comicId, err := uuid.Parse(comicIdStr)
	if err != nil {
		return "", errors.New("invalid comicId format")
	}

	// Verify comic exists and get slug
	comic, err := s.comicRepo.FindOne(ctx, []any{
		clause.Eq{Column: "id", Value: comicId},
	}, nil)
	if err != nil {
		return "", errors.New("comic not found")
	}

	// Build path
	path := "comics/" + comic.Slug + "/cover/" + uniqueFilename
	return path, nil
}
