package fileroute

import (
	"manga-go/internal/app/api/common"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"file-route",
	common.ProvideAsRoute(NewFileRoute),
	fx.Provide(NewFileHandler),
)
