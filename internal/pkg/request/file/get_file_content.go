package filerequest

type GetFileContentRequest struct {
	Size string `form:"size" binding:"omitempty,oneof=small medium large normal"`
}
