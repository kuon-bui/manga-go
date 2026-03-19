package userserivce

import (
	"context"
	"errors"
	"fmt"
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/mail/mailable"
	"manga-go/internal/pkg/utils"
	"time"

	"gorm.io/gorm/clause"
)

func (s *UserService) RequestResetPassword(ctx context.Context, email string) response.Result {
	if email == "" {
		return response.ResultInvalidRequestErr(errors.New("email is required"))
	}
	user, err := s.userRepo.FindOne(ctx, []any{
		clause.Eq{
			Column: "email",
			Value:  email,
		},
	}, nil)
	if err != nil {
		s.logger.Error("Failed to find user by email: ", err)
		return response.ResultErrDb(err)
	}

	token := utils.TokenGenerator(30)

	expiryTimeAt := time.Now().Add(time.Duration(s.config.ResetPassword.TokenExpiryMinutes) * time.Minute)

	err = s.userRepo.Update(ctx, []any{
		clause.Eq{
			Column: "id",
			Value:  user.ID,
		},
	}, map[string]any{
		"reset_password_token":     token,
		"reset_password_expiry_at": expiryTimeAt,
	})
	if err != nil {
		s.logger.Error("Failed to set reset password token: ", err)
		return response.ResultErrDb(err)
	}

	mailable.NewResetPasswordMail(mailable.ResetPasswordMailParams{
		UserName:         user.Name,
		ResetPasswordURL: fmt.Sprintf(s.config.ResetPassword.ResetPasswordURL, token),
		ExpiryMinutes:    s.config.ResetPassword.TokenExpiryMinutes,
		CurrentYear:      time.Now().Year(),
	}).AddTo(user.Email).
		Dispatch(s.asynqClient)

	return response.ResultSuccess("request reset password success", nil)
}
