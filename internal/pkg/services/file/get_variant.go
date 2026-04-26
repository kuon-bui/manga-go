package fileservice

import "context"

func (s *FileService) GetFileByVariant(ctx context.Context, filename, rawVariant string) ([]byte, string, error) {
	cleanedKey, err := sanitizeObjectKey(filename)
	if err != nil {
		return nil, "", err
	}

	variant, err := ParseImageVariant(rawVariant)
	if err != nil {
		return nil, "", err
	}

	variantKey := BuildVariantObjectKey(cleanedKey, variant)
	fileContent, err := s.objectStorage.GetFile(ctx, variantKey)
	if err == nil {
		return fileContent, variantKey, nil
	}

	if variantKey != cleanedKey && s.objectStorage.IsNotFoundError(err) {
		fallbackContent, fallbackErr := s.objectStorage.GetFile(ctx, cleanedKey)
		if fallbackErr == nil {
			return fallbackContent, cleanedKey, nil
		}
		return nil, "", fallbackErr
	}

	return nil, "", err
}
