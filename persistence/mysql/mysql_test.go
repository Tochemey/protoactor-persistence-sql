package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"testing"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/assert"
	"github.com/tochemey/protoactor-persistence-sql/persistence"
)

// solely for tests
var db *sql.DB
var containerPort int
var maxConnectionLifetime = 120

// Mainly to boot the database before all tests run
func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// create the postgres docker resource that will be used in the suite
	resource, err := pool.RunWithOptions(
		&dockertest.RunOptions{
			Repository: "mysql",
			Tag:        "latest",
			Env: []string{
				"MYSQL_ROOT_PASSWORD=test",
				"MYSQL_USER=test",
				"MYSQL_PASSWORD=test",
				"MYSQL_DATABASE=testdb",
			},
		}, func(config *docker.HostConfig) {
			// set AutoRemove to true so that stopped container goes away by itself
			config.AutoRemove = true
			config.RestartPolicy = docker.RestartPolicy{
				Name: "no",
			}
		},
	)

	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// get the container port to use
	containerPort, _ = strconv.Atoi(resource.GetPort("3306/tcp"))

	if err = pool.Retry(
		func() error {
			var err error
			db, err = sql.Open(
				"mysql", fmt.Sprintf("test:test@(localhost:%s)/testdb?parseTime=true", resource.GetPort("3306/tcp")),
			)
			if err != nil {
				return err
			}
			return db.Ping()
		},
	); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func TestConnection(t *testing.T) {
	ctx := context.TODO()
	testCases := map[string]struct {
		config      persistence.DBConfig
		expectError bool
	}{
		"valid connection settings": {
			config: persistence.NewDBConfig(
				"test",
				"test",
				"testdb",
				"public",
				"localhost",
				containerPort,
				maxConnectionLifetime,
			),
			expectError: false,
		},
		"database does not exist": {
			config: persistence.NewDBConfig(
				"test",
				"test",
				"test",
				"public",
				"localhost",
				containerPort,
				maxConnectionLifetime,
			),
			expectError: true,
		},
		"authentication failed": {
			config: persistence.NewDBConfig(
				"some-username",
				"some-password",
				"testdb",
				"public",
				"localhost",
				containerPort,
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

func TestJournalAndSnapshotTablesCreation(t *testing.T) {
	ctx := context.TODO()
	config := persistence.NewDBConfig(
		"test",
		"test",
		"testdb",
		"public",
		"localhost",
		containerPort,
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
			err = tableExists("testdb", "journal")
			assertions.NoError(err)
			assertions.Nil(err)
			err = tableExists("testdb", "snapshot")
			assertions.NoError(err)
			assertions.Nil(err)
		},
	)
}

func tableExists(schema, tableName string) error {
	var result string
	err := db.
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
