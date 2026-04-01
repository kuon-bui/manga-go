package readinghistoryrepo

import "go.uber.org/fx"

var Module = fx.Module(
	"readinghistory-repo",
	fx.Provide(NewReadingHistoryRepo),
)
