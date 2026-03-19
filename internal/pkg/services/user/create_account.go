package userserivce

import (
	"base-go/internal/app/api/common/response"
	"base-go/internal/pkg/hash"
	"base-go/internal/pkg/model"
	userrequest "base-go/internal/pkg/request/user"
	"context"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *UserService) CreateAccount(ctx context.Context, req *userrequest.CreateUserRequest) response.Result {
	// validate request
	userExists, err := s.userRepo.FindOne(ctx, []any{
		clause.Eq{Column: "email", Value: req.Email},
	}, nil)
	if err != nil && !errors.As(err, &gorm.ErrRecordNotFound) {
		s.logger.Error("Failed to check if user exists", "error", err)
		return response.ResultErrDb(err)
	}

	if userExists != nil {
		return response.ResultError("User with this email already exists")
	}

	user := model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hash.HashPassword(req.Password),
	}
	err = s.userRepo.Create(ctx, &user)
	if err != nil {
		s.logger.Error("Failed to create user account", "error", err)
		return response.ResultErrDb(err)
	}

	return response.ResultSuccess("User account created successfully", user)
}
