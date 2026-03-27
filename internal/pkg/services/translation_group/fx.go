package translationgroupservice

import "go.uber.org/fx"

var Module = fx.Module(
	"translation-group-service",
	fx.Provide(NewTranslationGroupService),
)
