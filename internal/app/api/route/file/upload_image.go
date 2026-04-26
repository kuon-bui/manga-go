package fileroute

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/constant"
	"manga-go/internal/pkg/fileprocess"
	filerequest "manga-go/internal/pkg/request/file"
	queueconstant "manga-go/internal/queue/queue_constant"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

const maxUploadImageSize int64 = 10 * 1024 * 1024
const defaultTemporaryImageCleanupDelayHours = 24

// @Summary      Upload image (chapter or comic cover)
// @Description  Queue image processing to generate WebP variants (small/medium/large/normal). For chapter type: if chapterSlug provided -> comics/{comicSlug}/chapters/{chapterSlug}/pages/page-{pageIdx}.webp. For cover -> comics/{comicSlug}/cover/{uuid}.webp
// @Tags         File
// @Accept       multipart/form-data
// @Produce      json
// @Param        file  formData  file  true  "Image file (max 10MB, image/* only)"
// @Param        type  formData  string  true  "Image type: 'chapter' or 'cover'"
// @Param        comicId  formData  string  true  "Comic ID (UUID)"
// @Param        chapterSlug  formData  string  false  "Chapter slug (required for chapter type)"
// @Param        pageIdx  formData  string  false  "Page index (required for chapter type)"
// @Success      200   {object}  response.Response
// @Failure      400   {object}  response.Response
// @Failure      401   {object}  response.Response
// @Failure      413   {object}  response.Response
// @Failure      500   {object}  response.Response
// @Router       /files/upload [post]
// @Security     AccessToken
func (h *FileHandler) uploadImage(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxUploadImageSize)

	var req filerequest.UploadFileRequest
	if err := c.ShouldBind(&req); err != nil {
		if strings.Contains(err.Error(), "http: request body too large") {
			response.ResultError("File size exceeds 10MB").ResponseResult(c)
			return
		}
		response.ResultInvalidRequestErr(err).ResponseResult(c)
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		if strings.Contains(err.Error(), "http: request body too large") {
			response.ResultError("File size exceeds 10MB").ResponseResult(c)
			return
		}
		response.ResultInvalidRequestErr(err).ResponseResult(c)
		return
	}
	req.File = *fileHeader

	if req.File.Size > maxUploadImageSize {
		response.ResultError("File size exceeds 10MB").ResponseResult(c)
		return
	}

	file, err := req.File.Open()
	if err != nil {
		response.ResultErrInternal(err).ResponseResult(c)
		return
	}
	defer file.Close()

	contentType := req.File.Header.Get("Content-Type")
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

	if req.Type != constant.UploadImageTypeChapter && req.Type != constant.UploadImageTypeComic && !strings.EqualFold(string(req.Type), "cover") {
		response.ResultError("'type' must be 'chapter' or 'cover'").ResponseResult(c)
		return
	}

	comicIdStr := req.ComicId.String()

	// Generate canonical WebP filename with UUID
	uniqueFilename := uuid.New().String() + ".webp"

	// Resolve slugs from IDs via fileService
	var filePath string
	switch req.Type {
	case constant.UploadImageTypeChapter:
		// For chapter type: chapterSlug must be provided.
		//
		// save to comics/{comicSlug}/chapters/{chapterSlug}/pages/

		if req.ChapterSlug == nil || strings.TrimSpace(*req.ChapterSlug) == "" {
			response.ResultError("'chapterSlug' parameter is required for chapter type").ResponseResult(c)
			return
		}

		if req.PageIdx == nil {
			response.ResultError("'pageIdx' parameter is required for chapter type").ResponseResult(c)
			return
		}

		chapterSlug := strings.TrimSpace(*req.ChapterSlug)

		if strings.Contains(chapterSlug, "/") || strings.Contains(chapterSlug, "\\") || strings.Contains(chapterSlug, "..") {
			response.ResultError("'chapterSlug' contains invalid characters").ResponseResult(c)
			return
		}

		if *req.PageIdx < 0 {
			response.ResultError("'pageIdx' must be a non-negative integer").ResponseResult(c)
			return
		}

		uniqueFilename = fmt.Sprintf("page-%d.webp", *req.PageIdx)
		// Backend resolves: comicId -> comicSlug
		path, err := h.fileService.BuildChapterImagePath(c.Request.Context(), comicIdStr, chapterSlug, uniqueFilename)
		if err != nil {
			response.ResultError(err.Error()).ResponseResult(c)
			return
		}
		filePath = path

	case constant.UploadImageTypeComic:
		// Build path: comics/{comicSlug}/cover/{uuid}.webp
		path, err := h.fileService.BuildCoverImagePath(c.Request.Context(), comicIdStr, uniqueFilename)
		if err != nil {
			response.ResultError(err.Error()).ResponseResult(c)
			return
		}
		filePath = path

	default:
		// Backward-compatible alias: accept cover as comic upload type.
		if strings.EqualFold(string(req.Type), "cover") {
			path, err := h.fileService.BuildCoverImagePath(c.Request.Context(), comicIdStr, uniqueFilename)
			if err != nil {
				response.ResultError(err.Error()).ResponseResult(c)
				return
			}
			filePath = path
		}
	}

	temporaryObjectKey := "tmp/image-process/" + uuid.NewString()
	if err := h.fileService.UploadFile(c.Request.Context(), temporaryObjectKey, file, req.File.Size, contentType); err != nil {
		response.ResultErrInternal(err).ResponseResult(c)
		return
	}

	payloadBytes, err := json.Marshal(fileprocess.ImageProcessPayload{
		FilePath:           filePath,
		TemporaryObjectKey: temporaryObjectKey,
	})
	if err != nil {
		_ = h.fileService.DeleteFile(c.Request.Context(), temporaryObjectKey)
		response.ResultErrInternal(err).ResponseResult(c)
		return
	}

	task := asynq.NewTask(
		queueconstant.IMAGE_PROCESS_TASK,
		payloadBytes,
		asynq.MaxRetry(5),
		asynq.Timeout(3*time.Minute),
	)

	taskInfo, err := h.asynqClient.Enqueue(task, asynq.Queue(queueconstant.IMAGE_PROCESS_QUEUE))
	if err != nil {
		_ = h.fileService.DeleteFile(c.Request.Context(), temporaryObjectKey)
		response.ResultErrInternal(err).ResponseResult(c)
		return
	}

	cleanupTaskScheduled := false
	cleanupTaskID := ""
	cleanupPayloadBytes, cleanupPayloadErr := json.Marshal(fileprocess.ImageProcessCleanupPayload{
		TemporaryObjectKey: temporaryObjectKey,
	})
	if cleanupPayloadErr == nil {
		cleanupDelayHours := h.config.Asynq.ImageProcessCleanupDelayHours
		if cleanupDelayHours <= 0 {
			cleanupDelayHours = defaultTemporaryImageCleanupDelayHours
		}

		cleanupTask := asynq.NewTask(
			queueconstant.IMAGE_PROCESS_CLEANUP_TASK,
			cleanupPayloadBytes,
			asynq.ProcessIn(time.Duration(cleanupDelayHours)*time.Hour),
			asynq.MaxRetry(3),
		)

		cleanupTaskInfo, cleanupEnqueueErr := h.asynqClient.Enqueue(cleanupTask, asynq.Queue(queueconstant.IMAGE_PROCESS_QUEUE))
		if cleanupEnqueueErr == nil {
			cleanupTaskScheduled = true
			cleanupTaskID = cleanupTaskInfo.ID
		}
	}

	response.ResultSuccess("Upload image queued successfully", map[string]any{
		"status":                 "queued",
		"taskId":                 taskInfo.ID,
		"cleanup_task_scheduled": cleanupTaskScheduled,
		"cleanup_task_id":        cleanupTaskID,
		"path":                   filePath,
		"filename":               uniqueFilename,
		"url":                    "/files/content/" + filePath,
		"content_type":           "image/webp",
	}).ResponseResult(c)
}
