package godatatype

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

type DateTime struct {
	sql.NullTime
}

const DateTimeFormat = "2006-01-02T15:04"

func (d DateTime) Value() (driver.Value, error) {
	if !d.Valid {
		return nil, nil
	}
	return d.Time, nil
}

func (d DateTime) MarshalJSON() ([]byte, error) {
	if !d.Valid {
		return []byte("null"), nil
	}

	timeString := d.Time.Format(DateTimeFormat)

	return json.Marshal(timeString)
}

func (d *DateTime) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte("null")) {
		d.Valid = false
		return nil
	}

	s := strings.TrimSpace(strings.Trim(string(data), `"`))

	if s == "" {
		return nil
	}

	dt, err := time.Parse(DateTimeFormat, s)
	if err != nil {
		return err
	}

	d.Time = dt
	d.Valid = true

	return nil
}

// UnmarshalGQL implements the graphql.Unmarshaler interface
func (d *DateTime) UnmarshalGQL(v interface{}) error {
	if v == nil {
		return nil
	}

	dts, ok := v.(string)
	if !ok {
		return fmt.Errorf("DateTime must be a string but got a %T", v)
	}

	dts = strings.TrimSpace(dts)

	if dts == "" {
		return nil
	}

	dt, err := time.Parse(DateTimeFormat, dts)
	if err != nil {
		return err
	}

	d.Time = dt
	d.Valid = true

	return nil
}

// MarshalGQL implements the graphql.Marshaler interface
func (d DateTime) MarshalGQL(w io.Writer) {
	if !d.Valid {
		w.Write([]byte("null"))
	} else {
		w.Write([]byte(fmt.Sprintf(`"%s"`, d.Time.Format(DateTimeFormat))))
	}
}

func (d DateTime) GormDataType() string {
	return "datetime(3)"
}
