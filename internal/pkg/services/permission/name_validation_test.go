package permissionservice

import (
	"context"
	"net/http"
	"testing"

	permissionrequest "manga-go/internal/pkg/request/permission"
)

func TestCreatePermissionRejectsMalformedName(t *testing.T) {
	s := newPermissionService(t, true)

	result := s.CreatePermission(context.Background(), &permissionrequest.CreatePermissionRequest{
		Name: "comics.read",
	})

	if result.HttpStatus != http.StatusBadRequest {
		t.Fatalf("expected 400 for a name the policy engine cannot parse, got %d", result.HttpStatus)
	}
}

func TestCreatePermissionAcceptsWellFormedName(t *testing.T) {
	s := newPermissionService(t, true)

	result := s.CreatePermission(context.Background(), &permissionrequest.CreatePermissionRequest{
		Name: "comic:read",
	})

	if result.HttpStatus != http.StatusOK {
		t.Fatalf("expected a well-formed name to be accepted, got %d", result.HttpStatus)
	}
}

func TestUpdatePermissionRejectsMalformedName(t *testing.T) {
	s := newPermissionService(t, true)

	created := s.CreatePermission(context.Background(), &permissionrequest.CreatePermissionRequest{
		Name: "comic:read",
	})
	if created.HttpStatus != http.StatusOK {
		t.Fatalf("setup failed: could not create permission, got %d", created.HttpStatus)
	}

	permissions, _, err := s.permissionRepo.FindPaginated(context.Background(), nil, nil, nil)
	if err != nil || len(permissions) == 0 {
		t.Fatalf("setup failed: could not read back the created permission: %v", err)
	}
	id := permissions[0].ID

	result := s.UpdatePermission(context.Background(), id, &permissionrequest.UpdatePermissionRequest{
		Name: "manage-everything",
	})

	if result.HttpStatus != http.StatusBadRequest {
		t.Fatalf("expected 400 for a name the policy engine cannot parse, got %d", result.HttpStatus)
	}
}
