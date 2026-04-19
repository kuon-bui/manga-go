package commentservice

import (
	"context"
	"encoding/json"
	"errors"
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/model"
	commentrequest "manga-go/internal/pkg/request/comment"
	"manga-go/internal/pkg/utils"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm/clause"
)

func (s *CommentService) ListComments(ctx context.Context, req *commentrequest.ListCommentsRequest) response.Result {
	// Validate: must provide either comicId or chapterId
	if strings.TrimSpace(req.ComicId) == "" && strings.TrimSpace(req.ChapterId) == "" {
		return response.ResultError("Either comicId or chapterId query parameter is required")
	}

	conditions := []any{
		clause.Eq{Column: "parent_id", Value: nil},
	}

	if strings.TrimSpace(req.ChapterId) != "" {
		// Chapter-level or page-level comments
		chapterID, err := parseChapterID(req.ChapterId)
		if err != nil {
			return response.ResultError("invalid chapterId")
		}
		conditions = append(conditions, clause.Eq{Column: "chapter_id", Value: chapterID})

		if req.PageIndex != nil {
			conditions = append(conditions, clause.Eq{Column: "page_index", Value: *req.PageIndex})
		}
	} else {
		// Comic-level comments: comicId provided, chapterId must be nil
		comicID, err := uuid.Parse(strings.TrimSpace(req.ComicId))
		if err != nil {
			return response.ResultError("invalid comicId")
		}
		conditions = append(conditions,
			clause.Eq{Column: "comic_id", Value: comicID},
			clause.Eq{Column: "chapter_id", Value: nil},
		)
	}

	comments, total, err := s.commentRepo.FindPaginated(ctx, conditions, &req.Paging, map[string]common.MoreKeyOption{
		"User":                                 {},
		"Replies":                              {},
		"Replies.User":                         {},
		"Replies.Replies":                      {},
		"Replies.Replies.User":                 {},
		"Replies.Replies.Replies":              {},
		"Replies.Replies.Replies.User":         {},
		"Replies.Replies.Replies.Replies":      {},
		"Replies.Replies.Replies.Replies.User": {},
	})
	if err != nil {
		s.logger.Error("Failed to list comments", "error", err)
		return response.ResultErrDb(err)
	}

	// Get reaction counts and user reactions for all comments
	commentIds := extractAllCommentIds(comments)
	reactionCounts, err := s.reactionRepo.CountByCommentIds(ctx, commentIds)
	if err != nil {
		s.logger.Error("Failed to fetch reaction counts", "error", err)
		return response.ResultErrDb(err)
	}

	userReactions := make(map[uuid.UUID]string)
	user, err := utils.GetCurrentUserFormContext(ctx)
	if err == nil && user != nil {
		userReactions, err = s.reactionRepo.GetUserReactionsByCommentIds(ctx, commentIds, user.ID)
		if err != nil {
			s.logger.Error("Failed to fetch user reactions", "error", err)
			return response.ResultErrDb(err)
		}
	}

	// Map to response DTOs
	commentResponses := mapCommentsToResponses(comments, reactionCounts, userReactions)

	return response.ResultPaginationData(commentResponses, total, "Comments retrieved successfully")
}

func parseChapterID(raw string) (uuid.UUID, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return uuid.Nil, errors.New("chapterId is empty")
	}

	if strings.HasPrefix(trimmed, "[") {
		var ids []string
		if err := json.Unmarshal([]byte(trimmed), &ids); err != nil {
			return uuid.Nil, err
		}
		if len(ids) == 0 {
			return uuid.Nil, errors.New("chapterId array is empty")
		}
		trimmed = ids[0]
	}

	trimmed = strings.Trim(trimmed, `"`)
	return uuid.Parse(trimmed)
}

func extractAllCommentIds(comments []*model.Comment) []uuid.UUID {
	ids := make([]uuid.UUID, 0)
	var extract func([]*model.Comment)
	extract = func(cmts []*model.Comment) {
		for _, c := range cmts {
			ids = append(ids, c.ID)
			if c.Replies != nil {
				extract(c.Replies)
			}
		}
	}
	extract(comments)
	return ids
}
