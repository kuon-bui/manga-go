package authorization

import (
	"context"
	"errors"
	"testing"
)

func misconfiguredRequest() Request {
	return Request{
		Subject: "11111111-1111-1111-1111-111111111111",
		Org:     OrgPlatform,
		Action:  ActionManage,
		Object:  ObjectRole,
		Context: CtxAny,
	}
}

func TestAuthorizerDeniesWhenAuthorizerIsNil(t *testing.T) {
	var a *Authorizer

	err := a.Enforce(context.Background(), misconfiguredRequest())
	if err == nil {
		t.Fatal("expected nil authorizer to deny the request, got allow")
	}
	if errors.Is(err, ErrForbidden) {
		t.Fatalf("expected a misconfiguration error so callers answer 500, got ErrForbidden")
	}
}

func TestAuthorizerDeniesWhenEnforcerIsMissing(t *testing.T) {
	a := &Authorizer{}

	err := a.Enforce(context.Background(), misconfiguredRequest())
	if err == nil {
		t.Fatal("expected authorizer without enforcer to deny the request, got allow")
	}
	if errors.Is(err, ErrForbidden) {
		t.Fatalf("expected a misconfiguration error so callers answer 500, got ErrForbidden")
	}
}

func TestAuthorizerEnforceAnyDeniesWhenEnforcerIsMissing(t *testing.T) {
	a := &Authorizer{}

	err := a.EnforceAny(context.Background(), misconfiguredRequest(), DefaultContexts())
	if err == nil {
		t.Fatal("expected authorizer without enforcer to deny the request, got allow")
	}
	if errors.Is(err, ErrForbidden) {
		t.Fatalf("expected a misconfiguration error so callers answer 500, got ErrForbidden")
	}
}

func TestAuthorizerStillDeniesMissingEnforcerForImplicitlyReadableRequests(t *testing.T) {
	a := &Authorizer{}

	// A published comic read is allowed implicitly when the enforcer works; a
	// broken enforcer must not turn that shortcut into a blanket allow.
	err := a.Enforce(context.Background(), Request{
		Subject: "11111111-1111-1111-1111-111111111111",
		Org:     OrgPlatform,
		Action:  ActionRead,
		Object:  ObjectComic,
		Context: CtxPublished,
	})
	if err == nil {
		t.Fatal("expected authorizer without enforcer to deny the request, got allow")
	}
}
