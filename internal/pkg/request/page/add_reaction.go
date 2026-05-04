package pagerequest

type AddReactionRequest struct {
	Type string `json:"type" binding:"required" example:"LIKE"`
}
