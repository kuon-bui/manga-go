package utils

import (
	"context"
	"fmt"
	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/model"

	"github.com/gin-gonic/gin"
)

func GetCurrentUserFromGinContext(c *gin.Context) (*model.User, error) {
	user := c.MustGet(string(common.CurrentUser)).(*model.User)
	if user == nil {
		return nil, fmt.Errorf("could not get current user")
	}

	return user, nil
}

func SetCurrentUserToGinContext(c *gin.Context, user *model.User) {
	c.Set(string(common.CurrentUser), user)
	ctx := context.WithValue(c.Request.Context(), common.CurrentUser, user)
	c.Request = c.Request.WithContext(ctx)
}

func GetCurrentUserFormContext(ctx context.Context) (*model.User, error) {
	user, ok := ctx.Value(common.CurrentUser).(*model.User)
	if !ok || user == nil {
		return nil, fmt.Errorf("could not get current user from context")
	}

	return user, nil
}

func SetTokenIdToGinContext(c *gin.Context, tokenId string) {
	c.Set(string(common.TokenId), tokenId)
}

func GetTokenIdFromGinContext(c *gin.Context) (string, error) {
	tokenId, ok := c.Get(string(common.TokenId))
	if !ok {
		return "", fmt.Errorf("could not get token id")
	}

	return tokenId.(string), nil
}
