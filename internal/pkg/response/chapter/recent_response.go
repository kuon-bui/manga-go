package chapterresponse

import "manga-go/internal/pkg/model"

type RecentUpdateResponse struct {
	Title   *model.Comic   `json:"title"`
	Chapter *model.Chapter `json:"chapter"`
}
