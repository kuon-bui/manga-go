package ratingrequest

type UpdateRatingRequest struct {
	Score int `json:"score" binding:"required,min=1,max=5"`
}
