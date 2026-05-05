package validator

import (
	notificationpkg "manga-go/internal/pkg/notification"

	"github.com/go-playground/validator/v10"
)

var ValidateNotificationCategory validator.Func = func(fl validator.FieldLevel) bool {
	category := notificationpkg.Category(fl.Field().String())

	if category == "" {
		return true // Allow empty value, use "required" tag to enforce presence
	}

	return notificationpkg.GetAllowedCategories()[category]
}
