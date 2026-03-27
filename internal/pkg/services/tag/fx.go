package tagservice

import "go.uber.org/fx"

var Module = fx.Module(
	"tag-service",
	fx.Provide(NewTagService),
)
