package translationgroupseeder

import (
	"errors"
	"fmt"
	"manga-go/internal/pkg/model"
	translationgrouprepo "manga-go/internal/pkg/repo/translation_group"
	userrepo "manga-go/internal/pkg/repo/user"
	seederutil "manga-go/internal/pkg/seeder/util"
	"strings"

	"github.com/jaswdr/faker/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const translationGroupCount = 3

type TranslationGroupSeeder struct {
	repo     *translationgrouprepo.TranslationGroupRepo
	userRepo *userrepo.UserRepository
	faker    faker.Faker
}

func NewTranslationGroupSeeder(
	repo *translationgrouprepo.TranslationGroupRepo,
	userRepo *userrepo.UserRepository,
	faker faker.Faker,
) *TranslationGroupSeeder {
	return &TranslationGroupSeeder{repo: repo, userRepo: userRepo, faker: faker}
}

func (s *TranslationGroupSeeder) Name() string {
	return "TranslationGroupSeeder"
}

func (s *TranslationGroupSeeder) Truncate(tx *gorm.DB) error {
	if err := tx.Exec("UPDATE users SET translation_group_id = NULL WHERE translation_group_id IS NOT NULL").Error; err != nil {
		return err
	}
	if err := tx.Exec("UPDATE comics SET translation_group_id = NULL WHERE translation_group_id IS NOT NULL").Error; err != nil {
		return err
	}
	return seederutil.TruncateTables(tx, "translation_groups")
}

func (s *TranslationGroupSeeder) Seed(tx *gorm.DB) error {
	users, err := s.userRepo.FindAllWithTx(tx, []any{func(db *gorm.DB) *gorm.DB {
		return db.Order("email ASC")
	}}, nil)
	if err != nil {
		return err
	}
	users = filterSeedUsers(users)
	if len(users) == 0 {
		return nil
	}

	groupCount := translationGroupCount
	if len(users) < groupCount {
		groupCount = len(users)
	}

	groups := make([]*model.TranslationGroup, 0, groupCount)
	for index := 1; index <= groupCount; index++ {
		owner := users[index-1]
		slug := fmt.Sprintf("seed-translation-group-%02d", index)
		group, err := s.repo.FindOneWithTransaction(tx, []any{clause.Eq{Column: "slug", Value: slug}}, nil)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			group = &model.TranslationGroup{OwnerID: owner.ID}
			group.Fake(s.faker)
			group.Name = fmt.Sprintf("%s Seed %02d", group.Name, index)
			group.Slug = slug
			logoURL := fmt.Sprintf("https://picsum.photos/seed/%s/300/300", slug)
			group.LogoUrl = &logoURL
			if err := s.repo.CreateWithTransaction(tx, group); err != nil {
				return err
			}
		} else {
			logoURL := fmt.Sprintf("https://picsum.photos/seed/%s/300/300", slug)
			if err := s.repo.UpdateWithTransaction(tx, []any{clause.Eq{Column: "id", Value: group.ID}}, map[string]any{
				"owner_id": owner.ID,
				"logo_url": logoURL,
			}); err != nil {
				return err
			}
		}

		groups = append(groups, group)
		if err := s.userRepo.UpdateWithTransaction(tx, []any{clause.Eq{Column: "id", Value: owner.ID}}, map[string]any{
			"translation_group_id": group.ID,
		}); err != nil {
			return err
		}
	}

	for index, user := range users[groupCount:] {
		group := groups[index%len(groups)]
		if err := s.userRepo.UpdateWithTransaction(tx, []any{clause.Eq{Column: "id", Value: user.ID}}, map[string]any{
			"translation_group_id": group.ID,
		}); err != nil {
			return err
		}
	}

	return nil
}

func filterSeedUsers(users []*model.User) []*model.User {
	filtered := make([]*model.User, 0, len(users))
	for _, user := range users {
		if user != nil && strings.HasPrefix(user.Email, "seed-user-") && strings.HasSuffix(user.Email, "@manga.local") {
			filtered = append(filtered, user)
		}
	}

	return filtered
}
