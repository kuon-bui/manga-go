package repo

import (
	userrepo "manga-go/internal/pkg/repo/user"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"repo",
	userrepo.Module,
)
