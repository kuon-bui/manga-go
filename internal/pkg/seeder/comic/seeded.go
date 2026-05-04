package comicseeder

import (
	"manga-go/internal/pkg/model"
	"strings"
)

func filterSeededUsers(users []*model.User) []*model.User {
	filtered := make([]*model.User, 0, len(users))
	for _, user := range users {
		if user != nil && strings.HasPrefix(user.Email, "seed-user-") && strings.HasSuffix(user.Email, "@manga.local") {
			filtered = append(filtered, user)
		}
	}

	return filtered
}

func filterSeededTranslationGroups(groups []*model.TranslationGroup) []*model.TranslationGroup {
	filtered := make([]*model.TranslationGroup, 0, len(groups))
	for _, group := range groups {
		if group != nil && strings.HasPrefix(group.Slug, "seed-translation-group-") {
			filtered = append(filtered, group)
		}
	}

	return filtered
}
