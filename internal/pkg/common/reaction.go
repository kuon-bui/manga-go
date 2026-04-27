package common

import "strings"

const (
	ReactionTypeLike  = "LIKE"
	ReactionTypeLove  = "LOVE"
	ReactionTypeHaha  = "HAHA"
	ReactionTypeWow   = "WOW"
	ReactionTypeSad   = "SAD"
	ReactionTypeAngry = "ANGRY"
)

var allowedReactionTypes = map[string]struct{}{
	ReactionTypeLike:  {},
	ReactionTypeLove:  {},
	ReactionTypeHaha:  {},
	ReactionTypeWow:   {},
	ReactionTypeSad:   {},
	ReactionTypeAngry: {},
}

func NormalizeReactionType(t string) string {
	return strings.ToUpper(strings.TrimSpace(t))
}

func IsValidReactionType(t string) bool {
	_, ok := allowedReactionTypes[NormalizeReactionType(t)]
	return ok
}
