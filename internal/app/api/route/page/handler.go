package pageroute

import (
	pageservice "manga-go/internal/pkg/services/page"

	"go.uber.org/fx"
)

type PageHandler struct {
	pageService *pageservice.PageService
}

type PageHandlerParams struct {
	fx.In

	PageService *pageservice.PageService
}

func NewPageHandler(p PageHandlerParams) *PageHandler {
	return &PageHandler{
		pageService: p.PageService,
	}
}
