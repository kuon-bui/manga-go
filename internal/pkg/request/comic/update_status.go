package comicrequest

import "manga-go/internal/pkg/constant"

type UpdateComicStatusRequest struct {
	Status constant.ComicStatus `json:"status" binding:"required,oneof=ongoing completed hiatus cancelled"`
}
