package pagerepo

import "go.uber.org/fx"

var Module = fx.Module(
	"page-repo",
	fx.Provide(NewPageRepo),
)
