package authorrequest

type UpdateAuthorRequest struct {
	Name string `json:"name" binding:"required"`
}
