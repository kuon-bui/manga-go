package translationgrouprequest

type JoinTranslationGroupRequest struct {
	Slug string `json:"slug" binding:"required"`
}
