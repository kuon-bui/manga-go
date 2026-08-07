package authorization

import (
	"context"

	"manga-go/internal/pkg/model"
)

type viewerContextKey struct{}

type ViewerContext struct {
	Subject string
	User    *model.User
}

func ViewerFromContext(ctx context.Context) ViewerContext {
	viewer, ok := ctx.Value(viewerContextKey{}).(ViewerContext)
	if !ok || viewer.Subject == "" {
		return ViewerContext{Subject: SubjectAnonymous}
	}
	return viewer
}

func HasViewer(ctx context.Context) bool {
	_, ok := ctx.Value(viewerContextKey{}).(ViewerContext)
	return ok
}

func WithViewer(ctx context.Context, user *model.User) context.Context {
	viewer := ViewerContext{Subject: SubjectAnonymous}
	if user != nil {
		viewer.Subject = Subject(user.ID)
		viewer.User = user
	}
	return context.WithValue(ctx, viewerContextKey{}, viewer)
}
