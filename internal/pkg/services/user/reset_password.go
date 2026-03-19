package userserivce

import (
	"context"
	"errors"
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/hash"
	userrequest "manga-go/internal/pkg/request/user"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *UserService) ResetPassword(ctx context.Context, req userrequest.ResetPasswordRequest) response.Result {
	now := time.Now()
	user, err := s.userRepo.FindOne(ctx, []any{
		// find user by reset password token and check if token is not expired
		clause.And(
			clause.Eq{
				Column: "reset_password_token",
				Value:  req.Token,
			},
			clause.Gt{
				Column: "reset_password_expiry_at",
				Value:  now,
			},
		),
	}, nil)
	if err != nil {
		s.logger.Error("Failed to find user by reset password token: ", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.ResultInvalidRequestErr(errors.New("invalid or expired reset password token"))
		}

		return response.ResultErrDb(err)
	}

	newPasswordHash := hash.HashPassword(req.NewPassword)
	if newPasswordHash == user.Password {
		return response.ResultInvalidRequestErr(errors.New("new password must be different from the old password"))
	}
	err = s.userRepo.Update(ctx, []any{
		clause.Eq{
			Column: "id",
			Value:  user.ID,
		},
	}, map[string]any{
		"password":                 newPasswordHash,
		"reset_password_token":     nil,
		"reset_password_expiry_at": nil,
	})
	if err != nil {
		s.logger.Error("Failed to reset password: ", err)
		return response.ResultErrDb(err)
	}

	return response.ResultSuccess("reset password success", nil)
}
