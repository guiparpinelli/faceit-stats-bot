package sqlite

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type UnixTime struct {
	time.Time
}

func (ut *UnixTime) Scan(src interface{}) error {
	if src == nil {
		ut.Time = time.Time{}
		return nil
	}

	switch v := src.(type) {
	case int64:
		ut.Time = time.Unix(v, 0)
	case int: // Handle int in case SQLite driver returns it
		ut.Time = time.Unix(int64(v), 0)
	default:
		return fmt.Errorf("unsupported type for UnixTime: %T, expected int64", src)
	}
	return nil
}

func (ut UnixTime) Value() (driver.Value, error) {
	return ut.Unix(), nil
}
