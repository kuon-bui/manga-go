package comicservice

import "go.uber.org/fx"

var Module = fx.Module(
	"comic-service",
	fx.Provide(NewComicService),
)
