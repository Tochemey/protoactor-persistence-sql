package persistencesql

import (
	"database/sql/driver"
	"errors"

	"google.golang.org/protobuf/reflect/protoreflect"
)

type Manifest protoreflect.FullName

// Value - Implementation of valuer for database/sql
func (m Manifest) Value() (driver.Value, error) {
	return string(m), nil
}

// Scan - Implement the database/sql scanner interface
func (m *Manifest) Scan(value interface{}) error {
	var source []byte
	switch value.(type) {
	case string:
		source = []byte(value.(string))
	case []byte:
		source = value.([]byte)
	default:
		return errors.New("incompatible type for Manifest")
	}

	str := string(source)
	*m = Manifest(str)
	return nil
}
