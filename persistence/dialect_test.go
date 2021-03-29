package persistence

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDialect(t *testing.T) {
	testCases := map[string]struct {
		config DBConfig
		driver Driver
		err    error
	}{
		// asserting the creation of Postgres SQLDialect
		"postgres": {
			config: DBConfig{},
			driver: POSTGRES,
			err:    nil,
		},
		// asserting the creation of MySQL SQLDialect
		"mysql": {
			config: DBConfig{},
			driver: MYSQL,
			err:    nil,
		},
		// asserting the creation of SQL Server SQLDialect
		"sqlserver": {
			config: DBConfig{},
			driver: SQLSERVER,
			err:    nil,
		},
		// asserting the creation of Oracle SQLDialect
		"oracle": {
			config: DBConfig{},
			driver: "ORACLE",
			err:    errors.New("invalid driver type"),
		},
		// asserting unknown driver type
		"unknown": {
			config: DBConfig{},
			driver: "unknown",
			err:    errors.New("invalid driver type"),
		},
	}

	for name, testCase := range testCases {
		t.Run(
			name, func(t *testing.T) {
				// get instance of assert
				assertions := assert.New(t)

				sqlDialect, err := NewDialect(testCase.config, testCase.driver)

				if testCase.err != nil {
					assertions.Equal(testCase.err, err)
				} else {
					assertions.NotNil(sqlDialect)
					_, ok := sqlDialect.(SQLDialect)
					assertions.True(ok)
				}
			},
		)
	}
}
