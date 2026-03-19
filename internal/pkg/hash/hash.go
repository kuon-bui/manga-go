package hash

import (
	"base-go/internal/pkg/constant"
	"base-go/internal/pkg/logger"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(pwd string) string {
	// Use GenerateFromPassword to hash & salt pwd
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), constant.HashPasswordCost)
	if err != nil {
		logger.GetLogger().Error(fmt.Errorf("hash and salt %v: %w", string(pwd), err))
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash)
}

func ComparePassword(hashedPwd string, plainPwd string) bool {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)
	bytePwd := []byte(plainPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, bytePwd)
	return err == nil
}
