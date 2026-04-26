package userrequest

type UpdateUserProfileRequest struct {
	Name   *string `json:"name"`
	Avatar *string `json:"avatar"`
}
