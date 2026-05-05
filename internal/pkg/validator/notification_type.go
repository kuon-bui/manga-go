package validator

import (
	notificationpkg "manga-go/internal/pkg/notification"

	"github.com/go-playground/validator/v10"
)

var ValidateNotificationType validator.Func = func(fl validator.FieldLevel) bool {
	notificationType := notificationpkg.Type(fl.Field().String())

	if notificationType == "" {
		return true // Allow empty value, use "required" tag to enforce presence
	}

	return notificationpkg.GetAllowedTypes()[notificationType]
}
