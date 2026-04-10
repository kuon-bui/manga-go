package route

import (
	authorroute "manga-go/internal/app/api/route/author"
	chapterroute "manga-go/internal/app/api/route/chapter"
	comicroute "manga-go/internal/app/api/route/comic"
	commentroute "manga-go/internal/app/api/route/comment"
	fileroute "manga-go/internal/app/api/route/file"
	genreroute "manga-go/internal/app/api/route/genre"
	permissionroute "manga-go/internal/app/api/route/permission"
	ratingroute "manga-go/internal/app/api/route/rating"
	readinghistoryroute "manga-go/internal/app/api/route/reading_history"
	roleroute "manga-go/internal/app/api/route/role"
	swaggerrouter "manga-go/internal/app/api/route/swagger"
	tagroute "manga-go/internal/app/api/route/tag"
	translationgrouproute "manga-go/internal/app/api/route/translation_group"
	userroute "manga-go/internal/app/api/route/user"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"route",
	userroute.Module,
	authorroute.Module,
	genreroute.Module,
	fileroute.Module,
	tagroute.Module,
	comicroute.Module,
	chapterroute.Module,
	translationgrouproute.Module,
	roleroute.Module,
	permissionroute.Module,
	ratingroute.Module,
	readinghistoryroute.Module,
	commentroute.Module,
	swaggerrouter.Module,
)
