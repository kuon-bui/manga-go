package activityseeder

import (
	"errors"
	"fmt"
	"manga-go/internal/pkg/bitset"
	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/model"
	chapterrepo "manga-go/internal/pkg/repo/chapter"
	comicrepo "manga-go/internal/pkg/repo/comic"
	comicfollowrepo "manga-go/internal/pkg/repo/comic_follow"
	commentrepo "manga-go/internal/pkg/repo/comment"
	commentreportrepo "manga-go/internal/pkg/repo/comment_report"
	pagerepo "manga-go/internal/pkg/repo/page"
	pagereactionrepo "manga-go/internal/pkg/repo/page_reaction"
	ratingrepo "manga-go/internal/pkg/repo/rating"
	reactionrepo "manga-go/internal/pkg/repo/reaction"
	readinghistoryrepo "manga-go/internal/pkg/repo/reading_history"
	readingprogressrepo "manga-go/internal/pkg/repo/reading_progress"
	userrepo "manga-go/internal/pkg/repo/user"
	usercomicreadrepo "manga-go/internal/pkg/repo/user_comic_read"
	seederutil "manga-go/internal/pkg/seeder/util"

	"github.com/google/uuid"
	"github.com/jaswdr/faker/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ActivitySeeder struct {
	userRepo            *userrepo.UserRepository
	comicRepo           *comicrepo.ComicRepo
	chapterRepo         *chapterrepo.ChapterRepo
	pageRepo            *pagerepo.PageRepo
	commentRepo         *commentrepo.CommentRepo
	ratingRepo          *ratingrepo.RatingRepo
	readingHistoryRepo  *readinghistoryrepo.ReadingHistoryRepo
	readingProgressRepo *readingprogressrepo.ReadingProgressRepo
	comicFollowRepo     *comicfollowrepo.ComicFollowRepo
	userComicReadRepo   *usercomicreadrepo.UserComicReadRepo
	reactionRepo        *reactionrepo.ReactionRepo
	pageReactionRepo    *pagereactionrepo.PageReactionRepo
	commentReportRepo   *commentreportrepo.CommentReportRepo
	faker               faker.Faker
}

func NewActivitySeeder(
	userRepo *userrepo.UserRepository,
	comicRepo *comicrepo.ComicRepo,
	chapterRepo *chapterrepo.ChapterRepo,
	pageRepo *pagerepo.PageRepo,
	commentRepo *commentrepo.CommentRepo,
	ratingRepo *ratingrepo.RatingRepo,
	readingHistoryRepo *readinghistoryrepo.ReadingHistoryRepo,
	readingProgressRepo *readingprogressrepo.ReadingProgressRepo,
	comicFollowRepo *comicfollowrepo.ComicFollowRepo,
	userComicReadRepo *usercomicreadrepo.UserComicReadRepo,
	reactionRepo *reactionrepo.ReactionRepo,
	pageReactionRepo *pagereactionrepo.PageReactionRepo,
	commentReportRepo *commentreportrepo.CommentReportRepo,
	faker faker.Faker,
) *ActivitySeeder {
	return &ActivitySeeder{
		userRepo:            userRepo,
		comicRepo:           comicRepo,
		chapterRepo:         chapterRepo,
		pageRepo:            pageRepo,
		commentRepo:         commentRepo,
		ratingRepo:          ratingRepo,
		readingHistoryRepo:  readingHistoryRepo,
		readingProgressRepo: readingProgressRepo,
		comicFollowRepo:     comicFollowRepo,
		userComicReadRepo:   userComicReadRepo,
		reactionRepo:        reactionRepo,
		pageReactionRepo:    pageReactionRepo,
		commentReportRepo:   commentReportRepo,
		faker:               faker,
	}
}

func (s *ActivitySeeder) Name() string {
	return "ActivitySeeder"
}

func (s *ActivitySeeder) Truncate(tx *gorm.DB) error {
	return seederutil.TruncateTables(
		tx,
		"comment_reports",
		"comment_reactions",
		"page_reactions",
		"ratings",
		"reading_histories",
		"reading_progresses",
		"user_comic_reads",
		"comments",
		"comic_follows",
	)
}

