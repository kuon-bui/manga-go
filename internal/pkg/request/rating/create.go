package ratingrequest

type CreateRatingRequest struct {
	Score int `json:"score" binding:"required,min=1,max=5"`
}
