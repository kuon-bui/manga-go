package readinghistoryservice

import "go.uber.org/fx"

var Module = fx.Module(
	"readinghistory-service",
	fx.Provide(NewReadingHistoryService),
)
