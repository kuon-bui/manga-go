package seeder

import (
	"context"
	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/logger"
	"runtime/debug"

	"gorm.io/gorm"
)

// Seeder defines the interface that every seeder must implement.
type Seeder interface {
	Name() string
	Seed(tx *gorm.DB) error
}

// SeederRunner holds an ordered list of seeders and runs them sequentially.
type SeederRunner struct {
	seeders []Seeder
	logger  *logger.Logger
	db      *gorm.DB
}

func NewSeederRunner(seeders []Seeder, logger *logger.Logger, db *gorm.DB) *SeederRunner {
	return &SeederRunner{
		seeders: seeders,
		logger:  logger,
		db:      db,
	}
}

// RunAll executes every seeder in order. It stops on first error.
func (r *SeederRunner) RunAll(ctx context.Context) (err error) {
	tx := r.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			common.ShowDebugTrace("running seeders", debug.Stack())
			panic(r)
		} else if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	r.logger.Info("Starting seeders...")
	for _, s := range r.seeders {
		r.logger.Infof("Running seeder: %s", s.Name())
		if err := s.Seed(tx); err != nil {
			r.logger.Errorf("Seeder %s failed: %v", s.Name(), err)
			return err
		}
		r.logger.Infof("Seeder %s completed", s.Name())
	}
	r.logger.Info("All seeders completed successfully")
	return nil
}
