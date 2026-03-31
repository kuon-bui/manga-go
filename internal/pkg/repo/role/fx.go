package rolerepo

import "go.uber.org/fx"

var Module = fx.Module(
	"role-repo",
	fx.Provide(NewRoleRepo),
)
