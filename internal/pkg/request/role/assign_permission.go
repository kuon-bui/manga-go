package rolerequest

// AssignPermissionRequest replaces the role's grants with exactly Permissions.
// Names come from the catalog served by GET /permissions, e.g. "comic:write".
type AssignPermissionRequest struct {
	Permissions []string `json:"permissions" binding:"required,min=1,dive,required"`
}
