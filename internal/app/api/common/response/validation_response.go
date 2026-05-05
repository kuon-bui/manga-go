package response

import (
	"fmt"
	"manga-go/internal/pkg/constant"
	notificationpkg "manga-go/internal/pkg/notification"

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
		return constant.ComicAgeRatingValidationMessage(validationErr.Field())
	case "comic_type":
		return constant.ComicTypeValidationMessage(validationErr.Field())
	case "comic_status":
		return constant.ComicStatusValidationMessage(validationErr.Field())
	case "follow_status":
		return constant.FollowStatusValidationMessage(validationErr.Field())
	case "notification_category":
		return notificationpkg.CategoryValidationMessage(validationErr.Field())
	case "notification_type":
		return notificationpkg.TypeValidationMessage(validationErr.Field())
	case "notification_entity_type":
		return notificationpkg.EntityTypeValidationMessage(validationErr.Field())
	default:
		return fmt.Sprintf("%s is invalid (%s)", validationErr.Field(), validationErr.Tag())
	}
}
