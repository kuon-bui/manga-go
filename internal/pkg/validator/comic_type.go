package validator

import (
	"manga-go/internal/pkg/constant"

	"github.com/go-playground/validator/v10"
)

var ValidateComicType validator.Func = func(fl validator.FieldLevel) bool {
	comicType := constant.ComicType(fl.Field().String())

	if comicType == "" {
		return true // Allow empty value, use "required" tag to enforce presence
	}
	allowedTypes := constant.GetAllowedComicTypes()

	return allowedTypes[comicType]
}
