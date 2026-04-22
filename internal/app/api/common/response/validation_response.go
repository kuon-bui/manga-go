package response

import (
	"fmt"
	"manga-go/internal/pkg/constant"

	"github.com/go-playground/validator/v10"
)

func parseValidationErrors(err error) []ValidationFieldError {
	if err == nil {
		return nil
	}

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return nil
	}

	result := make([]ValidationFieldError, 0, len(validationErrors))
	for _, validationErr := range validationErrors {
		result = append(result, ValidationFieldError{
			Field:   validationErr.Field(),
			Message: buildValidationMessage(validationErr),
		})
	}

	return result
}

func buildValidationMessage(validationErr validator.FieldError) string {
	switch validationErr.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", validationErr.Field())
	case "email":
		return fmt.Sprintf("%s must be a valid email", validationErr.Field())
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", validationErr.Field(), validationErr.Param())
	case "min":
		return fmt.Sprintf("%s must be at least %s", validationErr.Field(), validationErr.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s", validationErr.Field(), validationErr.Param())
	case "len":
		return fmt.Sprintf("%s must have length %s", validationErr.Field(), validationErr.Param())
	case "uuid", "uuid4":
		return fmt.Sprintf("%s must be a valid UUID", validationErr.Field())
	case "age_rating":
		return fmt.Sprintf(
			"%s must be a valid age rating (%s, %s,%s, %s)",
			validationErr.Field(),
			constant.AgeRatingAll,
			constant.AgeRating13Plus,
			constant.AgeRating16Plus,
			constant.AgeRating18Plus,
		)
	case "comic_type":
		return fmt.Sprintf(
			"%s must be a valid comic type (%s, %s, %s, %s, %s)",
			validationErr.Field(),
			constant.ComicTypeManga,
			constant.ComicTypeManhwa,
			constant.ComicTypeManhua,
			constant.ComicTypeComic,
			constant.ComicTypeNovel,
		)
	case "comic_status":
		return fmt.Sprintf(
			"%s must be a valid comic status (%s, %s, %s, %s)",
			validationErr.Field(),
			constant.ComicStatusOngoing,
			constant.ComicStatusCompleted,
			constant.ComicStatusHiatus,
			constant.ComicStatusCancelled,
		)
	case "order_check":
		return fmt.Sprintf(
			"%s must be a valid order direction (ASC or DESC)",
			validationErr.Field(),
		)
	case "comic_sort_by":
		return fmt.Sprintf(
			"%s must be a valid sort by field (lastChapterAt, createdAt, rating, followCount)",
			validationErr.Field(),
		)
	default:
		return fmt.Sprintf("%s is invalid (%s)", validationErr.Field(), validationErr.Tag())
	}
}
