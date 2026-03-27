package translationgrouprepo

import "go.uber.org/fx"

var Module = fx.Module(
	"translation-group-repo",
	fx.Provide(NewTranslationGroupRepo),
)
