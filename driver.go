package persistence

import (
	"errors"
	"fmt"
)

type Driver string

const (
	Postgres  Driver = "postgres"
	MySQL     Driver = "mysql"
	Oracle    Driver = "oracle"
	SqlServer Driver = "sqlserver"
)

// IsValid checks whether the given driver is valid or not
func (d Driver) IsValid() error {
	switch d {
	case Postgres, MySQL, Oracle, SqlServer:
		return nil
	}
	return errors.New("invalid driver type")
}

// String returns the actual value
func (d Driver) String() string {
	return string(d)
}

// ConnStr returns the connection string provided by the driver
func (d Driver) ConnStr(dbHost string, dbPort int, dbName, dbUser, dbPassword, dbSchema string) string {
	var connectionInfo string
	switch d {
	case Postgres:
		connectionInfo = fmt.Sprintf(
			"host=%s port=%d user=%s dbName=%s sslmode=disable search_path=%s", dbHost, dbPort, dbUser, dbName,
			dbSchema,
		)
		// The postgres driver gets confused in cases where the user has no password
		// set but a password is passed, so only set password if its non-empty
		if dbPassword != "" {
			connectionInfo += fmt.Sprintf(" password=%s", dbPassword)
		}
	case MySQL:
		connectionInfo = fmt.Sprintf(
			"%s:%s@tcp(%s:%v)/%s", dbUser, dbPassword, dbHost, dbPort, dbName,
		)
	}

	return connectionInfo
}

// SchemaSql returns the sql file to create schema for a given driver
func (d Driver) SchemaSql() string {
	return fmt.Sprintf("schemas/%s-schemas.sql", d.String())
}
