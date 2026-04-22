package validator

import (
	"manga-go/internal/pkg/constant"

	"github.com/go-playground/validator/v10"
)

var ValidateComicStatus validator.Func = func(fl validator.FieldLevel) bool {
	status := constant.ComicStatus(fl.Field().String())

	if status == "" {
		return true // Allow empty value, use "required" tag to enforce presence
	}

	allowedStatuses := map[constant.ComicStatus]bool{
		constant.ComicStatusOngoing:   true,
		constant.ComicStatusCompleted: true,
		constant.ComicStatusHiatus:    true,
		constant.ComicStatusCancelled: true,
	}

	return allowedStatuses[status]
}
