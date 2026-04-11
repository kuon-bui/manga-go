package comicservice

import (
	"manga-go/internal/pkg/model"
	"strings"

	"gorm.io/gorm"
)

func (s *ComicService) resolveOrCreateAuthorsByNames(tx *gorm.DB, names []string) ([]*model.Author, error) {
	normalizedNames := make([]string, 0, len(names))
	seen := make(map[string]struct{}, len(names))

	for _, name := range names {
		normalized := strings.TrimSpace(name)
		if normalized == "" {
			continue
		}
		if _, ok := seen[normalized]; ok {
			continue
		}

		seen[normalized] = struct{}{}
		normalizedNames = append(normalizedNames, normalized)
	}

	if len(normalizedNames) == 0 {
		return []*model.Author{}, nil
	}

	foundAuthors, err := s.authorRepo.FindByNamesWithTx(tx, normalizedNames, nil)
	if err != nil {
		return nil, err
	}

	authorByName := make(map[string]*model.Author, len(foundAuthors))
	for _, author := range foundAuthors {
		if _, exists := authorByName[author.Name]; !exists {
			authorByName[author.Name] = author
		}
	}

	resolvedAuthors := make([]*model.Author, 0, len(normalizedNames))
	for _, name := range normalizedNames {
		author, exists := authorByName[name]
		if exists {
			resolvedAuthors = append(resolvedAuthors, author)
			continue
		}

		author = &model.Author{Name: name}
		if err := s.authorRepo.CreateWithTransaction(tx, author); err != nil {
			return nil, err
		}

		authorByName[name] = author
		resolvedAuthors = append(resolvedAuthors, author)
	}

	return resolvedAuthors, nil
}
