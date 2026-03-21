package services

import (
	authorservice "manga-go/internal/pkg/services/author"
	userserivce "manga-go/internal/pkg/services/user"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"services",
	userserivce.Module,
	authorservice.Module,
)
