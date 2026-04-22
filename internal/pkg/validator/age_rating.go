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

	allowedRatings := map[constant.ComicAgeRating]bool{
		constant.AgeRatingAll:    true,
		constant.AgeRating13Plus: true,
		constant.AgeRating16Plus: true,
		constant.AgeRating18Plus: true,
	}

	return allowedRatings[ageRating]
}
