package translationgroupservice

import (
	"context"
	"io"
	"manga-go/internal/app/api/common/response"
	"path/filepath"

	"github.com/google/uuid"
	"gorm.io/gorm/clause"
)

func (s *TranslationGroupService) UploadLogo(ctx context.Context, groupID uuid.UUID, file io.Reader, size int64, contentType string, filename string) response.Result {
	// First ensure group exists
	group, err := s.translationGroupRepo.FindOne(ctx, []any{
		clause.Eq{Column: "id", Value: groupID},
	}, nil)
	if err != nil {
		s.logger.Error("Failed to find translation group", "error", err)
		return response.ResultErrDb(err)
	}

	ext := filepath.Ext(filename)
	if ext == "" {
		ext = ".png" // fallback
	}

	// Generate a unique filename for S3
	objectKey := "translation-groups/" + groupID.String() + "/logo" + ext

	// Upload to object storage
	if err := s.objectStorage.UploadFile(ctx, objectKey, file, size, contentType); err != nil {
		s.logger.Error("Failed to upload logo", "error", err)
		return response.ResultErrInternal(err)
	}

	// Update DB record
	// To construct public URL, usually it requires endpoint/bucket. We assume frontend prepends generic path,
	// or we can store just the object key and let FE use presigned URL logic.
	// But according to missing-api.md, it's just `logoUrl`. 
	// For simplicity, we store the objectKey. FE might format it, or we format it if we have config.
	// To be safe, we'll store the objectKey because missing-api.md says "redirected to file upload API", so usually the path is returned.
	// Wait! AWS S3 has a public domain format.
	publicURL := objectKey // Storing object key as logo url. The FileService direct upload returns the object key. Wait, FE can prefix.

	err = s.translationGroupRepo.Update(ctx, []any{
		clause.Eq{Column: "id", Value: groupID},
	}, map[string]any{
		"logo_url": publicURL,
	})
	if err != nil {
		s.logger.Error("Failed to update logo url in db", "error", err)
		return response.ResultErrDb(err)
	}
	
	group.LogoUrl = &publicURL

	return response.ResultSuccess("Logo uploaded successfully", map[string]string{
		"logoUrl": publicURL,
	})
}
