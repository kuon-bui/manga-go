package notificationservice

import (
	"context"
	"errors"
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/model"
	userrequest "manga-go/internal/pkg/request/user"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *NotificationService) GetUserConfig(ctx context.Context, userID uuid.UUID) response.Result {
	user, err := s.userRepo.FindOne(ctx, []any{
		clause.Eq{Column: "id", Value: userID},
	}, nil)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.ResultNotFound("User")
		}

		s.logger.Error("Failed to load user config", "error", err)
		return response.ResultErrDb(err)
	}

	return response.ResultSuccess("User config retrieved successfully", user.UserConfig.ToResponse())
}

func (s *NotificationService) UpdateUserConfig(ctx context.Context, userID uuid.UUID, req *userrequest.UpdateUserConfigRequest) response.Result {
	user, err := s.userRepo.FindOne(ctx, []any{
		clause.Eq{Column: "id", Value: userID},
	}, nil)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.ResultNotFound("User")
		}

		s.logger.Error("Failed to find user for config update", "error", err)
		return response.ResultErrDb(err)
	}

	config := user.UserConfig
	if len(config) == 0 {
		config = model.DefaultUserConfig()
	}

	if req.SeenNotificationCenter != nil {
		config.Set(model.UserConfigSeenNotificationCenter, *req.SeenNotificationCenter)
	}
	if req.EnableSSENotifications != nil {
		config.Set(model.UserConfigEnableSSENotifications, *req.EnableSSENotifications)
	}
	if req.EnableEmailNotifications != nil {
		config.Set(model.UserConfigEnableEmailNotifications, *req.EnableEmailNotifications)
	}
	if req.EnableComicNewChapterNotifications != nil {
		config.Set(model.UserConfigEnableComicNewChapterNotifications, *req.EnableComicNewChapterNotifications)
	}
	if req.EnableCommentReplyNotifications != nil {
		config.Set(model.UserConfigEnableCommentReplyNotifications, *req.EnableCommentReplyNotifications)
	}
	if req.EnableMentionNotifications != nil {
		config.Set(model.UserConfigEnableMentionNotifications, *req.EnableMentionNotifications)
	}
	if req.EnableSystemAnnouncements != nil {
		config.Set(model.UserConfigEnableSystemAnnouncements, *req.EnableSystemAnnouncements)
	}

	if err := s.userRepo.Update(ctx, []any{
		clause.Eq{Column: "id", Value: userID},
	}, map[string]any{
		"user_config": config,
	}); err != nil {
		s.logger.Error("Failed to update user config", "error", err)
		return response.ResultErrDb(err)
	}

	return response.ResultSuccess("User config updated successfully", config.ToResponse())
}
