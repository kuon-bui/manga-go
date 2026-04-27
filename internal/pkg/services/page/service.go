package pageservice

import (
	"manga-go/internal/pkg/logger"
	pagerepo "manga-go/internal/pkg/repo/page"
	pagereactionrepo "manga-go/internal/pkg/repo/page_reaction"

	"go.uber.org/fx"
)

type PageService struct {
	logger           *logger.Logger
	pageRepo         *pagerepo.PageRepo
	pageReactionRepo *pagereactionrepo.PageReactionRepo
}

type PageServiceParams struct {
	fx.In

	Logger           *logger.Logger
	PageRepo         *pagerepo.PageRepo
	PageReactionRepo *pagereactionrepo.PageReactionRepo
}

func NewPageService(p PageServiceParams) *PageService {
	return &PageService{
		logger:           p.Logger,
		pageRepo:         p.PageRepo,
		pageReactionRepo: p.PageReactionRepo,
	}
}
