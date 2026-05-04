package seederutil

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

func TruncateTables(tx *gorm.DB, tables ...string) error {
	if len(tables) == 0 {
		return nil
	}

	quoted := make([]string, 0, len(tables))
	for _, table := range tables {
		quoted = append(quoted, fmt.Sprintf(`"%s"`, table))
	}

	query := fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", strings.Join(quoted, ", "))
	return tx.Exec(query).Error
}
