package translationgrouproute

import (
	translationgroupservice "manga-go/internal/pkg/services/translation_group"

	"go.uber.org/fx"
)

type TranslationGroupHandler struct {
	translationGroupService *translationgroupservice.TranslationGroupService
}

type TranslationGroupHandlerParams struct {
	fx.In

	TranslationGroupService *translationgroupservice.TranslationGroupService
}

func NewTranslationGroupHandler(p TranslationGroupHandlerParams) *TranslationGroupHandler {
	return &TranslationGroupHandler{
		translationGroupService: p.TranslationGroupService,
	}
}
