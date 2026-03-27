package translationgrouproute

import (
	"manga-go/internal/app/api/common"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"translation-group-route",
	common.ProvideAsRoute(NewTranslationGroupRoute),
	fx.Provide(NewTranslationGroupHandler),
)
