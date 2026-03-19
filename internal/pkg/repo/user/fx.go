package userrepo

import "go.uber.org/fx"

var Module = fx.Module(
	"user-repo",
	fx.Provide(NewUserRepository),
)
