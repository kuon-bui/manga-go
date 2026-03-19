package route

import (
	userroute "base-go/internal/app/api/route/user"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"route",
	userroute.Module,
)
