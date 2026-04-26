package fileroute

import (
	"io"
	"net/http"
	"strings"

	"manga-go/internal/app/api/common/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const maxUploadImageSize int64 = 10 * 1024 * 1024

// @Summary      Upload image (chapter or comic cover)
// @Description  Upload image and automatically generate WebP variants (economy/small/clear/sharp). For chapter type: if chapterId provided -> comics/{comicSlug}/chapters/{chapterSlug}/pages/{uuid}.webp, else -> comics/{comicSlug}/temp-uploads/{uuid}.webp. For cover -> comics/{comicSlug}/cover/{uuid}.webp
// @Tags         File
// @Accept       multipart/form-data
// @Produce      json
// @Param        file  formData  file  true  "Image file (max 10MB, image/* only)"
// @Param        type  formData  string  true  "Image type: 'chapter' or 'cover'"
// @Param        comicId  formData  string  true  "Comic ID (UUID)"
// @Param        chapterId  formData  string  false  "Chapter ID (UUID) - optional. If not provided, image saved to temp folder for later assignment"
// @Success      200   {object}  response.Response
// @Failure      400   {object}  response.Response
// @Failure      401   {object}  response.Response
// @Failure      413   {object}  response.Response
// @Failure      500   {object}  response.Response
// @Router       /files/upload [post]
// @Security     AccessToken
func (h *FileHandler) uploadImage(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxUploadImageSize)

	fileHeader, err := c.FormFile("file")
	if err != nil {
		if strings.Contains(err.Error(), "http: request body too large") {
			response.ResultError("File size exceeds 10MB").ResponseResult(c)
			return
		}
		response.ResultInvalidRequestErr(err).ResponseResult(c)
		return
	}

	if fileHeader.Size > maxUploadImageSize {
		response.ResultError("File size exceeds 10MB").ResponseResult(c)
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		response.ResultErrInternal(err).ResponseResult(c)
		return
	}
	defer file.Close()

	contentType := fileHeader.Header.Get("Content-Type")
	if contentType == "" {
		header := make([]byte, 512)
		readN, readErr := file.Read(header)
		if readErr != nil && readErr != io.EOF {
			response.ResultErrInternal(readErr).ResponseResult(c)
			return
		}

		contentType = http.DetectContentType(header[:readN])
		if _, seekErr := file.Seek(0, io.SeekStart); seekErr != nil {
			response.ResultErrInternal(seekErr).ResponseResult(c)
			return
		}
	}

	if !strings.HasPrefix(contentType, "image/") {
		response.ResultError("Only image files are allowed").ResponseResult(c)
		return
	}

	// Get parameters from form
	uploadType := strings.TrimSpace(c.PostForm("type"))
	comicIdStr := strings.TrimSpace(c.PostForm("comicId"))
	chapterIdStr := strings.TrimSpace(c.PostForm("chapterId"))

	if uploadType == "" {
		response.ResultError("'type' parameter is required (chapter or cover)").ResponseResult(c)
		return
	}

	if comicIdStr == "" {
		response.ResultError("'comicId' parameter is required").ResponseResult(c)
		return
	}

	if uploadType != "chapter" && uploadType != "cover" {
		response.ResultError("'type' must be 'chapter' or 'cover'").ResponseResult(c)
		return
	}

	// For chapter type: chapterId is optional
	// - If provided: save to comics/{comicSlug}/chapters/{chapterSlug}/pages/
	// - If NOT provided: save to comics/{comicSlug}/temp-uploads/ (pending assignment)

	// Generate canonical WebP filename with UUID
	uniqueFilename := uuid.New().String() + ".webp"

	// Resolve slugs from IDs via fileService
	var filePath string
	if uploadType == "chapter" {
		if chapterIdStr != "" {
			// Backend resolves: comicId -> comicSlug, chapterId -> chapterSlug
			path, err := h.fileService.BuildChapterImagePath(c.Request.Context(), comicIdStr, chapterIdStr, uniqueFilename)
			if err != nil {
				response.ResultError(err.Error()).ResponseResult(c)
				return
			}
			filePath = path
		} else {
			// No chapterId provided: save to temp-uploads folder
			// Frontend will assign these images to chapter when creating
			path, err := h.fileService.BuildTempChapterImagePath(c.Request.Context(), comicIdStr, uniqueFilename)
			if err != nil {
				response.ResultError(err.Error()).ResponseResult(c)
				return
			}
			filePath = path
		}
	} else if uploadType == "cover" {
		// Build path: comics/{comicSlug}/cover/{uuid}.webp
		path, err := h.fileService.BuildCoverImagePath(c.Request.Context(), comicIdStr, uniqueFilename)
		if err != nil {
			response.ResultError(err.Error()).ResponseResult(c)
			return
		}
		filePath = path
	}

	uploadResult, err := h.fileService.UploadImageVariants(c.Request.Context(), filePath, file)
	if err != nil {
		response.ResultErrInternal(err).ResponseResult(c)
		return
	}

	response.ResultSuccess("Upload image successfully", uploadResult).ResponseResult(c)
}
