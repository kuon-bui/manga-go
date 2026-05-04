package authorseeder

import (
	"errors"
	"fmt"
	"manga-go/internal/pkg/model"
	authorrepo "manga-go/internal/pkg/repo/author"
	seederutil "manga-go/internal/pkg/seeder/util"

	"github.com/jaswdr/faker/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const fakeAuthorCount = 12

var authors = []string{
	"Eiichiro Oda",
	"Akira Toriyama",
	"Masashi Kishimoto",
	"Tite Kubo",
	"Hajime Isayama",
	"Kentaro Miura",
	"Naoki Urasawa",
	"Yoshihiro Togashi",
	"Rumiko Takahashi",
	"Hiromu Arakawa",
	"Kohei Horikoshi",
	"Gege Akutami",
	"Koyoharu Gotouge",
	"ONE",
	"Yusuke Murata",
}

type AuthorSeeder struct {
	repo  *authorrepo.AuthorRepo
	faker faker.Faker
}

func NewAuthorSeeder(repo *authorrepo.AuthorRepo, faker faker.Faker) *AuthorSeeder {
	return &AuthorSeeder{repo: repo, faker: faker}
}

func (s *AuthorSeeder) Name() string {
	return "AuthorSeeder"
}

func (s *AuthorSeeder) Truncate(tx *gorm.DB) error {
	return seederutil.TruncateTables(tx, "comic_authors", "comic_artists", "authors")
}

func (s *AuthorSeeder) Seed(tx *gorm.DB) error {
	for _, name := range authors {
		_, err := s.repo.FindOneWithTransaction(tx, []any{clause.Eq{Column: "name", Value: name}}, nil)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			author := &model.Author{}
			author.Fake(s.faker)
			author.Name = name
			if err := s.repo.CreateWithTransaction(tx, author); err != nil {
				return err
			}
		}
	}

	for index := 1; index <= fakeAuthorCount; index++ {
		author := &model.Author{}
		author.Fake(s.faker)
		author.Name = fmt.Sprintf("%s Seed %02d", author.Name, index)

		_, err := s.repo.FindOneWithTransaction(tx, []any{clause.Eq{Column: "name", Value: author.Name}}, nil)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := s.repo.CreateWithTransaction(tx, author); err != nil {
				return err
			}
		}
	}

	return nil
}
