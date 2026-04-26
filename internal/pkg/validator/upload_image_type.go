package validator

import (
	"manga-go/internal/pkg/constant"
	"strings"

	"github.com/go-playground/validator/v10"
)

var ValidateUploadImageType validator.Func = func(fl validator.FieldLevel) bool {
	imageType := fl.Field().String()

	if imageType == "" {
		return true // Allow empty value, use "required" tag to enforce presence
	}

	normalizedType := strings.ToLower(strings.TrimSpace(imageType))

	allowedTypes := map[string]bool{
		string(constant.UploadImageTypeComic):   true,
		string(constant.UploadImageTypeChapter): true,
		"cover":                                 true,
	}

	return allowedTypes[normalizedType]
}
