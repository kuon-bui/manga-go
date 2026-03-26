package fileservice

import "context"

func (s *FileService) GeneratePresignedURL(ctx context.Context, filename string) (string, error) {
	url, err := s.objectStorage.CreatePresignedURL(ctx, filename)
	if err != nil {
		return "", err
	}
	return url, nil
}
