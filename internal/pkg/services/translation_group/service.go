package translationgroupservice

import (
	"manga-go/internal/pkg/logger"
	translationgrouprepo "manga-go/internal/pkg/repo/translation_group"
	userrepo "manga-go/internal/pkg/repo/user"

	"go.uber.org/fx"
)

type TranslationGroupService struct {
	logger               *logger.Logger
	translationGroupRepo *translationgrouprepo.TranslationGroupRepo
	userRepo             *userrepo.UserRepository
}

type TranslationGroupServiceParams struct {
	fx.In
	Logger               *logger.Logger
	TranslationGroupRepo *translationgrouprepo.TranslationGroupRepo
	UserRepo             *userrepo.UserRepository
}

func NewTranslationGroupService(params TranslationGroupServiceParams) *TranslationGroupService {
	return &TranslationGroupService{
		logger:               params.Logger,
		translationGroupRepo: params.TranslationGroupRepo,
		userRepo:             params.UserRepo,
	}
}
