package usercomicreadrepo

import "go.uber.org/fx"

var Module = fx.Module(
	"user-comic-read-repo",
	fx.Provide(NewUserComicReadRepo),
)
