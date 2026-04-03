package readingprogressrepo

import "go.uber.org/fx"

var Module = fx.Module(
	"reading-progress-repo",
	fx.Provide(
		NewReadingProgressRepo,
	),
)
