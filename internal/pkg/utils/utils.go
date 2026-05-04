package utils

import (
	"crypto/rand"
	"fmt"
	"os"
	"strings"
)

func GetRootDir() string {
	workingDir, _ := os.Getwd()

	// loại bỏ phần "/cmd/api" nếu tồn tại
	rootDir := strings.TrimSuffix(workingDir, "/cmd/dev")
	rootDir = strings.TrimSuffix(rootDir, "/cmd/queue")

	return rootDir
}

func TokenGenerator(len int) string {
	b := make([]byte, len)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
