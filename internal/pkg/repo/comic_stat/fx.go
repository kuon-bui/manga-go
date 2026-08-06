package comicstatrepo

import "go.uber.org/fx"

var Module = fx.Module(
	"comic-stat-repo",
	fx.Provide(NewComicStatRepo),
)
