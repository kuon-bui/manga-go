package validator

import (
	"manga-go/internal/pkg/constant"

	"github.com/go-playground/validator/v10"
)

var ValidateAgeRating validator.Func = func(fl validator.FieldLevel) bool {
	ageRating := constant.ComicAgeRating(fl.Field().String())

	if ageRating == "" {
		return true // Allow empty value, use "required" tag to enforce presence
	}

	allowedRatings := constant.GetAllowedComicAgeRatings()

	return allowedRatings[ageRating]
}
