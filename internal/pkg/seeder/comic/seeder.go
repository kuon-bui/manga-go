package comicseeder

import (
	"errors"
	"manga-go/internal/pkg/constant"
	"manga-go/internal/pkg/model"
	authorrepo "manga-go/internal/pkg/repo/author"
	chapterrepo "manga-go/internal/pkg/repo/chapter"
	comicrepo "manga-go/internal/pkg/repo/comic"
	genrerepo "manga-go/internal/pkg/repo/genre"
	pagerepo "manga-go/internal/pkg/repo/page"
	tagrepo "manga-go/internal/pkg/repo/tag"
	translationgrouprepo "manga-go/internal/pkg/repo/translation_group"
	userrepo "manga-go/internal/pkg/repo/user"
	seederutil "manga-go/internal/pkg/seeder/util"

	"github.com/jaswdr/faker/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const fakeComicCount = 8

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
	Artists     []string
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
		Artists:     []string{"Eiichiro Oda"},
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
		Artists:     []string{"Masashi Kishimoto"},
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
		Artists:     []string{"Hajime Isayama"},
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
		Artists:     []string{"Hiromu Arakawa"},
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
		Artists:     []string{"Gege Akutami"},
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
	{
		Title:       "The Wandering Archivist",
		Slug:        "the-wandering-archivist",
		Description: "A young archivist travels across ruined kingdoms to recover forbidden books and uncover the truth behind a vanished empire.",
		Type:        constant.ComicTypeNovel,
		Status:      constant.ComicStatusOngoing,
		AgeRating:   constant.AgeRating13Plus,
		IsPublished: true,
		IsHot:       false,
		Authors:     []string{"Naoki Urasawa"},
		Artists:     []string{"Naoki Urasawa"},
		Genres:      []string{"fantasy", "mystery", "drama"},
		Tags:        []string{"magic", "time-travel"},
		Chapters: []chapterSeed{
			{
				Number: "1",
				Title:  "The Last Library",
				Pages: []pageSeed{
					{PageNumber: 1, ImageURL: "https://picsum.photos/seed/novel-1-1/800/1200"},
					{PageNumber: 2, ImageURL: "https://picsum.photos/seed/novel-1-2/800/1200"},
					{PageNumber: 3, ImageURL: "https://picsum.photos/seed/novel-1-3/800/1200"},
				},
			},
			{
				Number: "2",
				Title:  "Map of Ashes",
				Pages: []pageSeed{
					{PageNumber: 1, ImageURL: "https://picsum.photos/seed/novel-2-1/800/1200"},
					{PageNumber: 2, ImageURL: "https://picsum.photos/seed/novel-2-2/800/1200"},
					{PageNumber: 3, ImageURL: "https://picsum.photos/seed/novel-2-3/800/1200"},
				},
			},
		},
	},
}

type ComicSeeder struct {
	comicRepo            *comicrepo.ComicRepo
	authorRepo           *authorrepo.AuthorRepo
	genreRepo            *genrerepo.GenreRepo
	tagRepo              *tagrepo.TagRepo
	chapterRepo          *chapterrepo.ChapterRepo
	pageRepo             *pagerepo.PageRepo
	userRepo             *userrepo.UserRepository
	translationGroupRepo *translationgrouprepo.TranslationGroupRepo
	faker                faker.Faker
}

func NewComicSeeder(
	comicRepo *comicrepo.ComicRepo,
	authorRepo *authorrepo.AuthorRepo,
	genreRepo *genrerepo.GenreRepo,
	tagRepo *tagrepo.TagRepo,
	chapterRepo *chapterrepo.ChapterRepo,
	pageRepo *pagerepo.PageRepo,
	userRepo *userrepo.UserRepository,
	translationGroupRepo *translationgrouprepo.TranslationGroupRepo,
	faker faker.Faker,
) *ComicSeeder {
	return &ComicSeeder{
		comicRepo:            comicRepo,
		authorRepo:           authorRepo,
		genreRepo:            genreRepo,
		tagRepo:              tagRepo,
		chapterRepo:          chapterRepo,
		pageRepo:             pageRepo,
		userRepo:             userRepo,
		translationGroupRepo: translationGroupRepo,
		faker:                faker,
	}
}

func (s *ComicSeeder) Name() string {
	return "ComicSeeder"
}

func (s *ComicSeeder) Truncate(tx *gorm.DB) error {
	return seederutil.TruncateTables(tx, "pages", "chapters", "comic_artists", "comic_authors", "comic_genres", "comic_tags", "comics")
}

