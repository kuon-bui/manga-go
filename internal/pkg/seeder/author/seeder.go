package authorseeder

import (
	"errors"
	"manga-go/internal/pkg/model"
	authorrepo "manga-go/internal/pkg/repo/author"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

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
	repo *authorrepo.AuthorRepo
}

func NewAuthorSeeder(repo *authorrepo.AuthorRepo) *AuthorSeeder {
	return &AuthorSeeder{repo: repo}
}

func (s *AuthorSeeder) Name() string {
	return "AuthorSeeder"
}

func (s *AuthorSeeder) Seed(tx *gorm.DB) error {
	for _, name := range authors {
		_, err := s.repo.FindOneWithTransaction(tx, []any{clause.Eq{Column: "name", Value: name}}, nil)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			author := &model.Author{Name: name}
			if err := s.repo.CreateWithTransaction(tx, author); err != nil {
				return err
			}
		}
	}
	return nil
}
