package fileservice

import "context"

func (s *FileService) GetFile(ctx context.Context, filename string) ([]byte, error) {
	cleanedKey, err := sanitizeObjectKey(filename)
	if err != nil {
		return nil, err
	}

	return s.objectStorage.GetFile(ctx, cleanedKey)
}
