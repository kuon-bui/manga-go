package permissionroute

import (
	permissionservice "manga-go/internal/pkg/services/permission"

	"go.uber.org/fx"
)

type PermissionHandler struct {
	permissionService *permissionservice.PermissionService
}

type PermissionHandlerParams struct {
	fx.In

	PermissionService *permissionservice.PermissionService
}

func NewPermissionHandler(p PermissionHandlerParams) *PermissionHandler {
	return &PermissionHandler{
		permissionService: p.PermissionService,
	}
}
