package roleroute

import (
	roleservice "manga-go/internal/pkg/services/role"

	"go.uber.org/fx"
)

type RoleHandler struct {
	roleService *roleservice.RoleService
}

type RoleHandlerParams struct {
	fx.In

	RoleService *roleservice.RoleService
}

func NewRoleHandler(p RoleHandlerParams) *RoleHandler {
	return &RoleHandler{
		roleService: p.RoleService,
	}
}
