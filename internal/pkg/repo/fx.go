package repo

import (
	authorrepo "manga-go/internal/pkg/repo/author"
	genrerepo "manga-go/internal/pkg/repo/genre"
	userrepo "manga-go/internal/pkg/repo/user"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"repo",
	userrepo.Module,
	authorrepo.Module,
	genrerepo.Module,
)
