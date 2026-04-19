package fileservice

import (
	chapterrepo "manga-go/internal/pkg/repo/chapter"
	comicrepo "manga-go/internal/pkg/repo/comic"
	objectstorage "manga-go/internal/pkg/object_storage"

	"go.uber.org/fx"
)

type FileService struct {
	objectStorage *objectstorage.ObjectStorage
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
