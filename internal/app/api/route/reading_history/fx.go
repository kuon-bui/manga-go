package readinghistoryroute

import (
	"manga-go/internal/app/api/common"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"readinghistory-route",
	common.ProvideAsRoute(NewReadingHistoryRoute),
	fx.Provide(NewReadingHistoryHandler),
)
