package userroute

import (
	"manga-go/internal/pkg/config"
	userserivce "manga-go/internal/pkg/services/user"

	"go.uber.org/fx"
)

type userHandler struct {
	config      *config.Config
	userService *userserivce.UserService
}

type UserHandlerParams struct {
	fx.In

	Config      *config.Config
	UserService *userserivce.UserService
}

func NewUserHandler(p UserHandlerParams) *userHandler {
	return &userHandler{
		userService: p.UserService,
		config:      p.Config,
	}
}
