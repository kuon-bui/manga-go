package genreseeder

import (
	"errors"
	"fmt"
	"manga-go/internal/pkg/model"
	genrerepo "manga-go/internal/pkg/repo/genre"
	seederutil "manga-go/internal/pkg/seeder/util"

	"github.com/jaswdr/faker/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const fakeGenreCount = 10

type genreSeed struct {
	Name        string
	Slug        string
	Description string
}

var genres = []genreSeed{
	{Name: "Action", Slug: "action", Description: "High-energy stories featuring combat and adventure."},
	{Name: "Adventure", Slug: "adventure", Description: "Stories following characters on epic journeys."},
	{Name: "Comedy", Slug: "comedy", Description: "Humorous stories designed to entertain and amuse."},
	{Name: "Drama", Slug: "drama", Description: "Emotionally driven stories with serious themes."},
	{Name: "Fantasy", Slug: "fantasy", Description: "Stories set in magical or supernatural worlds."},
	{Name: "Horror", Slug: "horror", Description: "Dark stories intended to frighten and unsettle."},
	{Name: "Mystery", Slug: "mystery", Description: "Stories centered on solving puzzles or crimes."},
	{Name: "Romance", Slug: "romance", Description: "Stories focusing on love and relationships."},
	{Name: "Sci-Fi", Slug: "sci-fi", Description: "Stories exploring futuristic science and technology."},
	{Name: "Slice of Life", Slug: "slice-of-life", Description: "Realistic stories depicting everyday events."},
	{Name: "Sports", Slug: "sports", Description: "Stories centered around athletic competition."},
	{Name: "Thriller", Slug: "thriller", Description: "Suspenseful stories with high-stakes tension."},
	{Name: "Supernatural", Slug: "supernatural", Description: "Stories involving paranormal phenomena."},
	{Name: "Psychological", Slug: "psychological", Description: "Stories exploring the complexity of the human mind."},
	{Name: "Historical", Slug: "historical", Description: "Stories set in historical time periods."},
}

type GenreSeeder struct {
	repo  *genrerepo.GenreRepo
	faker faker.Faker
}

func NewGenreSeeder(repo *genrerepo.GenreRepo, faker faker.Faker) *GenreSeeder {
	return &GenreSeeder{repo: repo, faker: faker}
}

func (s *GenreSeeder) Name() string {
	return "GenreSeeder"
}

func (s *GenreSeeder) Truncate(tx *gorm.DB) error {
	return seederutil.TruncateTables(tx, "comic_genres", "genres")
}

func (s *GenreSeeder) Seed(tx *gorm.DB) error {
	for _, g := range genres {
		_, err := s.repo.FindOneWithTransaction(tx, []any{clause.Eq{Column: "slug", Value: g.Slug}}, nil)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			genre := &model.Genre{}
			genre.Fake(s.faker)
			genre.Name = g.Name
			genre.Slug = g.Slug
			genre.Description = g.Description
			if err := s.repo.CreateWithTransaction(tx, genre); err != nil {
				return err
			}
		}
	}

	for index := 1; index <= fakeGenreCount; index++ {
		genre := &model.Genre{}
		genre.Fake(s.faker)
		genre.Name = fmt.Sprintf("%s Seed %02d", genre.Name, index)
		genre.Slug = fmt.Sprintf("seed-genre-%02d-%s", index, genre.Slug)

		_, err := s.repo.FindOneWithTransaction(tx, []any{clause.Eq{Column: "slug", Value: genre.Slug}}, nil)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := s.repo.CreateWithTransaction(tx, genre); err != nil {
				return err
			}
		}
	}

	return nil
}
