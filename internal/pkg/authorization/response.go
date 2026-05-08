package authorization

import (
	"context"
	"errors"

	"manga-go/internal/app/api/common/response"
)

func EnforceAnyResult(ctx context.Context, authorizer *Authorizer, req Request, contexts []Context) *response.Result {
	if authorizer == nil {
		return nil
	}

	if err := authorizer.EnforceAny(ctx, req, contexts); err != nil {
		if errors.Is(err, ErrForbidden) {
			result := response.ResultForbidden()
			return &result
		}
		result := response.ResultErrInternal(err)
		return &result
	}

	return nil
}
