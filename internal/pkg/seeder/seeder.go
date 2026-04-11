package seeder

import (
	"context"
	"manga-go/internal/pkg/logger"
)

// Seeder defines the interface that every seeder must implement.
type Seeder interface {
	Name() string
	Seed(ctx context.Context) error
}

// SeederRunner holds an ordered list of seeders and runs them sequentially.
type SeederRunner struct {
	seeders []Seeder
	logger  *logger.Logger
}

func NewSeederRunner(seeders []Seeder, logger *logger.Logger) *SeederRunner {
	return &SeederRunner{
		seeders: seeders,
		logger:  logger,
	}
}

// RunAll executes every seeder in order. It stops on first error.
func (r *SeederRunner) RunAll(ctx context.Context) error {
	r.logger.Info("Starting seeders...")
	for _, s := range r.seeders {
		r.logger.Infof("Running seeder: %s", s.Name())
		if err := s.Seed(ctx); err != nil {
			r.logger.Errorf("Seeder %s failed: %v", s.Name(), err)
			return err
		}
		r.logger.Infof("Seeder %s completed", s.Name())
	}
	r.logger.Info("All seeders completed successfully")
	return nil
}
