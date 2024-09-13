package types

import (
	"database/sql"
	"time"
)

type NullTime struct {
	sql.NullTime
}

func (v NullTime) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return JSON.Marshal(v.Time)
	} else {
		return JSON.Marshal(nil)
	}
}

func (v NullTime) UnmarshalJSON(data []byte) error {
	var s *time.Time
	if err := JSON.Unmarshal(data, &s); err != nil {
		return err
	}
	if s != nil {
		v.Valid = true
		v.Time = *s
	} else {
		v.Valid = false
	}
	return nil
}
