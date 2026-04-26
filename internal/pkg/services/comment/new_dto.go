package commentservice

import (
	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/constant"
	"manga-go/internal/pkg/model"
	"time"

	"github.com/google/uuid"
)

type NewCommentUserInfo struct {
	ID     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	Avatar *string   `json:"avatar"`
}

type NewCommentComicInfo struct {
	ID     uuid.UUID            `json:"id"`
	Title  string               `json:"title"`
	Slug   string               `json:"slug"`
	Status constant.ComicStatus `json:"status"`
}

type NewCommentChapterInfo struct {
	ID     uuid.UUID `json:"id"`
	Title  string    `json:"title"`
	Slug   string    `json:"slug"`
	Number string    `json:"number"`
}

type NewCommentResponse struct {
	ID             uuid.UUID              `json:"id"`
	Content        string                 `json:"content"`
	CreatedAt      *time.Time             `json:"createdAt"`
	Author         *NewCommentUserInfo    `json:"author"`
	Comic          *NewCommentComicInfo   `json:"comic"`
	Chapter        *NewCommentChapterInfo `json:"chapter"`
	ReactionCounts ReactionCounts         `json:"reactionCounts"`
}

func mapNewCommentToResponse(comment *model.Comment, reactionCounts map[uuid.UUID]map[string]int64) *NewCommentResponse {
	counts := ReactionCounts{}
	if countMap, ok := reactionCounts[comment.ID]; ok {
		counts.LIKE = countMap["LIKE"]
		counts.LOVE = countMap["LOVE"]
		counts.HAHA = countMap["HAHA"]
		counts.WOW = countMap["WOW"]
		counts.SAD = countMap["SAD"]
		counts.ANGRY = countMap["ANGRY"]
	}

	var author *NewCommentUserInfo
	if comment.User != nil {
		var avatar *string
		if comment.User.Avatar != nil {
			avatarURL := common.AddFileContentPrefix(*comment.User.Avatar)
			avatar = &avatarURL
		}

		author = &NewCommentUserInfo{
			ID:     comment.User.ID,
			Name:   comment.User.Name,
			Avatar: avatar,
		}
	}

	var comic *NewCommentComicInfo
	if comment.Comic != nil {
		comic = &NewCommentComicInfo{
			ID:     comment.Comic.ID,
			Title:  comment.Comic.Title,
			Slug:   comment.Comic.Slug,
			Status: comment.Comic.Status,
		}
	}

	var chapter *NewCommentChapterInfo
	if comment.Chapter != nil {
		chapter = &NewCommentChapterInfo{
			ID:     comment.Chapter.ID,
			Title:  comment.Chapter.Title,
			Slug:   comment.Chapter.Slug,
			Number: comment.Chapter.Number,
		}
	}

	return &NewCommentResponse{
		ID:             comment.ID,
		Content:        comment.Content,
		CreatedAt:      comment.CreatedAt,
		Author:         author,
		Comic:          comic,
		Chapter:        chapter,
		ReactionCounts: counts,
	}
}

func mapNewCommentsToResponses(comments []*model.Comment, reactionCounts map[uuid.UUID]map[string]int64) []*NewCommentResponse {
	responses := make([]*NewCommentResponse, 0, len(comments))
	for _, comment := range comments {
		responses = append(responses, mapNewCommentToResponse(comment, reactionCounts))
	}

	return responses
}
