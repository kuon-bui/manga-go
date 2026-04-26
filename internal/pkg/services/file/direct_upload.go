package fileservice

import (
	"context"
	"io"
)

func (s *FileService) UploadFile(ctx context.Context, filename string, body io.Reader, contentLength int64, contentType string) error {
	return s.objectStorage.UploadFile(ctx, filename, body, contentLength, contentType)
}

func (s *FileService) DeleteFile(ctx context.Context, filename string) error {
	cleanedKey, err := sanitizeObjectKey(filename)
	if err != nil {
		return err
	}

	return s.objectStorage.DeleteFile(ctx, cleanedKey)
}

func (s *FileService) IsNotFoundError(err error) bool {
	return s.objectStorage.IsNotFoundError(err)
}
