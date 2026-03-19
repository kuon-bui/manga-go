package utils

import (
	"fmt"
	"html/template"
	"manga-go/internal/pkg/logger"
	"path/filepath"
	"strings"
	"sync"
)

var (
	mailTemplate map[string]*template.Template
	mu           sync.RWMutex
)

func LoadMailTemplate(logger *logger.Logger) {
	mu.Lock()
	defer mu.Unlock()

	// lấy thư mục gốc làm việc hiện tại
	rootDir := GetRootDir()
	dir := filepath.Join(rootDir, "resources", "mails", "*.html")
	fmt.Printf("Loading mail templates from directory: %s\n", dir)
	mailTemplate = make(map[string]*template.Template)
	files, _ := filepath.Glob(dir)
	for _, file := range files {
		name := strings.TrimSuffix(filepath.Base(file), ".html")
		template, err := template.ParseFiles(file)
		if err != nil {
			logger.Error("Failed to load template", "file", name, "error", err)
			panic(err)
		}

		mailTemplate[name] = template
	}
}

func GetMailTemplate(name string) (*template.Template, bool) {
	mu.RLock()
	defer mu.RUnlock()
	template, ok := mailTemplate[name]
	return template, ok
}
