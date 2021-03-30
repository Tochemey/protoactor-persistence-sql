package persistencesql

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMySQLConnection(t *testing.T) {
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
				mysqlContainerPort,
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
				mysqlContainerPort,
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
				mysqlContainerPort,
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
				dialect, err := NewMySQLDialect(testCase.config)
				assertions.NoError(err)
				assertions.NotNil(dialect)

				err = dialect.Connect(ctx)

				switch testCase.expectError {
				case true:
					assertions.Error(err)
				default:
					assertions.NoError(err)
				}
			},
		)
	}
}

func TestMySQLJournalAndSnapshotTablesCreation(t *testing.T) {
	ctx := context.TODO()
	config := NewDBConfig(
		"test",
		"test",
		"testdb",
		"public",
		"localhost",
		mysqlContainerPort,
		maxConnectionLifetime,
	)

	t.Run(
		"happy path", func(t *testing.T) {
			// get instance of assert
			assertions := assert.New(t)
			// create the dialect instance
			dialect, err := NewMySQLDialect(config)
			assertions.NoError(err)
			assertions.NotNil(dialect)

			// connect to the database
			err = dialect.Connect(ctx)
			assertions.NoError(err)

			// create the journal and snapshot table successfully
			err = dialect.CreateSchemasIfNotExist(ctx)
			assertions.NoError(err)

			// check whether both tables have been created
			err = checkMySQLTable("testdb", "journal")
			assertions.NoError(err)
			assertions.Nil(err)
			err = checkMySQLTable("testdb", "snapshot")
			assertions.NoError(err)
			assertions.Nil(err)
		},
	)
}

func checkMySQLTable(schema, tableName string) error {
	var result string
	err := mysqlHandle.
		QueryRow(
			fmt.Sprintf(
				"SELECT table_name FROM information_schema.tables WHERE table_schema = '%s' AND table_name = '%s' LIMIT 1; ",
				schema, tableName,
			),
		).
		Scan(&result)
	switch {
	case err == sql.ErrNoRows, err != nil:
		return err
	default:
		return nil
	}
}
