package repo

import (
	authorrepo "manga-go/internal/pkg/repo/author"
	userrepo "manga-go/internal/pkg/repo/user"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"repo",
	userrepo.Module,
	authorrepo.Module,
)
