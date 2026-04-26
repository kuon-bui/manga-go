package commentservice

import (
	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/model"
	"time"

	"github.com/google/uuid"
)

type ReactionCounts struct {
	LIKE  int64 `json:"LIKE"`
	LOVE  int64 `json:"LOVE"`
	HAHA  int64 `json:"HAHA"`
	WOW   int64 `json:"WOW"`
	SAD   int64 `json:"SAD"`
	ANGRY int64 `json:"ANGRY"`
}

type CommentUserInfo struct {
	ID     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	Avatar *string   `json:"avatar"`
}

type CommentResponse struct {
	ID             uuid.UUID          `json:"id"`
	Content        string             `json:"content"`
	Author         *CommentUserInfo   `json:"author"`
	CreatedAt      *time.Time         `json:"createdAt"`
	ReactionCounts ReactionCounts     `json:"reactionCounts"`
	UserReaction   *string            `json:"userReaction"` // nil = no reaction
	ReplyCount     int                `json:"replyCount"`
	Replies        []*CommentResponse `json:"replies,omitempty"`
}

func mapCommentToResponse(comment *model.Comment, reactionCounts map[uuid.UUID]map[string]int64, userReactions map[uuid.UUID]string) *CommentResponse {
	var author *CommentUserInfo
	if comment.User != nil {
		var avatar *string
		if comment.User.Avatar != nil {
			avatarURL := common.AddFileContentPrefix(*comment.User.Avatar)
			avatar = &avatarURL
		}

		author = &CommentUserInfo{
			ID:     comment.User.ID,
			Name:   comment.User.Name,
			Avatar: avatar,
		}
	}

	replyCount := 0
	if comment.Replies != nil {
		replyCount = len(comment.Replies)
	}

	counts := ReactionCounts{}
	if countMap, ok := reactionCounts[comment.ID]; ok {
		counts.LIKE = countMap["LIKE"]
		counts.LOVE = countMap["LOVE"]
		counts.HAHA = countMap["HAHA"]
		counts.WOW = countMap["WOW"]
		counts.SAD = countMap["SAD"]
		counts.ANGRY = countMap["ANGRY"]
	}

	var userReaction *string
	if reaction, ok := userReactions[comment.ID]; ok && reaction != "" {
		userReaction = &reaction
	}

	var replies []*CommentResponse
	if comment.Replies != nil {
		replies = make([]*CommentResponse, 0, len(comment.Replies))
		for _, reply := range comment.Replies {
			replies = append(replies, mapCommentToResponse(reply, reactionCounts, userReactions))
		}
	}

	return &CommentResponse{
		ID:             comment.ID,
		Content:        comment.Content,
		Author:         author,
		CreatedAt:      comment.CreatedAt,
		ReactionCounts: counts,
		UserReaction:   userReaction,
		ReplyCount:     replyCount,
		Replies:        replies,
	}
}

func mapCommentsToResponses(comments []*model.Comment, reactionCounts map[uuid.UUID]map[string]int64, userReactions map[uuid.UUID]string) []*CommentResponse {
	responses := make([]*CommentResponse, 0, len(comments))
	for _, comment := range comments {
		responses = append(responses, mapCommentToResponse(comment, reactionCounts, userReactions))
	}
	return responses
}
