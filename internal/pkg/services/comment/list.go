package commentservice

import (
	"context"
	"encoding/json"
	"errors"
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/common"
	commentrequest "manga-go/internal/pkg/request/comment"
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

	return response.ResultPaginationData(comments, total, "Comments retrieved successfully")
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
