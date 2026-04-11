package comicseeder

import (
	"context"
	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/constant"
	"manga-go/internal/pkg/model"
	authorrepo "manga-go/internal/pkg/repo/author"
	chapterrepo "manga-go/internal/pkg/repo/chapter"
	comicrepo "manga-go/internal/pkg/repo/comic"
	genrerepo "manga-go/internal/pkg/repo/genre"
	tagrepo "manga-go/internal/pkg/repo/tag"

	"gorm.io/gorm/clause"
)

type chapterSeed struct {
	Number string
	Title  string
}

type comicSeed struct {
	Title       string
	Slug        string
	Description string
	Type        constant.ComicType
	Status      constant.ComicStatus
	AgeRating   constant.ComicAgeRating
	IsPublished bool
	IsHot       bool
	Authors     []string
	Genres      []string
	Tags        []string
	Chapters    []chapterSeed
}

var comics = []comicSeed{
	{
		Title:       "One Piece",
		Slug:        "one-piece",
		Description: "A young pirate dreams of becoming the King of the Pirates by finding the legendary One Piece treasure.",
		Type:        constant.ComicTypeManga,
		Status:      constant.ComicStatusOngoing,
		AgeRating:   constant.AgeRatingAll,
		IsPublished: true,
		IsHot:       true,
		Authors:     []string{"Eiichiro Oda"},
		Genres:      []string{"action", "adventure", "comedy"},
		Tags:        []string{"sword-art", "demons"},
		Chapters: []chapterSeed{
			{Number: "1", Title: "Romance Dawn"},
			{Number: "2", Title: "They Call Him Straw Hat Luffy"},
			{Number: "3", Title: "Enter Zoro: Pirate Hunter"},
		},
	},
	{
		Title:       "Naruto",
		Slug:        "naruto",
		Description: "A young ninja with dreams of becoming the greatest ninja and leader of his village.",
		Type:        constant.ComicTypeManga,
		Status:      constant.ComicStatusCompleted,
		AgeRating:   constant.AgeRatingAll,
		IsPublished: true,
		IsHot:       false,
		Authors:     []string{"Masashi Kishimoto"},
		Genres:      []string{"action", "adventure", "fantasy"},
		Tags:        []string{"martial-arts", "demons", "magic"},
		Chapters: []chapterSeed{
			{Number: "1", Title: "Uzumaki Naruto"},
			{Number: "2", Title: "Konohamaru!"},
			{Number: "3", Title: "Sasuke Uchiha"},
		},
	},
	{
		Title:       "Attack on Titan",
		Slug:        "attack-on-titan",
		Description: "Humanity fights for survival against gigantic humanoid creatures called Titans.",
		Type:        constant.ComicTypeManga,
		Status:      constant.ComicStatusCompleted,
		AgeRating:   constant.AgeRating16Plus,
		IsPublished: true,
		IsHot:       true,
		Authors:     []string{"Hajime Isayama"},
		Genres:      []string{"action", "drama", "thriller"},
		Tags:        []string{"military", "post-apocalyptic"},
		Chapters: []chapterSeed{
			{Number: "1", Title: "To You, 2000 Years From Now"},
			{Number: "2", Title: "That Day"},
			{Number: "3", Title: "Night of the Disbanding Ceremony"},
		},
	},
	{
		Title:       "Fullmetal Alchemist",
		Slug:        "fullmetal-alchemist",
		Description: "Two brothers use alchemy to restore their bodies after a failed ritual, uncovering a dark conspiracy.",
		Type:        constant.ComicTypeManga,
		Status:      constant.ComicStatusCompleted,
		AgeRating:   constant.AgeRating13Plus,
		IsPublished: true,
		IsHot:       false,
		Authors:     []string{"Hiromu Arakawa"},
		Genres:      []string{"action", "adventure", "fantasy"},
		Tags:        []string{"magic", "military"},
		Chapters: []chapterSeed{
			{Number: "1", Title: "The Two Alchemists"},
			{Number: "2", Title: "Body of the Sanctioned"},
		},
	},
	{
		Title:       "Jujutsu Kaisen",
		Slug:        "jujutsu-kaisen",
		Description: "A high school student joins a secret organization of jujutsu sorcerers to eliminate cursed spirits.",
		Type:        constant.ComicTypeManga,
		Status:      constant.ComicStatusOngoing,
		AgeRating:   constant.AgeRating16Plus,
		IsPublished: true,
		IsHot:       true,
		Authors:     []string{"Gege Akutami"},
		Genres:      []string{"action", "supernatural", "horror"},
		Tags:        []string{"demons", "martial-arts"},
		Chapters: []chapterSeed{
			{Number: "1", Title: "Ryomen Sukuna"},
			{Number: "2", Title: "For Myself"},
			{Number: "3", Title: "Girl of Steel"},
		},
	},
}

