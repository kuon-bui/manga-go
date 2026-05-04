package model

import (
	"encoding/json"
	"manga-go/internal/pkg/common"
	"manga-go/internal/pkg/hash"
	"time"

	"github.com/google/uuid"
	"github.com/jaswdr/faker/v2"
)

type User struct {
	common.SqlModel
	Name                  string     `json:"name" gorm:"column:name"`
	Avatar                *string    `json:"avatar" gorm:"column:avatar"`
	Email                 string     `json:"email" gorm:"column:email"`
	Password              string     `json:"-" gorm:"column:password"`
	ResetPasswordToken    string     `json:"-" gorm:"column:reset_password_token"`
	ResetPasswordExpiryAt *time.Time `json:"-" gorm:"column:reset_password_expiry_at"`
	TranslationGroupID    *uuid.UUID `json:"translationGroupId,omitempty" gorm:"column:translation_group_id"`
	UserConfig            UserConfig `json:"-" gorm:"column:user_config;type:bytea"`

	TranslationGroup *TranslationGroup `json:"translationGroup,omitempty" gorm:"foreignKey:TranslationGroupID"`
	Roles            []*Role           `json:"roles,omitempty" gorm:"many2many:users_roles;"`
}

func (User) TableName() string {
	return "users"
}

func (u User) MarshalJSON() ([]byte, error) {
	type alias User
	temp := alias(u)

	if temp.Avatar != nil {
		avatar := common.AddFileContentPrefix(*temp.Avatar)
		temp.Avatar = &avatar
	}

	return json.Marshal(temp)
}

func (u *User) GetUserId() uuid.UUID {
	return u.ID
}

func (u *User) GetEmail() string {
	return u.Email
}

func (u *User) Fake(f faker.Faker) {
	u.Name = f.Person().Name()
	u.Email = f.Internet().Email()
	u.Password = hash.HashPassword("12345678")
	u.UserConfig = DefaultUserConfig()
}
