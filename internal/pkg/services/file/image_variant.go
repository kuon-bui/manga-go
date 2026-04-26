package fileservice

import (
	"errors"
	"path"
	"strings"
)

type ImageVariant string

const (
	ImageVariantEconomy ImageVariant = "economy"
	ImageVariantSmall   ImageVariant = "small"
	ImageVariantClear   ImageVariant = "clear"
	ImageVariantSharp   ImageVariant = "sharp"
)

type imageVariantPreset struct {
	Width int
}

var imageVariantOrder = []ImageVariant{
	ImageVariantEconomy,
	ImageVariantSmall,
	ImageVariantClear,
	ImageVariantSharp,
}

var imageVariantPresets = map[ImageVariant]imageVariantPreset{
	ImageVariantEconomy: {Width: 480},
	ImageVariantSmall:   {Width: 720},
	ImageVariantClear:   {Width: 1080},
	ImageVariantSharp:   {Width: 0},
}

func ParseImageVariant(raw string) (ImageVariant, error) {
	v := ImageVariant(strings.TrimSpace(raw))
	if v == "" {
		return ImageVariantSharp, nil
	}

	switch v {
	case ImageVariantEconomy, ImageVariantSmall, ImageVariantClear, ImageVariantSharp:
		return v, nil
	default:
		return "", errors.New("invalid variant")
	}
}

func BuildVariantObjectKey(canonicalObjectKey string, variant ImageVariant) string {
	if variant == ImageVariantSharp {
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
