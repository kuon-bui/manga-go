package translationgroupservice

import (
	casbinpkg "manga-go/internal/pkg/casbin"
	"manga-go/internal/pkg/logger"
	translationgrouprepo "manga-go/internal/pkg/repo/translation_group"
	userrepo "manga-go/internal/pkg/repo/user"

	"go.uber.org/fx"
)

type TranslationGroupService struct {
	logger               *logger.Logger
	translationGroupRepo *translationgrouprepo.TranslationGroupRepo
	userRepo             *userrepo.UserRepository
	enforcer             *casbinpkg.Enforcer
}

type TranslationGroupServiceParams struct {
	fx.In
	Logger               *logger.Logger
	TranslationGroupRepo *translationgrouprepo.TranslationGroupRepo
	UserRepo             *userrepo.UserRepository
	Enforcer             *casbinpkg.Enforcer
}

func NewTranslationGroupService(params TranslationGroupServiceParams) *TranslationGroupService {
	return &TranslationGroupService{
		logger:               params.Logger,
		translationGroupRepo: params.TranslationGroupRepo,
		userRepo:             params.UserRepo,
		enforcer:             params.Enforcer,
	}
}
