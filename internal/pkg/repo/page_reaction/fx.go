package pagereactionrepo

import "go.uber.org/fx"

var Module = fx.Module(
	"page-reaction-repo",
	fx.Provide(NewPageReactionRepo),
)
