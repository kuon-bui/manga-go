package commentrequest

type UpdateCommentRequest struct {
	Content   string `json:"content" binding:"required"`
	PageIndex *int   `json:"pageIndex,omitempty"`
}
