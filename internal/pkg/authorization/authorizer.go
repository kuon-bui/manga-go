package authorization

import (
	"context"
	"errors"

	casbinpkg "manga-go/internal/pkg/casbin"
)

var ErrForbidden = errors.New("forbidden")

type Request struct {
	Subject string
	Org     string
	Action  string
	Object  string
	Context string
}

type Authorizer struct {
	enforcer *casbinpkg.Enforcer
}

func NewAuthorizer(enforcer *casbinpkg.Enforcer) *Authorizer {
	return &Authorizer{enforcer: enforcer}
}

func (a *Authorizer) Enforce(ctx context.Context, req Request) error {
	if a == nil || a.enforcer == nil {
		return nil
	}

	if req.Org == "" {
		req.Org = OrgPlatform
	}
	if req.Context == "" {
		req.Context = CtxAny
	}

	ok, err := a.enforcer.Enforce(req.Subject, req.Org, req.Action, req.Object, req.Context)
	if err != nil {
		return err
	}
	if !ok {
		return ErrForbidden
	}
	return nil
}

func (a *Authorizer) EnforceAny(ctx context.Context, req Request, contexts []string) error {
	if a == nil || a.enforcer == nil {
		return nil
	}

	if len(contexts) == 0 {
		contexts = []string{req.Context}
	}

	var lastErr error
	for _, requestContext := range contexts {
		req.Context = requestContext
		if err := a.Enforce(ctx, req); err == nil {
			return nil
		} else if !errors.Is(err, ErrForbidden) {
			lastErr = err
		}
	}

	if lastErr != nil {
		return lastErr
	}
	return ErrForbidden
}
