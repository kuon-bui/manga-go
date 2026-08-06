package comicstatsworker

import (
	"context"
	"fmt"
	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/logger"
	"manga-go/internal/pkg/model"
	chapterrepo "manga-go/internal/pkg/repo/chapter"
	comicfollowrepo "manga-go/internal/pkg/repo/comic_follow"
	comicstatrepo "manga-go/internal/pkg/repo/comic_stat"
	ratingrepo "manga-go/internal/pkg/repo/rating"
	"runtime/debug"

	"github.com/google/uuid"
	"go.uber.org/fx"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type StatsServiceParams struct {
	fx.In
	Logger          *logger.Logger
	DB              *gorm.DB
	ComicStatRepo   *comicstatrepo.ComicStatRepo
	ComicFollowRepo *comicfollowrepo.ComicFollowRepo
	RatingRepo      *ratingrepo.RatingRepo
	ChapterRepo     *chapterrepo.ChapterRepo
}

type StatsService struct {
	logger          *logger.Logger
	db              *gorm.DB
	comicStatRepo   *comicstatrepo.ComicStatRepo
	comicFollowRepo *comicfollowrepo.ComicFollowRepo
	ratingRepo      *ratingrepo.RatingRepo
	chapterRepo     *chapterrepo.ChapterRepo
}

func NewStatsService(p StatsServiceParams) *StatsService {
	return &StatsService{
		logger:          p.Logger,
		db:              p.DB,
		comicStatRepo:   p.ComicStatRepo,
		comicFollowRepo: p.ComicFollowRepo,
		ratingRepo:      p.RatingRepo,
		chapterRepo:     p.ChapterRepo,
	}
}

func (s *StatsService) RecomputeComicStats(ctx context.Context, comicID uuid.UUID) (err error) {
	tx := s.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = fmt.Errorf("panic occurred: %v", r)
			common.ShowDebugTrace("RecomputeComicStats panic", debug.Stack())
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit().Error
		}
	}()

	followCount, err := s.comicFollowRepo.CountAllWithTransaction(tx, []any{
		clause.Eq{Column: "comic_id", Value: comicID},
	})
	if err != nil {
		s.logger.Error("Failed to count comic follows", "comicID", comicID, "error", err)
		return err
	}

	ratingCount, err := s.ratingRepo.CountAllWithTransaction(tx, []any{
		clause.Eq{Column: "comic_id", Value: comicID},
	})
	if err != nil {
		s.logger.Error("Failed to count ratings", "comicID", comicID, "error", err)
		return err
	}

	chapterCount, err := s.chapterRepo.CountAllWithTransaction(tx, []any{
		clause.Eq{Column: "comic_id", Value: comicID},
	})
	if err != nil {
		s.logger.Error("Failed to count chapters", "comicID", comicID, "error", err)
		return err
	}

	avgRating, err := s.ratingRepo.FindAvgRatingByComicIDWithTransaction(tx, comicID)
	if err != nil {
		s.logger.Error("Failed to compute average rating", "comicID", comicID, "error", err)
		return err
	}

	s.logger.Info("Computed comic stats",
		"comicID", comicID,
		"followCount", followCount,
		"ratingCount", ratingCount,
		"chapterCount", chapterCount,
		"avgRating", avgRating,
	)

	stat := &model.ComicStat{
		ComicID:      comicID,
		FollowCount:  int(followCount),
		RatingCount:  int(ratingCount),
		ChapterCount: int(chapterCount),
		AvgRating:    avgRating,
	}

	if err := s.comicStatRepo.UpsertWithTransaction(
		tx,
		stat,
		[]string{"comic_id"},
		[]string{"follow_count", "rating_count", "chapter_count", "avg_rating", "updated_at"},
	); err != nil {
		s.logger.Error("Failed to upsert comic stats", "comicID", comicID, "error", err)
		return err
	}

	return nil
}
