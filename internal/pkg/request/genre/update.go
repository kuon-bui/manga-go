package genrerequest

type UpdateGenreRequest struct {
	Name string `json:"name" binding:"required"`
}
