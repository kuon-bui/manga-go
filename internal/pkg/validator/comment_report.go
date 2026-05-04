package validator

import (
	"manga-go/internal/pkg/constant"

	"github.com/go-playground/validator/v10"
)

var ValidateCommentReportReason validator.Func = func(fl validator.FieldLevel) bool {
	reason := fl.Field().String()
	if reason == "" {
		return true // Allow empty value, use "required" tag to enforce presence
	}

	allowedReasons := constant.GetAllowedCommentReportReasons()

	return allowedReasons[constant.CommentReportReason(reason)]
}
