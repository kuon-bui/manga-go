package authorization

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

// PermissionDefinition is one grantable permission. Name is the wire format
// ("<object>:<action>"); Grants is what it expands to when written as policy —
// the shorthand "write" covers create, update and publish.
type PermissionDefinition struct {
	Name     string    `json:"name"`
	Object   Object    `json:"object"`
	Action   Action    `json:"action"`
	Grants   []Action  `json:"grants"`
	Contexts []Context `json:"contexts"`
}

// actionWrite is a grant keyword rather than a request action: no request ever
// arrives with act=write, it is only ever a shorthand in a permission name.
const actionWrite Action = "write"

// grantableActions lists, per object, the permissions an administrator can hand
// out. It is the single source of truth for the permission vocabulary: anything
// absent here cannot be granted, which is what stops a typo from being stored
// as a policy rule that can never match a request.
var grantableActions = map[Object][]Action{
	ObjectAuditLog:         {ActionRead},
	ObjectAuthor:           {ActionRead, actionWrite, ActionDelete, ActionManage},
	ObjectChapter:          {ActionRead, actionWrite, ActionDelete, ActionManage},
	ObjectComic:            {ActionRead, actionWrite, ActionDelete, ActionManage},
	ObjectComment:          {ActionRead, actionWrite, ActionDelete, ActionManage},
	ObjectFile:             {ActionRead, actionWrite, ActionDelete, ActionManage},
	ObjectGenre:            {ActionRead, actionWrite, ActionDelete, ActionManage},
	ObjectNotification:     {ActionRead, actionWrite, ActionDelete, ActionManage},
	ObjectPage:             {ActionRead, actionWrite, ActionDelete, ActionManage},
	ObjectPermission:       {ActionRead, ActionManage},
	ObjectRating:           {ActionRead, actionWrite, ActionDelete, ActionManage},
	ObjectReadingHistory:   {ActionRead, actionWrite, ActionDelete, ActionManage},
	ObjectRole:             {ActionRead, ActionManage},
	ObjectTag:              {ActionRead, actionWrite, ActionDelete, ActionManage},
	ObjectTranslationGroup: {ActionRead, actionWrite, ActionDelete, ActionManage},
	ObjectUser:             {ActionRead, actionWrite, ActionDelete, ActionManage},
}

var (
	catalogOnce sync.Once
	catalogList []PermissionDefinition
	catalogByID map[string]PermissionDefinition
)

func buildCatalog() {
	catalogByID = make(map[string]PermissionDefinition)

	for object, actions := range grantableActions {
		for _, action := range actions {
			def := PermissionDefinition{
				Name:     fmt.Sprintf("%s:%s", object, action),
				Object:   object,
				Action:   action,
				Grants:   expandAction(action),
				Contexts: DefaultContexts(),
			}
			catalogByID[def.Name] = def
			catalogList = append(catalogList, def)
		}
	}

	sort.Slice(catalogList, func(i, j int) bool {
		return catalogList[i].Name < catalogList[j].Name
	})
}

func expandAction(action Action) []Action {
	if action == actionWrite {
		return []Action{ActionCreate, ActionUpdate, ActionPublish}
	}
	return []Action{action}
}

// Catalog returns every grantable permission, sorted by name.
func Catalog() []PermissionDefinition {
	catalogOnce.Do(buildCatalog)

	out := make([]PermissionDefinition, len(catalogList))
	copy(out, catalogList)
	return out
}

// LookupPermission resolves a permission name. Names are matched exactly:
// unknown objects, unknown actions and malformed strings are all rejected.
func LookupPermission(name string) (PermissionDefinition, bool) {
	catalogOnce.Do(buildCatalog)

	if strings.Count(name, ":") != 1 {
		return PermissionDefinition{}, false
	}
	def, ok := catalogByID[name]
	return def, ok
}

// ValidatePermissionName reports whether name is a permission this system can
// translate into policy rules.
func ValidatePermissionName(name string) error {
	if _, ok := LookupPermission(name); !ok {
		return fmt.Errorf("%w: %q is not a known permission", ErrInvalidPermissionName, name)
	}
	return nil
}
