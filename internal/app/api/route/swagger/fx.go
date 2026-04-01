package swaggerrouter

import (
	"manga-go/internal/app/api/common"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"swagger-route",
	common.ProvideAsRoute(NewSwaggerRoute),
)
