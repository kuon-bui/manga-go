package authorization

import (
	"context"
	"testing"

	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/model"

	"github.com/google/uuid"
)

func TestViewerContext(t *testing.T) {
	if viewer := ViewerFromContext(context.Background()); viewer.Subject != SubjectAnonymous || viewer.User != nil {
		t.Fatalf("expected anonymous default, got %#v", viewer)
	}

	user := &model.User{SqlModel: common.SqlModel{ID: uuid.New()}}
	ctx := WithViewer(context.Background(), user)
	viewer := ViewerFromContext(ctx)
	if !HasViewer(ctx) || viewer.Subject != user.ID.String() || viewer.User != user {
		t.Fatalf("expected authenticated viewer, got %#v", viewer)
	}
}
