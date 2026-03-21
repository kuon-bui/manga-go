package route

import (
	authorroute "manga-go/internal/app/api/route/author"
	userroute "manga-go/internal/app/api/route/user"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"route",
	userroute.Module,
	authorroute.Module,
)
