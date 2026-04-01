package commentservice

import "go.uber.org/fx"

var Module = fx.Module(
	"comment-service",
	fx.Provide(NewCommentService),
)
