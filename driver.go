package persistencesql

import (
	"errors"
	"fmt"
)

// Driver defines a type of SQL driver accepted.
// This will be used by the golang sql library to load a specific driver
type Driver string

const (
	// POSTGRES driver type
	POSTGRES Driver = "postgres"
	// MYSQL driver type
	MYSQL Driver = "mysql"
)

// IsValid checks whether the given driver is valid or not
func (d Driver) IsValid() error {
	switch d {
	case POSTGRES, MYSQL:
		return nil
	}
	return errors.New("invalid driver type")
}

// String returns the actual value
func (d Driver) String() string {
	if err := d.IsValid(); err != nil {
		return ""
	}
	return string(d)
}

// ConnStr returns the connection string provided by the driver
func (d Driver) ConnStr(dbHost string, dbPort int, dbName, dbUser, dbPassword, dbSchema string) string {
	var connectionInfo string
	switch d {
	case POSTGRES:
		connectionInfo = fmt.Sprintf(
			"host=%s port=%d user=%s dbname=%s sslmode=disable search_path=%s", dbHost, dbPort, dbUser, dbName,
			dbSchema,
		)
		// The POSTGRES driver gets confused in cases where the user has no password
		// set but a password is passed, so only set password if its non-empty
		if dbPassword != "" {
			connectionInfo += fmt.Sprintf(" password=%s", dbPassword)
		}
	case MYSQL:
		connectionInfo = fmt.Sprintf(
			"%s:%s@tcp(%s:%v)/%s", dbUser, dbPassword, dbHost, dbPort, dbName,
		)
	}

	return connectionInfo
}

// SQLFile returns the sql file to create schema for a given driver
func (d Driver) SQLFile() string {
	switch d {
	case POSTGRES:
		return postgresSQL
	case MYSQL:
		return mysqlSQL
	}

	return ""
}
