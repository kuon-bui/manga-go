package commentrequest

type AddReactionRequest struct {
	Type string `json:"type" binding:"required" example:"like"`
}
