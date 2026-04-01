package readinghistoryrepo

import "go.uber.org/fx"

var Module = fx.Module(
	"reading-history-repo",
	fx.Provide(NewReadingHistoryRepo),
)
