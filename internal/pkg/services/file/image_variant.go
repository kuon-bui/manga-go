package fileservice

import (
	"errors"
	"path"
	"strings"
)

type ImageVariant string

const (
	ImageVariantSmall  ImageVariant = "small"
	ImageVariantMedium ImageVariant = "medium"
	ImageVariantLarge  ImageVariant = "large"
	ImageVariantNormal ImageVariant = "normal"
)

type imageVariantPreset struct {
	Width int
}

var imageVariantOrder = []ImageVariant{
	ImageVariantSmall,
	ImageVariantMedium,
	ImageVariantLarge,
	ImageVariantNormal,
}

var imageVariantPresets = map[ImageVariant]imageVariantPreset{
	ImageVariantSmall:  {Width: 480},
	ImageVariantMedium: {Width: 720},
	ImageVariantLarge:  {Width: 1080},
	ImageVariantNormal: {Width: 0},
}

func ParseImageVariant(raw string) (ImageVariant, error) {
	v := ImageVariant(strings.TrimSpace(raw))
	if v == "" {
		return ImageVariantNormal, nil
	}

	switch v {
	case ImageVariantSmall, ImageVariantMedium, ImageVariantLarge, ImageVariantNormal:
		return v, nil
	default:
		return "", errors.New("invalid size")
	}
}

func BuildVariantObjectKey(canonicalObjectKey string, variant ImageVariant) string {
	if variant == ImageVariantNormal {
		return canonicalObjectKey
	}

	ext := path.Ext(canonicalObjectKey)
	base := strings.TrimSuffix(canonicalObjectKey, ext)
	return base + "__" + string(variant) + ".webp"
}

func sanitizeObjectKey(raw string) (string, error) {
	key := strings.TrimSpace(strings.TrimPrefix(raw, "/"))
	if key == "" || strings.Contains(key, "\\") || strings.Contains(key, "\x00") {
		return "", errors.New("invalid filename")
	}

	cleaned := path.Clean(key)
	if cleaned == "." || cleaned == "/" || cleaned == "" {
		return "", errors.New("invalid filename")
	}

	if strings.HasPrefix(cleaned, "../") || cleaned == ".." || strings.HasPrefix(cleaned, "/") {
		return "", errors.New("invalid filename")
	}

	return cleaned, nil
}
