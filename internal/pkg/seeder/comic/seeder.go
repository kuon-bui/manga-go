package comicseeder

import (
	"context"
	"errors"
	"manga-go/internal/pkg/constant"
	"manga-go/internal/pkg/model"
	authorrepo "manga-go/internal/pkg/repo/author"
	chapterrepo "manga-go/internal/pkg/repo/chapter"
	comicrepo "manga-go/internal/pkg/repo/comic"
	genrerepo "manga-go/internal/pkg/repo/genre"
	pagerepo "manga-go/internal/pkg/repo/page"
	tagrepo "manga-go/internal/pkg/repo/tag"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type pageSeed struct {
	PageNumber int
	ImageURL   string
}

type chapterSeed struct {
	Number string
	Title  string
	Pages  []pageSeed
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
			{
				Number: "1",
				Title:  "Romance Dawn",
				Pages: []pageSeed{
					{PageNumber: 1, ImageURL: "https://picsum.photos/seed/op-1-1/800/1200"},
					{PageNumber: 2, ImageURL: "https://picsum.photos/seed/op-1-2/800/1200"},
					{PageNumber: 3, ImageURL: "https://picsum.photos/seed/op-1-3/800/1200"},
					{PageNumber: 4, ImageURL: "https://picsum.photos/seed/op-1-4/800/1200"},
					{PageNumber: 5, ImageURL: "https://picsum.photos/seed/op-1-5/800/1200"},
				},
			},
			{
				Number: "2",
				Title:  "They Call Him Straw Hat Luffy",
				Pages: []pageSeed{
					{PageNumber: 1, ImageURL: "https://picsum.photos/seed/op-2-1/800/1200"},
					{PageNumber: 2, ImageURL: "https://picsum.photos/seed/op-2-2/800/1200"},
					{PageNumber: 3, ImageURL: "https://picsum.photos/seed/op-2-3/800/1200"},
					{PageNumber: 4, ImageURL: "https://picsum.photos/seed/op-2-4/800/1200"},
					{PageNumber: 5, ImageURL: "https://picsum.photos/seed/op-2-5/800/1200"},
				},
			},
			{
				Number: "3",
				Title:  "Enter Zoro: Pirate Hunter",
				Pages: []pageSeed{
					{PageNumber: 1, ImageURL: "https://picsum.photos/seed/op-3-1/800/1200"},
					{PageNumber: 2, ImageURL: "https://picsum.photos/seed/op-3-2/800/1200"},
					{PageNumber: 3, ImageURL: "https://picsum.photos/seed/op-3-3/800/1200"},
					{PageNumber: 4, ImageURL: "https://picsum.photos/seed/op-3-4/800/1200"},
					{PageNumber: 5, ImageURL: "https://picsum.photos/seed/op-3-5/800/1200"},
				},
			},
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
			{
				Number: "1",
				Title:  "Uzumaki Naruto",
				Pages: []pageSeed{
					{PageNumber: 1, ImageURL: "https://picsum.photos/seed/na-1-1/800/1200"},
					{PageNumber: 2, ImageURL: "https://picsum.photos/seed/na-1-2/800/1200"},
					{PageNumber: 3, ImageURL: "https://picsum.photos/seed/na-1-3/800/1200"},
					{PageNumber: 4, ImageURL: "https://picsum.photos/seed/na-1-4/800/1200"},
					{PageNumber: 5, ImageURL: "https://picsum.photos/seed/na-1-5/800/1200"},
				},
			},
			{
				Number: "2",
				Title:  "Konohamaru!",
				Pages: []pageSeed{
					{PageNumber: 1, ImageURL: "https://picsum.photos/seed/na-2-1/800/1200"},
					{PageNumber: 2, ImageURL: "https://picsum.photos/seed/na-2-2/800/1200"},
					{PageNumber: 3, ImageURL: "https://picsum.photos/seed/na-2-3/800/1200"},
					{PageNumber: 4, ImageURL: "https://picsum.photos/seed/na-2-4/800/1200"},
					{PageNumber: 5, ImageURL: "https://picsum.photos/seed/na-2-5/800/1200"},
				},
			},
			{
				Number: "3",
				Title:  "Sasuke Uchiha",
				Pages: []pageSeed{
					{PageNumber: 1, ImageURL: "https://picsum.photos/seed/na-3-1/800/1200"},
					{PageNumber: 2, ImageURL: "https://picsum.photos/seed/na-3-2/800/1200"},
					{PageNumber: 3, ImageURL: "https://picsum.photos/seed/na-3-3/800/1200"},
					{PageNumber: 4, ImageURL: "https://picsum.photos/seed/na-3-4/800/1200"},
					{PageNumber: 5, ImageURL: "https://picsum.photos/seed/na-3-5/800/1200"},
				},
			},
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
			{
				Number: "1",
				Title:  "To You, 2000 Years From Now",
				Pages: []pageSeed{
					{PageNumber: 1, ImageURL: "https://picsum.photos/seed/aot-1-1/800/1200"},
					{PageNumber: 2, ImageURL: "https://picsum.photos/seed/aot-1-2/800/1200"},
					{PageNumber: 3, ImageURL: "https://picsum.photos/seed/aot-1-3/800/1200"},
					{PageNumber: 4, ImageURL: "https://picsum.photos/seed/aot-1-4/800/1200"},
					{PageNumber: 5, ImageURL: "https://picsum.photos/seed/aot-1-5/800/1200"},
				},
			},
			{
				Number: "2",
				Title:  "That Day",
				Pages: []pageSeed{
					{PageNumber: 1, ImageURL: "https://picsum.photos/seed/aot-2-1/800/1200"},
					{PageNumber: 2, ImageURL: "https://picsum.photos/seed/aot-2-2/800/1200"},
					{PageNumber: 3, ImageURL: "https://picsum.photos/seed/aot-2-3/800/1200"},
					{PageNumber: 4, ImageURL: "https://picsum.photos/seed/aot-2-4/800/1200"},
					{PageNumber: 5, ImageURL: "https://picsum.photos/seed/aot-2-5/800/1200"},
				},
			},
			{
				Number: "3",
				Title:  "Night of the Disbanding Ceremony",
				Pages: []pageSeed{
					{PageNumber: 1, ImageURL: "https://picsum.photos/seed/aot-3-1/800/1200"},
					{PageNumber: 2, ImageURL: "https://picsum.photos/seed/aot-3-2/800/1200"},
					{PageNumber: 3, ImageURL: "https://picsum.photos/seed/aot-3-3/800/1200"},
					{PageNumber: 4, ImageURL: "https://picsum.photos/seed/aot-3-4/800/1200"},
					{PageNumber: 5, ImageURL: "https://picsum.photos/seed/aot-3-5/800/1200"},
				},
			},
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
			{
				Number: "1",
				Title:  "The Two Alchemists",
				Pages: []pageSeed{
					{PageNumber: 1, ImageURL: "https://picsum.photos/seed/fma-1-1/800/1200"},
					{PageNumber: 2, ImageURL: "https://picsum.photos/seed/fma-1-2/800/1200"},
					{PageNumber: 3, ImageURL: "https://picsum.photos/seed/fma-1-3/800/1200"},
					{PageNumber: 4, ImageURL: "https://picsum.photos/seed/fma-1-4/800/1200"},
					{PageNumber: 5, ImageURL: "https://picsum.photos/seed/fma-1-5/800/1200"},
				},
			},
			{
				Number: "2",
				Title:  "Body of the Sanctioned",
				Pages: []pageSeed{
					{PageNumber: 1, ImageURL: "https://picsum.photos/seed/fma-2-1/800/1200"},
					{PageNumber: 2, ImageURL: "https://picsum.photos/seed/fma-2-2/800/1200"},
					{PageNumber: 3, ImageURL: "https://picsum.photos/seed/fma-2-3/800/1200"},
					{PageNumber: 4, ImageURL: "https://picsum.photos/seed/fma-2-4/800/1200"},
					{PageNumber: 5, ImageURL: "https://picsum.photos/seed/fma-2-5/800/1200"},
				},
			},
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
			{
				Number: "1",
				Title:  "Ryomen Sukuna",
				Pages: []pageSeed{
					{PageNumber: 1, ImageURL: "https://picsum.photos/seed/jjk-1-1/800/1200"},
					{PageNumber: 2, ImageURL: "https://picsum.photos/seed/jjk-1-2/800/1200"},
					{PageNumber: 3, ImageURL: "https://picsum.photos/seed/jjk-1-3/800/1200"},
					{PageNumber: 4, ImageURL: "https://picsum.photos/seed/jjk-1-4/800/1200"},
					{PageNumber: 5, ImageURL: "https://picsum.photos/seed/jjk-1-5/800/1200"},
				},
			},
			{
				Number: "2",
				Title:  "For Myself",
				Pages: []pageSeed{
					{PageNumber: 1, ImageURL: "https://picsum.photos/seed/jjk-2-1/800/1200"},
					{PageNumber: 2, ImageURL: "https://picsum.photos/seed/jjk-2-2/800/1200"},
					{PageNumber: 3, ImageURL: "https://picsum.photos/seed/jjk-2-3/800/1200"},
					{PageNumber: 4, ImageURL: "https://picsum.photos/seed/jjk-2-4/800/1200"},
					{PageNumber: 5, ImageURL: "https://picsum.photos/seed/jjk-2-5/800/1200"},
				},
			},
			{
				Number: "3",
				Title:  "Girl of Steel",
				Pages: []pageSeed{
					{PageNumber: 1, ImageURL: "https://picsum.photos/seed/jjk-3-1/800/1200"},
					{PageNumber: 2, ImageURL: "https://picsum.photos/seed/jjk-3-2/800/1200"},
					{PageNumber: 3, ImageURL: "https://picsum.photos/seed/jjk-3-3/800/1200"},
					{PageNumber: 4, ImageURL: "https://picsum.photos/seed/jjk-3-4/800/1200"},
					{PageNumber: 5, ImageURL: "https://picsum.photos/seed/jjk-3-5/800/1200"},
				},
			},
		},
	},
}

