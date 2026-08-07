package authorization

import (
	"testing"
)

func TestCatalogIsNotEmpty(t *testing.T) {
	if len(Catalog()) == 0 {
		t.Fatal("expected the catalog to define the grantable permissions")
	}
}

func TestCatalogNamesAreUnique(t *testing.T) {
	seen := map[string]bool{}
	for _, def := range Catalog() {
		if seen[def.Name] {
			t.Errorf("permission %q is defined twice", def.Name)
		}
		seen[def.Name] = true
	}
}

func TestCatalogEntriesExpandToConcreteActions(t *testing.T) {
	for _, def := range Catalog() {
		if def.Object == "" {
			t.Errorf("%q has no object", def.Name)
		}
		if len(def.Grants) == 0 {
			t.Errorf("%q grants no actions, so it could never authorize anything", def.Name)
		}
		if def.Name != string(def.Object)+":"+string(def.Action) {
			t.Errorf("%q does not match its object/action pair (%s/%s)", def.Name, def.Object, def.Action)
		}
	}
}

func TestLookupPermissionExpandsWriteShorthand(t *testing.T) {
	def, ok := LookupPermission("comic:write")
	if !ok {
		t.Fatal("expected comic:write to be a known permission")
	}

	want := map[Action]bool{ActionCreate: true, ActionUpdate: true, ActionPublish: true}
	if len(def.Grants) != len(want) {
		t.Fatalf("expected write to expand to create/update/publish, got %v", def.Grants)
	}
	for _, action := range def.Grants {
		if !want[action] {
			t.Errorf("unexpected action %q in the expansion of comic:write", action)
		}
	}
}

func TestLookupPermissionKeepsSingleActionsAsThemselves(t *testing.T) {
	def, ok := LookupPermission("role:manage")
	if !ok {
		t.Fatal("expected role:manage to be a known permission")
	}
	if len(def.Grants) != 1 || def.Grants[0] != ActionManage {
		t.Fatalf("expected role:manage to grant exactly manage, got %v", def.Grants)
	}
}

func TestLookupPermissionRejectsUnknownNames(t *testing.T) {
	// Structurally malformed, plus names that parse but name nothing real. The
	// second group is the important one: the old parser accepted any action
	// string and produced a policy rule that could never match a request.
	for _, name := range []string{
		"", "comic", "comic:", ":read", "comic:read:extra", "comics.read",
		"comic:frobnicate", "unicorn:read", "COMIC:READ",
	} {
		if _, ok := LookupPermission(name); ok {
			t.Errorf("expected %q to be rejected, but the catalog accepted it", name)
		}
	}
}

func TestValidatePermissionNameAgreesWithCatalog(t *testing.T) {
	for _, def := range Catalog() {
		if err := ValidatePermissionName(def.Name); err != nil {
			t.Errorf("catalog entry %q failed validation: %v", def.Name, err)
		}
	}
	if err := ValidatePermissionName("comic:frobnicate"); err == nil {
		t.Error("expected an unknown action to fail validation")
	}
}

// The seeded roles must be expressible in the catalog, otherwise seeding a
// fresh environment fails.
func TestCatalogCoversSeededPermissionNames(t *testing.T) {
	for _, name := range []string{
		"comic:read", "comic:write", "comic:delete",
		"chapter:read", "chapter:write", "chapter:delete",
		"user:read", "user:manage",
		"role:manage",
		"tag:write", "tag:delete",
		"genre:write", "genre:delete",
		"author:write", "author:delete",
		"translation_group:write", "translation_group:delete",
		"comment:delete",
		"rating:delete",
	} {
		if _, ok := LookupPermission(name); !ok {
			t.Errorf("seeded permission %q is missing from the catalog", name)
		}
	}
}

func TestCatalogContainsAuditLogRead(t *testing.T) {
	definition, ok := LookupPermission("audit_log:read")
	if !ok {
		t.Fatal("expected audit_log:read in the catalog")
	}
	if definition.Object != ObjectAuditLog || definition.Action != ActionRead {
		t.Fatalf("unexpected definition: %#v", definition)
	}
}
