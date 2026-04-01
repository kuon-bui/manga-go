package repo

import (
	authorrepo "manga-go/internal/pkg/repo/author"
	chapterrepo "manga-go/internal/pkg/repo/chapter"
	comicrepo "manga-go/internal/pkg/repo/comic"
	genrerepo "manga-go/internal/pkg/repo/genre"
	permissionrepo "manga-go/internal/pkg/repo/permission"
	readinghistoryrepo "manga-go/internal/pkg/repo/readinghistory"
	rolerepo "manga-go/internal/pkg/repo/role"
	tagrepo "manga-go/internal/pkg/repo/tag"
	translationgrouprepo "manga-go/internal/pkg/repo/translation_group"
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
	chapterrepo.Module,
	translationgrouprepo.Module,
	rolerepo.Module,
	permissionrepo.Module,
	readinghistoryrepo.Module,
)