func (s *ComicSeeder) Seed(tx *gorm.DB) error {
	users, err := s.userRepo.FindAllWithTx(tx, []any{func(db *gorm.DB) *gorm.DB {
		return db.Order("email ASC")
	}}, nil)
	if err != nil {
		return err
	}
	translationGroups, err := s.translationGroupRepo.FindAllWithTx(tx, []any{func(db *gorm.DB) *gorm.DB {
		return db.Order("slug ASC")
	}}, nil)
	if err != nil {
		return err
	}

	for _, cs := range comics {
		createdComic := false
		// Find or create comic by slug
		comic, err := s.comicRepo.FindOneWithTransaction(tx, []any{clause.Eq{Column: "slug", Value: cs.Slug}}, nil)
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
			if err := s.comicRepo.CreateWithTransaction(tx, comic); err != nil {
				return err
			}
			createdComic = true
		}

		if createdComic {
			if err := s.assignComicOwnership(tx, comic, users, translationGroups); err != nil {
				return err
			}
		}

		authors, err := s.lookupAuthors(tx, cs.Authors)
		if err != nil {
			return err
		}

		artists, err := s.lookupAuthors(tx, cs.Artists)
		if err != nil {
			return err
		}

		genres, err := s.lookupGenres(tx, cs.Genres)
		if err != nil {
			return err
		}
		tags, err := s.lookupTags(tx, cs.Tags)
		if err != nil {
			return err
		}

		if createdComic {
			if err := s.comicRepo.UpdateComicWithTransaction(tx, comic.ID, map[string]any{}, map[string]any{
				"Authors": authors,
				"Artists": artists,
				"Genres":  genres,
				"Tags":    tags,
			}); err != nil {
				return err
			}
		}

		// Seed chapters and their pages
		for chapterIndex, ch := range cs.Chapters {
			publishedAt := resolveChapterPublishedAt(comic.Slug, chapterIndex+1)
			createdChapter := false
			chapter, err := s.chapterRepo.FindOneWithTransaction(tx, []any{
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
					Slug:        comic.Slug + "-ch-" + ch.Number,
					IsPublished: true,
					PublishedAt: publishedAt,
				}
				if err := s.chapterRepo.CreateWithTransaction(tx, chapter); err != nil {
					return err
				}
				createdChapter = true
			}

			// Seed pages for this chapter
			pageCount := resolveChapterPageCount(chapter.Slug)
			for pageNumber := 1; pageNumber <= pageCount; pageNumber++ {
				if err := s.upsertSeedPage(tx, comic.Type, chapter, pageNumber, findSeedImageURL(ch.Pages, pageNumber), false, createdChapter); err != nil {
					return err
				}
			}
		}
	}

	return s.seedFakeComics(tx, users, translationGroups)
}

func (s *ComicSeeder) assignComicOwnership(tx *gorm.DB, comic *model.Comic, users []*model.User, translationGroups []*model.TranslationGroup) error {
	data := map[string]any{}
	if len(users) > 0 {
		uploadedByID := users[len(comic.Slug)%len(users)].ID
		data["uploaded_by_id"] = uploadedByID
	}
	if len(translationGroups) > 0 {
		translationGroupID := translationGroups[len(comic.Title)%len(translationGroups)].ID
		data["translation_group_id"] = translationGroupID
	}

	if len(data) == 0 {
		return nil
	}

	return s.comicRepo.UpdateWithTransaction(tx, []any{clause.Eq{Column: "id", Value: comic.ID}}, data)
}

func (s *ComicSeeder) lookupAuthors(tx *gorm.DB, names []string) ([]*model.Author, error) {
	result := make([]*model.Author, 0, len(names))
	for _, name := range names {
		a, err := s.authorRepo.FindOneWithTransaction(tx, []any{clause.Eq{Column: "name", Value: name}}, nil)
		if err != nil {
			return nil, err
		}
		result = append(result, a)
	}
	return result, nil
}

func (s *ComicSeeder) lookupGenres(tx *gorm.DB, slugs []string) ([]*model.Genre, error) {
	result := make([]*model.Genre, 0, len(slugs))
	for _, slug := range slugs {
		g, err := s.genreRepo.FindOneWithTransaction(tx, []any{clause.Eq{Column: "slug", Value: slug}}, nil)
		if err != nil {
			return nil, err
		}
		result = append(result, g)
	}
	return result, nil
}

func (s *ComicSeeder) lookupTags(tx *gorm.DB, slugs []string) ([]*model.Tag, error) {
	result := make([]*model.Tag, 0, len(slugs))
	for _, slug := range slugs {
		t, err := s.tagRepo.FindOneWithTransaction(tx, []any{clause.Eq{Column: "slug", Value: slug}}, nil)
		if err != nil {
			return nil, err
		}
		result = append(result, t)
	}
	return result, nil
}
