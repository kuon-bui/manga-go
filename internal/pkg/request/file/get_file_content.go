package filerequest

type GetFileContentRequest struct {
	Variant string `form:"variant" binding:"omitempty,oneof=economy small clear sharp"`
}
