package comicfollowrepo

import "go.uber.org/fx"

var Module = fx.Module(
	"comic-follow-repo",
	fx.Provide(NewComicFollowRepo),
)
