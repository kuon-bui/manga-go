package tagrepo

import "go.uber.org/fx"

var Module = fx.Module(
	"tag-repo",
	fx.Provide(NewTagRepo),
)
