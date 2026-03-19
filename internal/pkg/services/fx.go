package services

import (
	userserivce "base-go/internal/pkg/services/user"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"services",
	userserivce.Module,
)
