package fileservice

import (
	"bytes"
	"context"
	"errors"
	"io"
	"path"
	"strings"
)

const webpContentType = "image/webp"

type UploadedImageVariant struct {
	Variant     string `json:"variant"`
	Width       int    `json:"width"`
	Path        string `json:"path"`
	URL         string `json:"url"`
	Size        int64  `json:"size"`
	ContentType string `json:"content_type"`
}

type UploadImageVariantsResult struct {
	Filename    string                 `json:"filename"`
	Path        string                 `json:"path"`
	URL         string                 `json:"url"`
	Size        int64                  `json:"size"`
	ContentType string                 `json:"content_type"`
	Variants    []UploadedImageVariant `json:"variants"`
}

func (s *FileService) UploadImageVariants(ctx context.Context, canonicalObjectKey string, body io.Reader) (*UploadImageVariantsResult, error) {
	if body == nil {
		return nil, errors.New("empty file")
	}

	cleanedKey, err := sanitizeObjectKey(canonicalObjectKey)
	if err != nil {
		return nil, err
	}

	raw, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}

	decodedImage, err := decodeImage(raw)
	if err != nil {
		return nil, err
	}

	variants := make([]UploadedImageVariant, 0, len(imageVariantOrder))
	var canonicalSize int64

	for _, variant := range imageVariantOrder {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		preset := imageVariantPresets[variant]
		encoded, width, err := encodeWebPVariant(decodedImage, preset)
		if err != nil {
			return nil, err
		}

		variantObjectKey := BuildVariantObjectKey(cleanedKey, variant)
		if err := s.objectStorage.UploadFile(ctx, variantObjectKey, bytes.NewReader(encoded), int64(len(encoded)), webpContentType); err != nil {
			return nil, err
		}

		if variant == ImageVariantNormal {
			canonicalSize = int64(len(encoded))
		}

		variants = append(variants, UploadedImageVariant{
			Variant:     string(variant),
			Width:       width,
			Path:        variantObjectKey,
			URL:         "/files/content/" + cleanedKey,
			Size:        int64(len(encoded)),
			ContentType: webpContentType,
		})
	}

	filename := path.Base(cleanedKey)
	if strings.TrimSpace(filename) == "" {
		filename = cleanedKey
	}

	return &UploadImageVariantsResult{
		Filename:    filename,
		Path:        cleanedKey,
		URL:         "/files/content/" + cleanedKey,
		Size:        canonicalSize,
		ContentType: webpContentType,
		Variants:    variants,
	}, nil
}