func (s *ActivitySeeder) Seed(tx *gorm.DB) error {
	comics, err := s.comicRepo.FindAllWithTx(tx, []any{func(db *gorm.DB) *gorm.DB {
		return db.Order("title ASC")
	}}, map[string]common.MoreKeyOption{
		"Chapters": {
			Custom: func(db *gorm.DB) *gorm.DB {
				return db.Order("chapter_idx ASC")
			},
		},
	})
	if err != nil {
		return err
	}
	users, err := s.userRepo.FindAllWithTx(tx, []any{func(db *gorm.DB) *gorm.DB {
		return db.Order("email ASC")
	}}, nil)
	if err != nil {
		return err
	}
	pages, err := s.pageRepo.FindAllWithTx(tx, []any{func(db *gorm.DB) *gorm.DB {
		return db.Order("page_number ASC")
	}}, nil)
	if err != nil {
		return err
	}
	if len(users) == 0 || len(comics) == 0 {
		return nil
	}

	pagesByChapter := make(map[uuid.UUID][]*model.Page)
	for _, page := range pages {
		pagesByChapter[page.ChapterID] = append(pagesByChapter[page.ChapterID], page)
	}

	for index, user := range users {
		comic := comics[index%len(comics)]
		if len(comic.Chapters) == 0 {
			continue
		}
		chapter := comic.Chapters[index%len(comic.Chapters)]
		chapterPages := pagesByChapter[chapter.ID]

		if err := s.seedFollow(tx, user, comic); err != nil {
			return err
		}
		comment, err := s.seedCommentThread(tx, users, user, comic, chapter, index)
		if err != nil {
			return err
		}
		if err := s.seedRating(tx, user, comic, index); err != nil {
			return err
		}
		if err := s.seedReadingData(tx, user, comic, chapter, index); err != nil {
			return err
		}
		if len(chapterPages) > 0 {
			if err := s.seedPageReaction(tx, users[(index+2)%len(users)], chapterPages[0], index); err != nil {
				return err
			}
		}
		if comment != nil {
			if err := s.seedCommentReaction(tx, users[(index+1)%len(users)], comment, index); err != nil {
				return err
			}
			if index%3 == 0 {
				if err := s.seedCommentReport(tx, users[(index+2)%len(users)], comment, index); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (s *ActivitySeeder) seedFollow(tx *gorm.DB, user *model.User, comic *model.Comic) error {
	_, err := s.comicFollowRepo.FindOneWithTransaction(tx, []any{
		clause.Eq{Column: "user_id", Value: user.ID},
		clause.Eq{Column: "comic_id", Value: comic.ID},
	}, nil)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		follow := &model.ComicFollow{UserID: user.ID, ComicID: comic.ID}
		follow.Fake(s.faker)
		return s.comicFollowRepo.CreateWithTransaction(tx, follow)
	}

	return nil
}

func (s *ActivitySeeder) seedRating(tx *gorm.DB, user *model.User, comic *model.Comic, index int) error {
	_, err := s.ratingRepo.FindOneWithTransaction(tx, []any{
		clause.Eq{Column: "user_id", Value: user.ID},
		clause.Eq{Column: "comic_id", Value: comic.ID},
	}, nil)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		rating := &model.Rating{UserId: user.ID, ComicId: comic.ID}
		rating.Fake(s.faker)
		comment := fmt.Sprintf("[seed-rating-%02d] %s", index+1, *rating.Comment)
		rating.Comment = &comment
		return s.ratingRepo.CreateWithTransaction(tx, rating)
	}

	return nil
}

func (s *ActivitySeeder) seedReadingData(tx *gorm.DB, user *model.User, comic *model.Comic, chapter *model.Chapter, index int) error {
	_, err := s.readingHistoryRepo.FindOneWithTransaction(tx, []any{
		clause.Eq{Column: "user_id", Value: user.ID},
		clause.Eq{Column: "chapter_id", Value: chapter.ID},
	}, nil)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		history := &model.ReadingHistory{UserID: user.ID, ComicID: comic.ID, ChapterID: chapter.ID}
		history.Fake(s.faker)
		if err := s.readingHistoryRepo.CreateWithTransaction(tx, history); err != nil {
			return err
		}
	}

	_, err = s.readingProgressRepo.FindOneWithTransaction(tx, []any{
		clause.Eq{Column: "user_id", Value: user.ID},
		clause.Eq{Column: "comic_id", Value: comic.ID},
		clause.Eq{Column: "chapter_id", Value: chapter.ID},
	}, nil)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		progress := &model.ReadingProgress{UserID: user.ID, ComicID: comic.ID, ChapterID: chapter.ID}
		progress.Fake(s.faker)
		if err := s.readingProgressRepo.CreateWithTransaction(tx, progress); err != nil {
			return err
		}
	}

	_, err = s.userComicReadRepo.FindOneWithTransaction(tx, []any{
		clause.Eq{Column: "user_id", Value: user.ID},
		clause.Eq{Column: "comic_id", Value: comic.ID},
	}, nil)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		read := &model.UserComicRead{UserID: user.ID, ComicID: comic.ID}
		read.Fake(s.faker)
		read.ReadData = bitset.NewReadBitset(len(comic.Chapters))
		limit := index%len(comic.Chapters) + 1
		for chapterIndex := 0; chapterIndex < limit; chapterIndex++ {
			read.ReadData.Mark(chapterIndex)
		}
		if err := s.userComicReadRepo.CreateWithTransaction(tx, read); err != nil {
			return err
		}
	}

	return nil
}