type ComicSeeder struct {
	comicRepo   *comicrepo.ComicRepo
	authorRepo  *authorrepo.AuthorRepo
	genreRepo   *genrerepo.GenreRepo
	tagRepo     *tagrepo.TagRepo
	chapterRepo *chapterrepo.ChapterRepo
}

func NewComicSeeder(
	comicRepo *comicrepo.ComicRepo,
	authorRepo *authorrepo.AuthorRepo,
	genreRepo *genrerepo.GenreRepo,
	tagRepo *tagrepo.TagRepo,
	chapterRepo *chapterrepo.ChapterRepo,
) *ComicSeeder {
	return &ComicSeeder{
		comicRepo:   comicRepo,
		authorRepo:  authorRepo,
		genreRepo:   genreRepo,
		tagRepo:     tagRepo,
		chapterRepo: chapterRepo,
	}
}

func (s *ComicSeeder) Name() string {
	return "ComicSeeder"
}

func (s *ComicSeeder) Seed(ctx context.Context) error {
	for _, cs := range comics {
		desc := cs.Description
		comic := model.Comic{
			Title:       cs.Title,
			Slug:        cs.Slug,
			Description: &desc,
			Type:        cs.Type,
			Status:      cs.Status,
			AgeRating:   cs.AgeRating,
			IsPublished: cs.IsPublished,
			IsHot:       cs.IsHot,
		}
		if err := s.comicRepo.DB.WithContext(ctx).
			Where(clause.Eq{Column: "slug", Value: cs.Slug}).
			FirstOrCreate(&comic).Error; err != nil {
			return err
		}

		authors, err := s.lookupAuthors(ctx, cs.Authors)
		if err != nil {
			return err
		}
		genres, err := s.lookupGenres(ctx, cs.Genres)
		if err != nil {
			return err
		}
		tags, err := s.lookupTags(ctx, cs.Tags)
		if err != nil {
			return err
		}

		// Replace associations (idempotent)
		if err := s.comicRepo.DB.WithContext(ctx).Model(&model.Comic{SqlModel: common.SqlModel{ID: comic.ID}}).
			Association("Authors").Replace(authors); err != nil {
			return err
		}
		if err := s.comicRepo.DB.WithContext(ctx).Model(&model.Comic{SqlModel: common.SqlModel{ID: comic.ID}}).
			Association("Genres").Replace(genres); err != nil {
			return err
		}
		if err := s.comicRepo.DB.WithContext(ctx).Model(&model.Comic{SqlModel: common.SqlModel{ID: comic.ID}}).
			Association("Tags").Replace(tags); err != nil {
			return err
		}

		// Seed chapters
		for _, ch := range cs.Chapters {
			chapter := model.Chapter{
				ComicID:     comic.ID,
				Number:      ch.Number,
				Title:       ch.Title,
				Slug:        comic.Slug + "-chapter-" + ch.Number,
				IsPublished: true,
			}
			if err := s.chapterRepo.DB.WithContext(ctx).
				Where(clause.Eq{Column: "comic_id", Value: comic.ID}).
				Where(clause.Eq{Column: "number", Value: ch.Number}).
				FirstOrCreate(&chapter).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *ComicSeeder) lookupAuthors(ctx context.Context, names []string) ([]*model.Author, error) {
	result := make([]*model.Author, 0, len(names))
	for _, name := range names {
		var a model.Author
		if err := s.authorRepo.DB.WithContext(ctx).
			Where(clause.Eq{Column: "name", Value: name}).
			First(&a).Error; err != nil {
			return nil, err
		}
		result = append(result, &a)
	}
	return result, nil
}

func (s *ComicSeeder) lookupGenres(ctx context.Context, slugs []string) ([]*model.Genre, error) {
	result := make([]*model.Genre, 0, len(slugs))
	for _, slug := range slugs {
		var g model.Genre
		if err := s.genreRepo.DB.WithContext(ctx).
			Where(clause.Eq{Column: "slug", Value: slug}).
			First(&g).Error; err != nil {
			return nil, err
		}
		result = append(result, &g)
	}
	return result, nil
}

func (s *ComicSeeder) lookupTags(ctx context.Context, slugs []string) ([]*model.Tag, error) {
	result := make([]*model.Tag, 0, len(slugs))
	for _, slug := range slugs {
		var t model.Tag
		if err := s.tagRepo.DB.WithContext(ctx).
			Where(clause.Eq{Column: "slug", Value: slug}).
			First(&t).Error; err != nil {
			return nil, err
		}
		result = append(result, &t)
	}
	return result, nil
}
