package fileservice

import "go.uber.org/fx"

var Module = fx.Module(
	"file_service",
	fx.Provide(
		NewFileService,
	),
)
