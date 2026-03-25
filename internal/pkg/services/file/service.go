package fileservice

import (
	objectstorage "manga-go/internal/pkg/object_storage"

	"go.uber.org/fx"
)

type FileService struct {
	objectStorage *objectstorage.ObjectStorage
}

type FileServiceParams struct {
	fx.In
	ObjectStorage *objectstorage.ObjectStorage
}

func NewFileService(params FileServiceParams) *FileService {
	return &FileService{
		objectStorage: params.ObjectStorage,
	}
}
