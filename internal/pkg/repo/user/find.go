package userrepo

import (
	"context"
	"encoding/json"
	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/model"

	"gorm.io/gorm/clause"
)

const userCacheKeyPrefix = "user:"

func (r *UserRepository) FindByEmail(ctx context.Context, email string, moreKeys map[string]common.MoreKeyOption) (*model.User, error) {
	userString, err := r.redis.Client().Get(ctx, userCacheKeyPrefix+email).Result()
	if userString != "" && err == nil {
		jsonUser := &model.User{}
		err = json.Unmarshal([]byte(userString), jsonUser)
		if err != nil {
			return nil, err
		}

		return jsonUser, nil
	}
	user, err := r.FindOne(ctx, []any{
		clause.Eq{
			Column: "email",
			Value:  email,
		},
	}, moreKeys)
	if err != nil {
		return nil, err
	}

	userBytes, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	r.redis.Client().Set(ctx, userCacheKeyPrefix+email, string(userBytes), 0).Err()

	return user, nil
}
