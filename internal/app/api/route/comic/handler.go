package comicroute

import (
	comicservice "manga-go/internal/pkg/services/comic"

	"go.uber.org/fx"
)

type ComicHandler struct {
	comicService *comicservice.ComicService
}

type ComicHandlerParams struct {
	fx.In

	ComicService *comicservice.ComicService
}

func NewComicHandler(p ComicHandlerParams) *ComicHandler {
	return &ComicHandler{
		comicService: p.ComicService,
	}
}
