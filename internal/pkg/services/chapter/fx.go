package chapterserivce

import "go.uber.org/fx"

var Module = fx.Module(
	"chapter-service",
	fx.Provide(NewChapterService),
)
