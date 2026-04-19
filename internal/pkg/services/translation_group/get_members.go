package translationgroupservice

import (
	"context"
	"manga-go/internal/app/api/common/response"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm/clause"
)

type GroupMemberResponse struct {
	ID       string  `json:"id"`
	UserID   string  `json:"userId"`
	User     UserObj `json:"user"`
	Role     string  `json:"role"`
	JoinedAt string  `json:"joinedAt"`
}

type UserObj struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	AvatarUrl *string `json:"avatarUrl"`
}

func (s *TranslationGroupService) GetMembers(ctx context.Context, groupID uuid.UUID) response.Result {
	group, err := s.translationGroupRepo.FindOne(ctx, []any{
		clause.Eq{Column: "id", Value: groupID},
	}, nil)
	if err != nil {
		s.logger.Error("Failed to find translation group", "error", err)
		return response.ResultErrDb(err)
	}

	users, err := s.userRepo.FindAll(ctx, []any{
		clause.Eq{Column: "translation_group_id", Value: groupID},
	}, nil)
	if err != nil {
		s.logger.Error("Failed to list translation group members", "error", err)
		return response.ResultErrDb(err)
	}

	var res []GroupMemberResponse
	for _, u := range users {
		role := "member"
		if u.ID == group.OwnerID {
			role = "admin"
		}

		res = append(res, GroupMemberResponse{
			ID:       u.ID.String(),
			UserID:   u.ID.String(),
			User: UserObj{
				ID:        u.ID.String(),
				Name:      u.Name,
				AvatarUrl: nil, // Note: Users currently don't have AvatarUrl in DB
			},
			Role:     role,
			JoinedAt: u.CreatedAt.Format(time.RFC3339),
		})
	}

	return response.ResultSuccess("Members retrieved successfully", map[string]interface{}{
		"data":  res,
		"total": len(res),
	})
}
