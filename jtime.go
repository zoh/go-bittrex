package bittrex

import (
	"encoding/json"
	"time"
	"fmt"
)

const TIME_FORMAT = "2006-01-02T15:04:05"

type jTime struct {
	time.Time
	Valid bool
}

func (jt *jTime) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s == "" {
		return nil
	}

	t, err := time.Parse(TIME_FORMAT, s)
	if err != nil {
		return err
	}
	jt.Time = t
	jt.Valid = true
	return nil
}

func (jt jTime) MarshalJSON() ([]byte, error) {
	if !jt.Valid {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf(`"%s"`, jt.Format(TIME_FORMAT))), nil
}
