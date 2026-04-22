package validator

import "github.com/go-playground/validator/v10"

var ValidateOrderDirection validator.Func = func(fl validator.FieldLevel) bool {
	order := fl.Field().String()

	if order == "" {
		return true // Allow empty value, use "required" tag to enforce presence
	}

	allowedOrders := map[string]bool{
		"asc":  true,
		"desc": true,
		"ASC":  true,
		"DESC": true,
	}

	return allowedOrders[order]
}

var ValidateComicSortBy validator.Func = func(fl validator.FieldLevel) bool {
	sortBy := fl.Field().String()

	if sortBy == "" {
		return true // Allow empty value, use "required" tag to enforce presence
	}

	allowedSortBy := map[string]bool{
		"lastChapterAt": true,
		"createdAt":     true,
		"rating":        true,
		"followCount":   true,
	}

	return allowedSortBy[sortBy]
}
