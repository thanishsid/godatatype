package model

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"
)

type Duration struct {
	Duration time.Duration
	Valid    bool
}

func (d Duration) Value() (driver.Value, error) {
	if d.Valid {
		return d.Duration, nil
	}

	return nil, nil
}

func (d *Duration) Scan(value any) error {
	if value == nil {
		d.Duration, d.Valid = 0, false
		return nil
	}

	d.Valid = true

	duration, ok := value.(int64)

	if !ok {
		return fmt.Errorf("Duration scan error: expected int64 but got %T", duration)
	}

	d.Duration = time.Duration(duration)

	return nil
}

func (Duration) GormDataType() string {
	return "BIGINT"
}

func (d Duration) MarshalJSON() ([]byte, error) {
	if !d.Valid {
		return []byte("null"), nil
	}

	durationString := d.Duration.String()

	return json.Marshal(durationString)
}

func (d *Duration) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte("null")) {
		d.Valid = false
		return nil
	}

	s := strings.TrimSpace(strings.Trim(string(data), `"`))

	if s == "" {
		return nil
	}

	duration, err := time.ParseDuration(s)
	if err != nil {
		return err
	}

	d.Duration = duration
	d.Valid = true

	return nil
}

// UnmarshalGQL implements the graphql.Unmarshaler interface
func (d *Duration) UnmarshalGQL(v any) error {
	durationStr, ok := v.(string)
	if !ok {
		return fmt.Errorf("Duration must be a string but got a %T", v)
	}

	durationStr = strings.TrimSpace(durationStr)

	if durationStr == "" {
		return nil
	}

	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return err
	}

	d.Duration = duration
	d.Valid = true

	return nil
}

// MarshalGQL implements the graphql.Marshaler interface
func (d Duration) MarshalGQL(w io.Writer) {
	if !d.Valid {
		w.Write([]byte("null"))
	} else {
		w.Write([]byte(d.Duration.String()))
	}
}
