package authorrequest

type CreateAuthorRequest struct {
	Name string `json:"name" binding:"required"`
}
