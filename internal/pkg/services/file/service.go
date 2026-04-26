package fileservice

import (
	"context"
	"io"

	objectstorage "manga-go/internal/pkg/object_storage"
	chapterrepo "manga-go/internal/pkg/repo/chapter"
	comicrepo "manga-go/internal/pkg/repo/comic"

	"go.uber.org/fx"
)

type fileStorage interface {
	CreatePresignedURL(ctx context.Context, key string) (string, error)
	GetFile(ctx context.Context, fileName string) ([]byte, error)
	UploadFile(ctx context.Context, fileName string, body io.Reader, contentLength int64, contentType string) error
	IsNotFoundError(err error) bool
}

type FileService struct {
	objectStorage fileStorage
	comicRepo     *comicrepo.ComicRepo
	chapterRepo   *chapterrepo.ChapterRepo
}

type FileServiceParams struct {
	fx.In
	ObjectStorage *objectstorage.ObjectStorage
	ComicRepo     *comicrepo.ComicRepo
	ChapterRepo   *chapterrepo.ChapterRepo
}

func NewFileService(params FileServiceParams) *FileService {
	return &FileService{
		objectStorage: params.ObjectStorage,
		comicRepo:     params.ComicRepo,
		chapterRepo:   params.ChapterRepo,
	}
}
