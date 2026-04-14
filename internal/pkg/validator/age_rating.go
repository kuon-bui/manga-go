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

	switch ageRating {
	case constant.AgeRatingAll, constant.AgeRating13Plus, constant.AgeRating16Plus, constant.AgeRating18Plus:
		return true
	default:
		return false
	}
}
