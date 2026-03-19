package model

import (
	"manga-go/internal/pkg/common"
	"time"

	"github.com/google/uuid"
)

type User struct {
	common.SqlModel
	Name                  string     `json:"name" gorm:"column:name"`
	Email                 string     `json:"email" gorm:"column:email"`
	Password              string     `json:"-" gorm:"column:password"`
	ResetPasswordToken    string     `json:"-" gorm:"column:reset_password_token"`
	ResetPasswordExpiryAt *time.Time `json:"-" gorm:"column:reset_password_expiry_at"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) GetUserId() uuid.UUID {
	return u.ID
}

func (u *User) GetEmail() string {
	return u.Email
}
