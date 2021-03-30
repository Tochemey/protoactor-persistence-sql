package persistencesql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPostgresConnection(t *testing.T) {
	ctx := context.TODO()
	testCases := map[string]struct {
		config      DBConfig
		expectError bool
	}{
		"valid connection settings": {
			config: NewDBConfig(
				"test",
				"test",
				"testdb",
				"public",
				"localhost",
				postgresContainerPort,
				maxConnectionLifetime,
			),
			expectError: false,
		},
		"database does not exist": {
			config: NewDBConfig(
				"test",
				"test",
				"test",
				"public",
				"localhost",
				postgresContainerPort,
				maxConnectionLifetime,
			),
			expectError: true,
		},
		"authentication failed": {
			config: NewDBConfig(
				"some-username",
				"some-password",
				"testdb",
				"public",
				"localhost",
				postgresContainerPort,
				maxConnectionLifetime,
			),
			expectError: true,
		},
	}

	for testName, testCase := range testCases {
		t.Run(
			testName, func(t *testing.T) {
				// get instance of assert
				assertions := assert.New(t)

				// create the dialect instance
				dialect, err := NewPostgresDialect(testCase.config)
				assertions.NoError(err)
				assertions.NotNil(dialect)

				err = dialect.Connect(ctx)
				switch testCase.expectError {
				case false:
					assertions.NoError(err)
				default:
					assertions.Error(err)
				}
			},
		)
	}
}

func TestPostgresJournalAndSnapshotTablesCreation(t *testing.T) {
	ctx := context.TODO()
	config := NewDBConfig(
		"test",
		"test",
		"testdb",
		"public",
		"localhost",
		postgresContainerPort,
		maxConnectionLifetime,
	)

	t.Run(
		"happy path", func(t *testing.T) {
			// get instance of assert
			assertions := assert.New(t)
			// create the dialect instance
			dialect, err := NewPostgresDialect(config)
			assertions.NoError(err)
			assertions.NotNil(dialect)

			// connect to the database
			err = dialect.Connect(ctx)
			assertions.NoError(err)

			// create the journal and snapshot table successfully
			err = dialect.CreateSchemasIfNotExist(ctx)
			assertions.NoError(err)

			// check whether both tables have been created
			err = checkPostgresTable("public", "journal")
			assertions.NoError(err)
			assertions.Nil(err)
			err = checkPostgresTable("public", "snapshot")
			assertions.NoError(err)
			assertions.Nil(err)

			// insert some data into the journal and snapshot tables

		},
	)
}

func checkPostgresTable(schema, tableName string) error {
	var result string
	err := postgresHandle.
		QueryRow(fmt.Sprintf("SELECT to_regclass('%s.%s');", schema, tableName)).
		Scan(&result)
	switch {
	case err == sql.ErrNoRows, err != nil:
		return err
	default:
		if strings.EqualFold(result, "null") {
			return errors.New(fmt.Sprintf("table %s.%s does not exist", schema, tableName))
		}

		return nil
	}
}
