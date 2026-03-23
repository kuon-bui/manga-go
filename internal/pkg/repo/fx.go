package repo

import (
	authorrepo "manga-go/internal/pkg/repo/author"
	comicrepo "manga-go/internal/pkg/repo/comic"
	genrerepo "manga-go/internal/pkg/repo/genre"
	tagrepo "manga-go/internal/pkg/repo/tag"
	userrepo "manga-go/internal/pkg/repo/user"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"repo",
	userrepo.Module,
	authorrepo.Module,
	genrerepo.Module,
	tagrepo.Module,
	comicrepo.Module,
)
