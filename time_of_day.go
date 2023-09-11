package model

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"
)

const TimeOfDayFormat = "15:04"
const TimeOfDayFormatExtended = "15:04:05"

type TimeOfDay struct {
	sql.NullTime
}

func (t TimeOfDay) Value() (driver.Value, error) {
	if !t.Valid {
		return nil, nil
	}
	return t.Time.Format(TimeOfDayFormat), nil
}

func (t *TimeOfDay) Scan(value any) error {
	tym, err := time.Parse(TimeOfDayFormatExtended, string(value.([]byte)))
	if err != nil {
		return err
	}

	t.Time = tym
	t.Valid = true

	return nil
}

func (TimeOfDay) GormDataType() string {
	return "TIME"
}

func (t TimeOfDay) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return []byte("null"), nil
	}

	timeString := t.Time.Format(TimeOfDayFormat)

	return json.Marshal(timeString)
}

func (t *TimeOfDay) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte("null")) {
		t.Valid = false
		return nil
	}

	s := strings.TrimSpace(strings.Trim(string(data), `"`))

	if s == "" {
		return nil
	}

	dt, err := time.Parse(TimeOfDayFormat, s)
	if err != nil {
		return err
	}

	t.Time = dt
	t.Valid = true

	return nil
}

// UnmarshalGQL implements the graphql.Unmarshaler interface
func (t *TimeOfDay) UnmarshalGQL(v interface{}) error {
	tmfs, ok := v.(string)
	if !ok {
		return fmt.Errorf("TimeOfDay must be a string but got a %T", v)
	}

	tmfs = strings.TrimSpace(tmfs)

	if tmfs == "" {
		return nil
	}

	tmf, err := time.Parse(TimeOfDayFormat, tmfs)
	if err != nil {
		return err
	}

	t.Time = tmf
	t.Valid = true

	return nil
}

// MarshalGQL implements the graphql.Marshaler interface
func (t TimeOfDay) MarshalGQL(w io.Writer) {
	if !t.Valid {
		w.Write([]byte("null"))
	} else {
		w.Write([]byte(fmt.Sprintf("%q", t.Time.Format(TimeOfDayFormat))))
	}
}
