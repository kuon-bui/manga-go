package reactionrepo

import "go.uber.org/fx"

var Module = fx.Module(
	"reaction-repo",
	fx.Provide(NewReactionRepo),
)
