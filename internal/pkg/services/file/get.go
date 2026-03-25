package fileservice

import "context"

func (s *FileService) GetFile(ctx context.Context, filename string) ([]byte, error) {
	return s.objectStorage.GetFile(ctx, filename)
}
