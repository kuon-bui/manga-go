package userroute

import (
	"manga-go/internal/pkg/config"
	comicservice "manga-go/internal/pkg/services/comic"
	notificationservice "manga-go/internal/pkg/services/notification"
	userserivce "manga-go/internal/pkg/services/user"

	"go.uber.org/fx"
)

type userHandler struct {
	comicService        *comicservice.ComicService
	config              *config.Config
	notificationService *notificationservice.NotificationService
	userService         *userserivce.UserService
}

type UserHandlerParams struct {
	fx.In

	ComicService        *comicservice.ComicService
	Config              *config.Config
	NotificationService *notificationservice.NotificationService
	UserService         *userserivce.UserService
}

func NewUserHandler(p UserHandlerParams) *userHandler {
	return &userHandler{
		comicService:        p.ComicService,
		userService:         p.UserService,
		config:              p.Config,
		notificationService: p.NotificationService,
	}
}
