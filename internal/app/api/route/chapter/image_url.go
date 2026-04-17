package chapterhandler

import (
	"net/url"
	"strings"

	"manga-go/internal/app/api/common/response"
	"manga-go/internal/pkg/model"

	"github.com/gin-gonic/gin"
)

func (h *ChapterHandler) normalizeChapterImageURLs(c *gin.Context, result *response.Result) {
	if result == nil || !result.Success || result.Data == nil {
		return
	}

	switch chapter := result.Data.(type) {
	case *model.Chapter:
		h.normalizePageImageURLs(c, chapter.Pages)
	case model.Chapter:
		h.normalizePageImageURLs(c, chapter.Pages)
		result.Data = chapter
	}
}

func (h *ChapterHandler) normalizePageImageURLs(c *gin.Context, pages []*model.Page) {
	if len(pages) == 0 {
		return
	}

	for _, page := range pages {
		if page == nil {
			continue
		}
		page.ImageURL = h.normalizeImageURL(c, page.ImageURL)
	}
}

func (h *ChapterHandler) normalizeImageURL(c *gin.Context, rawURL string) string {
	trimmed := strings.TrimSpace(rawURL)
	if trimmed == "" {
		return trimmed
	}

	if strings.HasPrefix(trimmed, "/files/content/") {
		return buildAbsoluteURL(c, trimmed)
	}
	if strings.HasPrefix(trimmed, "files/content/") {
		return buildAbsoluteURL(c, "/"+trimmed)
	}

	if !strings.HasPrefix(trimmed, "http://") && !strings.HasPrefix(trimmed, "https://") {
		return trimmed
	}

	u, err := url.Parse(trimmed)
	if err != nil {
		return trimmed
	}

	// Convert MinIO internal/presigned URLs into API file-content URLs for FE access.
	if u.Hostname() != "minio" {
		return trimmed
	}

	key := extractObjectKeyFromMinioPath(strings.TrimPrefix(u.EscapedPath(), "/"))
	if key == "" {
		return trimmed
	}

	return buildAbsoluteURL(c, "/files/content/"+key)
}

func extractObjectKeyFromMinioPath(path string) string {
	if path == "" {
		return ""
	}

	parts := strings.SplitN(path, "/", 2)
	if len(parts) < 2 {
		return path
	}

	return parts[1]
}

func buildAbsoluteURL(c *gin.Context, path string) string {
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}

	return scheme + "://" + c.Request.Host + path
}
