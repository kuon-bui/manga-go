package authorizationaudit

import "go.uber.org/fx"

var Module = fx.Module(
	"authorization-audit-repo",
	fx.Provide(NewRepo),
)
