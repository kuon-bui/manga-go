package authorizationroute

import (
	authorizationadmin "manga-go/internal/pkg/services/authorization_admin"

	"go.uber.org/fx"
)

type Handler struct {
	authorizationAdmin *authorizationadmin.Service
}

type HandlerParams struct {
	fx.In

	AuthorizationAdmin *authorizationadmin.Service
}

func NewHandler(params HandlerParams) *Handler {
	return &Handler{authorizationAdmin: params.AuthorizationAdmin}
}
