package validator

import (
	notificationpkg "manga-go/internal/pkg/notification"

	"github.com/go-playground/validator/v10"
)

var ValidateNotificationEntityType validator.Func = func(fl validator.FieldLevel) bool {
	entityType := notificationpkg.EntityType(fl.Field().String())

	if entityType == "" {
		return true // Allow empty value, use "required" tag to enforce presence
	}

	return notificationpkg.GetAllowedEntityTypes()[entityType]
}
