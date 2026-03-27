package userservice

import (
	"context"
	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/hash"
	jwtprovider "manga-go/internal/pkg/jwt_provider"
	userrequest "manga-go/internal/pkg/request/user"

	"gorm.io/gorm/clause"
)

func (us *UserService) SignIn(ctx context.Context, req *userrequest.SignInRequest) (*jwtprovider.Token, *jwtprovider.Token, response.Result) {
	user, err := us.userRepo.FindOne(ctx, []any{
		clause.Eq{Column: "email", Value: req.Email},
	}, nil)
	if err != nil {
		us.logger.Error("Failed to find user by email", "error", err)
		return nil, nil, response.ResultErrDb(err)
	}

	if !hash.ComparePassword(user.Password, req.Password) {
		return nil, nil, response.ResultError("Invalid email or password")
	}

	userPayload := jwtprovider.UserPayload{
		UserID:   user.ID,
		FullName: user.Name,
		Email:    user.Email,
	}

	accessToken, refreshToken, err := us.jwtProvider.GenerateToken(userPayload)
	if err != nil {
		us.logger.Error("Failed to generate JWT tokens", "error", err)
		return nil, nil, response.ResultErrInternal(err)
	}

	return accessToken, refreshToken, response.ResultSuccess("Sign in successful", map[string]any{
		"user": user,
	})
}
