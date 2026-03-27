package chapterhandler

import (
	"manga-go/internal/app/api/common"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"chapter-route",
	common.ProvideAsRoute(NewChapterRoute),
	fx.Provide(NewChapterHandler),
)
