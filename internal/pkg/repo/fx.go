package repo

import (
	authorrepo "manga-go/internal/pkg/repo/author"
	chapterrepo "manga-go/internal/pkg/repo/chapter"
	comicrepo "manga-go/internal/pkg/repo/comic"
	commentrepo "manga-go/internal/pkg/repo/comment"
	genrerepo "manga-go/internal/pkg/repo/genre"
	permissionrepo "manga-go/internal/pkg/repo/permission"
	reactionrepo "manga-go/internal/pkg/repo/reaction"
	readinghistoryrepo "manga-go/internal/pkg/repo/reading_history"
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
	commentrepo.Module,
	reactionrepo.Module,
)