func (s *ActivitySeeder) seedCommentThread(tx *gorm.DB, users []*model.User, author *model.User, comic *model.Comic, chapter *model.Chapter, index int) (*model.Comment, error) {
	content := fmt.Sprintf("[seed-comment-%02d] %s", index+1, s.faker.Lorem().Sentence(12))
	comment, err := s.commentRepo.FindOneWithTransaction(tx, []any{
		clause.Eq{Column: "user_id", Value: author.ID},
		clause.Eq{Column: "comic_id", Value: comic.ID},
		clause.Eq{Column: "content", Value: content},
	}, nil)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		comment = &model.Comment{UserId: author.ID, ComicId: comic.ID, ChapterId: &chapter.ID, Content: content}
		comment.Fake(s.faker)
		comment.Content = content
		pageIndex := index % 4
		comment.PageIndex = &pageIndex
		if err := s.commentRepo.CreateWithTransaction(tx, comment); err != nil {
			return nil, err
		}
	}

	if len(users) > 1 {
		responder := users[(index+1)%len(users)]
		if responder.ID != author.ID {
			replyContent := fmt.Sprintf("[seed-reply-%02d] %s", index+1, s.faker.Lorem().Sentence(10))
			_, err := s.commentRepo.FindOneWithTransaction(tx, []any{
				clause.Eq{Column: "user_id", Value: responder.ID},
				clause.Eq{Column: "parent_id", Value: comment.ID},
				clause.Eq{Column: "content", Value: replyContent},
			}, nil)
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, err
			}
			if errors.Is(err, gorm.ErrRecordNotFound) {
				reply := &model.Comment{UserId: responder.ID, ComicId: comic.ID, ChapterId: &chapter.ID, ParentId: &comment.ID, Content: replyContent}
				reply.Fake(s.faker)
				reply.Content = replyContent
				if err := s.commentRepo.CreateWithTransaction(tx, reply); err != nil {
					return nil, err
				}
			}
		}
	}

	return comment, nil
}

func (s *ActivitySeeder) seedCommentReaction(tx *gorm.DB, user *model.User, comment *model.Comment, index int) error {
	_, err := s.reactionRepo.FindOneWithTransaction(tx, []any{
		clause.Eq{Column: "user_id", Value: user.ID},
		clause.Eq{Column: "comment_id", Value: comment.ID},
	}, nil)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		reaction := &model.CommentReaction{CommentId: comment.ID}
		reaction.Fake(s.faker)
		reaction.UserId = user.ID
		return s.reactionRepo.CreateWithTransaction(tx, reaction)
	}

	return nil
}

func (s *ActivitySeeder) seedPageReaction(tx *gorm.DB, user *model.User, page *model.Page, index int) error {
	_, err := s.pageReactionRepo.FindOneWithTransaction(tx, []any{
		clause.Eq{Column: "user_id", Value: user.ID},
		clause.Eq{Column: "page_id", Value: page.ID},
	}, nil)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		reaction := &model.PageReaction{PageId: page.ID}
		reaction.Fake(s.faker)
		reaction.UserId = user.ID
		return s.pageReactionRepo.CreateWithTransaction(tx, reaction)
	}

	return nil
}

func (s *ActivitySeeder) seedCommentReport(tx *gorm.DB, user *model.User, comment *model.Comment, index int) error {
	if user.ID == comment.UserId {
		return nil
	}

	_, err := s.commentReportRepo.FindOneWithTransaction(tx, []any{
		clause.Eq{Column: "user_id", Value: user.ID},
		clause.Eq{Column: "comment_id", Value: comment.ID},
	}, nil)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		report := &model.CommentReport{UserId: user.ID, CommentId: comment.ID}
		report.Fake(s.faker)
		return s.commentReportRepo.CreateWithTransaction(tx, report)
	}

	return nil
}
