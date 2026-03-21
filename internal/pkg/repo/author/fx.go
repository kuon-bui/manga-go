package authorrepo

import "go.uber.org/fx"

var Module = fx.Module(
	"author-repo",
	fx.Provide(NewAuthorRepo),
)
