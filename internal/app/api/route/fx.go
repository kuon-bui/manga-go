package route

import (
	authorroute "manga-go/internal/app/api/route/author"
	fileroute "manga-go/internal/app/api/route/file"
	genreroute "manga-go/internal/app/api/route/genre"
	userroute "manga-go/internal/app/api/route/user"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"route",
	userroute.Module,
	authorroute.Module,
	genreroute.Module,
	fileroute.Module,
)
