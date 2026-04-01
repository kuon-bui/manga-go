package fileroute

import (
	"io"
	"net/http"
	"strings"

	"manga-go/internal/app/api/common/response"

	"github.com/gin-gonic/gin"
)

const maxUploadImageSize int64 = 10 * 1024 * 1024

// @Summary      Upload file
// @Description  Upload a file to object storage
// @Tags         File
// @Accept       multipart/form-data
// @Produce      json
// @Param        file  formData  file  true  "File to upload (max 10MB)"
// @Success      200   {object}  response.Response
// @Failure      400   {object}  response.Response
// @Failure      401   {object}  response.Response
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

	filename := strings.TrimPrefix(c.PostForm("filename"), "/")
	if filename == "" {
		filename = strings.TrimPrefix(fileHeader.Filename, "/")
	}
	if filename == "" {
		response.ResultError("Invalid filename").ResponseResult(c)
		return
	}

	err = h.fileService.UploadFile(c.Request.Context(), filename, file, fileHeader.Size, contentType)
	if err != nil {
		response.ResultErrInternal(err).ResponseResult(c)
		return
	}

	response.ResultSuccess("Upload image successfully", map[string]any{
		"filename":     filename,
		"content_type": contentType,
		"size":         fileHeader.Size,
	}).ResponseResult(c)
}
