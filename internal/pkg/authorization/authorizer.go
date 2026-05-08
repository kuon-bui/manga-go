package authorization

import (
	"context"
	"errors"

	casbinpkg "manga-go/internal/pkg/casbin"
)

var ErrForbidden = errors.New("forbidden")

type Request struct {
	Subject string
	Org     Org
	Action  Action
	Object  Object
	Context Context
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

	if isImplicitReaderAllowed(req) {
		return nil
	}

	ok, err := a.enforcer.Enforce(
		req.Subject,
		string(req.Org),
		string(req.Action),
		string(req.Object),
		string(req.Context),
	)
	if err != nil {
		return err
	}
	if !ok {
		return ErrForbidden
	}
	return nil
}

func (a *Authorizer) EnforceAny(ctx context.Context, req Request, contexts []Context) error {
	if a == nil || a.enforcer == nil {
		return nil
	}

	if len(contexts) == 0 {
		contexts = []Context{req.Context}
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

func isImplicitReaderAllowed(req Request) bool {
	if req.Subject == "" {
		return false
	}

	switch req.Object {
	case ObjectComic, ObjectChapter:
		return req.Action == ActionRead && req.Context == CtxPublished
	case ObjectComment, ObjectRating:
		if req.Org != OrgPlatform {
			return false
		}
		return req.Action == ActionCreate && req.Context == CtxAny ||
			(req.Action == ActionUpdate || req.Action == ActionDelete) && req.Context == CtxOwner
	case ObjectReadingHistory:
		if req.Org != OrgPlatform {
			return false
		}
		return req.Action == ActionCreate && req.Context == CtxAny ||
			(req.Action == ActionRead || req.Action == ActionUpdate || req.Action == ActionDelete) && req.Context == CtxOwner
	case ObjectUser:
		if req.Org != OrgPlatform {
			return false
		}
		return req.Action == ActionUpdate && req.Context == CtxSelf
	default:
		return false
	}
}
