package authorizationrevision

import "go.uber.org/fx"

var Module = fx.Module(
	"authorization-revision-repo",
	fx.Provide(NewRepo),
)
