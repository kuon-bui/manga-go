package repo

import (
	userrepo "base-go/internal/pkg/repo/user"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"repo",
	userrepo.Module,
)
