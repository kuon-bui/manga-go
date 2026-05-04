package pageservice

import "go.uber.org/fx"

var Module = fx.Module(
	"page-service",
	fx.Provide(NewPageService),
)
