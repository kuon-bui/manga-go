package comicrepo

import "go.uber.org/fx"

var Module = fx.Module(
	"comic-repo",
	fx.Provide(NewComicRepo),
)
