package sqlx

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

func JsonValue(i interface{}) (driver.Value, error) {
	data, err := json.Marshal(i)
	return string(data), err
}

func JsonScan(dst, src interface{}) error {
	var data []byte
	switch v := src.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return fmt.Errorf("failed to unmarshal JSONB value: %v", src)
	}
	return json.Unmarshal(data, dst)
}
