package chapterrepo

import "go.uber.org/fx"

var Module = fx.Module(
	"chapter-repo",
	fx.Provide(NewChapterRepo),
)
