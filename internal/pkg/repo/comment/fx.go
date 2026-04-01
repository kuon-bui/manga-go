package commentrepo

import "go.uber.org/fx"

var Module = fx.Module(
	"commentrepo",
	fx.Provide(NewCommentRepo),
)
