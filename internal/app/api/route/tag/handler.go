package tagroute

import (
	tagservice "manga-go/internal/pkg/services/tag"

	"go.uber.org/fx"
)

type TagHandler struct {
	tagService *tagservice.TagService
}

type TagHandlerParams struct {
	fx.In

	TagService *tagservice.TagService
}

func NewTagHandler(p TagHandlerParams) *TagHandler {
	return &TagHandler{
		tagService: p.TagService,
	}
}
