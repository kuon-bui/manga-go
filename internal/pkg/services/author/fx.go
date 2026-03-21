package authorservice

import "go.uber.org/fx"

var Module = fx.Module(
	"author-service",
	fx.Provide(NewAuthorService),
)
