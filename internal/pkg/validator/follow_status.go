package validator

import (
	"manga-go/internal/pkg/constant"

	"github.com/go-playground/validator/v10"
)

var ValidateFollowStatus validator.Func = func(fl validator.FieldLevel) bool {
	followStatus := fl.Field().String()

	if followStatus == "" {
		return true // Allow empty value, use "required" tag to enforce presence
	}

	allowedStatuses := map[constant.FollowStatus]bool{
		constant.FollowStatusReading:   true,
		constant.FollowStatusPlanned:   true,
		constant.FollowStatusCompleted: true,
		constant.FollowStatusDropped:   true,
		constant.FollowStatusFavorite:  true,
	}

	if allowedStatuses[constant.FollowStatus(followStatus)] {
		return true
	}

	return false
}
