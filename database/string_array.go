package database

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type StringArray []string

func (s StringArray) Value() (driver.Value, error) {
	if s == nil {
		return "[]", nil
	}
	data, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	return string(data), nil
}

func (s *StringArray) Scan(src any) error {
	if src == nil {
		return nil
	}
	data, ok := src.(string)
	if !ok {
		return fmt.Errorf("StringArray not encoded as a String")
	}
	return json.Unmarshal([]byte(data), s)
}
