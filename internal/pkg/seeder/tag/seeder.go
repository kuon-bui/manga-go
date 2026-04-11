package tagseeder

import (
	"context"
	"manga-go/internal/pkg/model"
	tagrepo "manga-go/internal/pkg/repo/tag"

	"gorm.io/gorm/clause"
)

type tagSeed struct {
	Name string
	Slug string
}

var tags = []tagSeed{
	{Name: "Isekai", Slug: "isekai"},
	{Name: "Reincarnation", Slug: "reincarnation"},
	{Name: "Magic", Slug: "magic"},
	{Name: "Sword Art", Slug: "sword-art"},
	{Name: "School Life", Slug: "school-life"},
	{Name: "Overpowered Protagonist", Slug: "overpowered-protagonist"},
	{Name: "Time Travel", Slug: "time-travel"},
	{Name: "Demons", Slug: "demons"},
	{Name: "Gods", Slug: "gods"},
	{Name: "Martial Arts", Slug: "martial-arts"},
	{Name: "Harem", Slug: "harem"},
	{Name: "Cooking", Slug: "cooking"},
	{Name: "Music", Slug: "music"},
	{Name: "Military", Slug: "military"},
	{Name: "Post-Apocalyptic", Slug: "post-apocalyptic"},
	{Name: "Mecha", Slug: "mecha"},
	{Name: "Vampire", Slug: "vampire"},
	{Name: "Zombie", Slug: "zombie"},
	{Name: "System", Slug: "system"},
	{Name: "Dungeons", Slug: "dungeons"},
}

type TagSeeder struct {
	repo *tagrepo.TagRepo
}

func NewTagSeeder(repo *tagrepo.TagRepo) *TagSeeder {
	return &TagSeeder{repo: repo}
}

func (s *TagSeeder) Name() string {
	return "TagSeeder"
}

func (s *TagSeeder) Seed(ctx context.Context) error {
	for _, t := range tags {
		tag := model.Tag{
			Name: t.Name,
			Slug: t.Slug,
		}
		if err := s.repo.DB.WithContext(ctx).
			Where(clause.Eq{Column: "slug", Value: t.Slug}).
			FirstOrCreate(&tag).Error; err != nil {
			return err
		}
	}
	return nil
}
