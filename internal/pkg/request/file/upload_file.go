package filerequest

import (
	"manga-go/internal/pkg/constant"
	"mime/multipart"

	"github.com/google/uuid"
)

type UploadFileRequest struct {
	Type        constant.UploadImageType `form:"type" binding:"required,upload_image_type"`
	ComicId     uuid.UUID                `form:"comicId" binding:"required,uuid"`
	ChapterSlug *string                  `form:"chapterSlug"`
	PageIdx     *int                     `form:"pageIdx"`

	File multipart.FileHeader `form:"file" binding:"required"` // 10MB max
}
