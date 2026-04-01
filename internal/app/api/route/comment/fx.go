package commentroute

import (
	"manga-go/internal/app/api/common"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"comment-route",
	common.ProvideAsRoute(NewCommentRoute),
	fx.Provide(NewCommentHandler),
)
