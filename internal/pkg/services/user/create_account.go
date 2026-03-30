package userservice

import (
	"context"
	"errors"
	"manga-go/internal/app/api/common/response"
	casbinpkg "manga-go/internal/pkg/casbin"
	"manga-go/internal/pkg/hash"
	"manga-go/internal/pkg/model"
	userrequest "manga-go/internal/pkg/request/user"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *UserService) CreateAccount(ctx context.Context, req *userrequest.CreateUserRequest) response.Result {
	// validate request
	userExists, err := s.userRepo.FindOne(ctx, []any{
		clause.Eq{Column: "email", Value: req.Email},
	}, nil)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
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

	// Assign the global "user" role in Casbin so the user can create/join groups
	if _, err := s.enforcer.AddRoleForUserInDomain(user.ID.String(), "user", casbinpkg.GlobalDomain); err != nil {
		s.logger.Errorf("Failed to assign user role for %s: %v", user.ID, err)
		return response.ResultErrInternal(err)
	}

	return response.ResultSuccess("User account created successfully", user)
}
