package common

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type JSONMap map[string]any

func (m JSONMap) Value() (driver.Value, error) {
	if m == nil {
		return []byte("{}"), nil
	}

	return json.Marshal(m)
}

func (m *JSONMap) Scan(value any) error {
	if value == nil {
		*m = JSONMap{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("JSONMap.Scan: expected []byte, got %T", value)
	}

	if len(bytes) == 0 {
		*m = JSONMap{}
		return nil
	}

	return json.Unmarshal(bytes, m)
}
