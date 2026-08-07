package authorization

import (
	"context"
	"errors"

	casbinpkg "manga-go/internal/pkg/casbin"
)

var ErrForbidden = errors.New("forbidden")

// ErrAuthorizerUnavailable means the authorization engine could not answer the
// request at all. It is deliberately distinct from ErrForbidden so callers
// answer 500 instead of 403: a broken enforcer is a server fault, not a
// decision about the caller. Either way the request is denied.
var ErrAuthorizerUnavailable = errors.New("authorizer unavailable")

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
		return ErrAuthorizerUnavailable
	}

	if req.Org == "" {
		req.Org = OrgPlatform
	}
	if req.Context == "" {
		req.Context = CtxAny
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
		return ErrAuthorizerUnavailable
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
