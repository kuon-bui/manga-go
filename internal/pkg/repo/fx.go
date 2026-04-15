package repo

import (
	authorrepo "manga-go/internal/pkg/repo/author"
	chapterrepo "manga-go/internal/pkg/repo/chapter"
	comicrepo "manga-go/internal/pkg/repo/comic"
	comicfollowrepo "manga-go/internal/pkg/repo/comic_follow"
	commentrepo "manga-go/internal/pkg/repo/comment"
	genrerepo "manga-go/internal/pkg/repo/genre"
	notificationrepo "manga-go/internal/pkg/repo/notification"
	pagerepo "manga-go/internal/pkg/repo/page"
	permissionrepo "manga-go/internal/pkg/repo/permission"
	ratingrepo "manga-go/internal/pkg/repo/rating"
	reactionrepo "manga-go/internal/pkg/repo/reaction"
	readinghistoryrepo "manga-go/internal/pkg/repo/reading_history"
	readingprogressrepo "manga-go/internal/pkg/repo/reading_progress"
	rolerepo "manga-go/internal/pkg/repo/role"
	tagrepo "manga-go/internal/pkg/repo/tag"
	translationgrouprepo "manga-go/internal/pkg/repo/translation_group"
	userrepo "manga-go/internal/pkg/repo/user"
	usercomicreadrepo "manga-go/internal/pkg/repo/user_comic_read"
	usernotificationrepo "manga-go/internal/pkg/repo/user_notification"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"repo",
	userrepo.Module,
	authorrepo.Module,
	genrerepo.Module,
	tagrepo.Module,
	comicfollowrepo.Module,
	comicrepo.Module,
	chapterrepo.Module,
	pagerepo.Module,
	translationgrouprepo.Module,
	rolerepo.Module,
	permissionrepo.Module,
	ratingrepo.Module,
	readinghistoryrepo.Module,
	commentrepo.Module,
	reactionrepo.Module,
	readingprogressrepo.Module,
	usercomicreadrepo.Module,
	notificationrepo.Module,
	usernotificationrepo.Module,
)
