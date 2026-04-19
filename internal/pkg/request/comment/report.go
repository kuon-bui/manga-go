package commentrequest

type ReportCommentRequest struct {
	Reason  string  `json:"reason" binding:"required,oneof=SPAM OFFENSIVE HARASSMENT ADULT_CONTENT"`
	Details *string `json:"details" binding:"max=500"`
}
