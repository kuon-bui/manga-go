package fileservice

import (
	"context"
	"io"
)

func (s *FileService) UploadFile(ctx context.Context, filename string, body io.Reader, contentLength int64, contentType string) error {
	return s.objectStorage.UploadFile(ctx, filename, body, contentLength, contentType)
}
