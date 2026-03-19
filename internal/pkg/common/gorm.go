package common

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SqlModel struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	CreatedAt *time.Time     `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt *time.Time     `json:"updatedAt" gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"column:deleted_at"`
}

type MoreKeyOption struct {
	Unscoped bool                    // Bỏ qua soft delete
	Order    *clause.OrderBy         // Sắp xếp trong Preload
	Where    *clause.Where           // Điều kiện lọc trong Preload
	Select   []string                // Các cột cần lấy
	Limit    *int                    // Giới hạn số lượng preload
	Custom   func(*gorm.DB) *gorm.DB // Hàm preload tùy chỉnh
}

type JoinExpr struct {
	SQL  string
	Vars []any
}

func GenerateModelCode(counter int, numberLength int, prefixCode string) string {
	numberStr := fmt.Sprintf("%0*d", numberLength, counter)
	return prefixCode + numberStr
}