type ComicSeeder struct {
	comicRepo   *comicrepo.ComicRepo
	authorRepo  *authorrepo.AuthorRepo
	genreRepo   *genrerepo.GenreRepo
	tagRepo     *tagrepo.TagRepo
	chapterRepo *chapterrepo.ChapterRepo
	pageRepo    *pagerepo.PageRepo
}

func NewComicSeeder(
	comicRepo *comicrepo.ComicRepo,
	authorRepo *authorrepo.AuthorRepo,
	genreRepo *genrerepo.GenreRepo,
	tagRepo *tagrepo.TagRepo,
	chapterRepo *chapterrepo.ChapterRepo,
	pageRepo *pagerepo.PageRepo,
) *ComicSeeder {
	return &ComicSeeder{
		comicRepo:   comicRepo,
		authorRepo:  authorRepo,
		genreRepo:   genreRepo,
		tagRepo:     tagRepo,
		chapterRepo: chapterRepo,
		pageRepo:    pageRepo,
	}
}

func (s *ComicSeeder) Name() string {
	return "ComicSeeder"
}

func (s *ComicSeeder) Seed(ctx context.Context) error {
	for _, cs := range comics {
		// Find or create comic by slug
		comic, err := s.comicRepo.FindOne(ctx, []any{clause.Eq{Column: "slug", Value: cs.Slug}}, nil)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			description := cs.Description
			comic = &model.Comic{
				Title:       cs.Title,
				Slug:        cs.Slug,
				Description: &description,
				Type:        cs.Type,
				Status:      cs.Status,
				AgeRating:   cs.AgeRating,
				IsPublished: cs.IsPublished,
				IsHot:       cs.IsHot,
			}
			if err := s.comicRepo.Create(ctx, comic); err != nil {
				return err
			}
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
		if err := s.comicRepo.UpdateComicWithTransaction(ctx, comic.ID, map[string]any{}, map[string]any{
			"Authors": authors,
			"Genres":  genres,
			"Tags":    tags,
		}); err != nil {
			return err
		}

		// Seed chapters and their pages
		for _, ch := range cs.Chapters {
			chapter, err := s.chapterRepo.FindOne(ctx, []any{
				clause.Eq{Column: "comic_id", Value: comic.ID},
				clause.Eq{Column: "number", Value: ch.Number},
			}, nil)
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
			if errors.Is(err, gorm.ErrRecordNotFound) {
				chapter = &model.Chapter{
					ComicID:     comic.ID,
					Number:      ch.Number,
					Title:       ch.Title,
					Slug:        comic.Slug + "-chapter-" + ch.Number,
					IsPublished: true,
				}
				if err := s.chapterRepo.Create(ctx, chapter); err != nil {
					return err
				}
			}

			// Seed pages for this chapter
			for _, pg := range ch.Pages {
				_, err := s.pageRepo.FindOne(ctx, []any{
					clause.Eq{Column: "chapter_id", Value: chapter.ID},
					clause.Eq{Column: "page_number", Value: pg.PageNumber},
				}, nil)
				if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
					return err
				}
				if errors.Is(err, gorm.ErrRecordNotFound) {
					page := &model.Page{
						ChapterID:  chapter.ID,
						PageNumber: pg.PageNumber,
						ImageURL:   pg.ImageURL,
					}
					if err := s.pageRepo.Create(ctx, page); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func (s *ComicSeeder) lookupAuthors(ctx context.Context, names []string) ([]*model.Author, error) {
	result := make([]*model.Author, 0, len(names))
	for _, name := range names {
		a, err := s.authorRepo.FindOne(ctx, []any{clause.Eq{Column: "name", Value: name}}, nil)
		if err != nil {
			return nil, err
		}
		result = append(result, a)
	}
	return result, nil
}

func (s *ComicSeeder) lookupGenres(ctx context.Context, slugs []string) ([]*model.Genre, error) {
	result := make([]*model.Genre, 0, len(slugs))
	for _, slug := range slugs {
		g, err := s.genreRepo.FindOne(ctx, []any{clause.Eq{Column: "slug", Value: slug}}, nil)
		if err != nil {
			return nil, err
		}
		result = append(result, g)
	}
	return result, nil
}

func (s *ComicSeeder) lookupTags(ctx context.Context, slugs []string) ([]*model.Tag, error) {
	result := make([]*model.Tag, 0, len(slugs))
	for _, slug := range slugs {
		t, err := s.tagRepo.FindOne(ctx, []any{clause.Eq{Column: "slug", Value: slug}}, nil)
		if err != nil {
			return nil, err
		}
		result = append(result, t)
	}
	return result, nil
}
