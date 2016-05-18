package zero

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"github.com/pobri19/null-extended/convert"
)

type NullInt16 struct {
	Int16 int16
	Valid bool
}

// Int16 is a nullable int16.
// JSON marshals to zero if null.
// Considered null to SQL if zero.
type Int16 struct {
	NullInt16
}

// NewInt16 creates a new Int16
func NewInt16(i int16, valid bool) Int16 {
	return Int16{
		NullInt16: NullInt16{
			Int16: i,
			Valid: valid,
		},
	}
}

// Int16From creates a new Int16 that will be null if zero.
func Int16From(i int16) Int16 {
	return NewInt16(i, i != 0)
}

// Int16FromPtr creates a new Int16 that be null if i is nil.
func Int16FromPtr(i *int16) Int16 {
	if i == nil {
		return NewInt16(0, false)
	}
	n := NewInt16(*i, true)
	return n
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports number and null input.
// 0 will be considered a null Int16.
// It also supports unmarshalling a sql.NullInt16.
func (i *Int16) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch v.(type) {
	case float64:
		// Unmarshal again, directly to int16, to avoid intermediate float64
		err = json.Unmarshal(data, &i.Int16)
	case map[string]interface{}:
		err = json.Unmarshal(data, &i.NullInt16)
	case nil:
		i.Valid = false
		return nil
	default:
		err = fmt.Errorf("json: cannot unmarshal %v into Go value of type zero.Int16", reflect.TypeOf(v).Name())
	}
	i.Valid = (err == nil) && (i.Int16 != 0)
	return err
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null Int16 if the input is a blank, zero, or not an integer.
// It will return an error if the input is not an integer, blank, or "null".
func (i *Int16) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		i.Valid = false
		return nil
	}
	var err error
	res, err := strconv.ParseInt(string(text), 10, 16)
	i.Int16 = int16(res)
	i.Valid = (err == nil) && (i.Int16 != 0)
	return err
}

// MarshalJSON implements json.Marshaler.
// It will encode 0 if this Int16 is null.
func (i Int16) MarshalJSON() ([]byte, error) {
	n := i.Int16
	if !i.Valid {
		n = 0
	}
	return []byte(strconv.FormatInt(int64(n), 10)), nil
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a zero if this Int16 is null.
func (i Int16) MarshalText() ([]byte, error) {
	n := i.Int16
	if !i.Valid {
		n = 0
	}
	return []byte(strconv.FormatInt(int64(n), 10)), nil
}

// SetValid changes this Int16's value and also sets it to be non-null.
func (i *Int16) SetValid(n int16) {
	i.Int16 = n
	i.Valid = true
}

// Ptr returns a pointer to this Int16's value, or a nil pointer if this Int16 is null.
func (i Int16) Ptr() *int16 {
	if !i.Valid {
		return nil
	}
	return &i.Int16
}

// IsZero returns true for null or zero Int16s, for future omitempty support (Go 1.4?)
func (i Int16) IsZero() bool {
	return !i.Valid || i.Int16 == 0
}

// Scan implements the Scanner interface.
func (n *NullInt16) Scan(value interface{}) error {
	if value == nil {
		n.Int16, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return convert.ConvertAssign(&n.Int16, value)
}

// Value implements the driver Valuer interface.
func (n NullInt16) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Int16, nil
}
