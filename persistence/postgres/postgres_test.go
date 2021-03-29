package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/assert"
	"github.com/tochemey/protoactor-persistence-sql/persistence"
)

// solely for tests
var db *sql.DB
var database = "testdb"
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
			Repository: "postgres",
			Tag:        "latest",
			Env: []string{
				"POSTGRES_USER=test",
				"POSTGRES_PASSWORD=test",
				"POSTGRES_DB=testdb",
				"listen_addresses = '*'",
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
	containerPort, _ = strconv.Atoi(resource.GetPort("5432/tcp"))

	if err = pool.Retry(
		func() error {
			var err error
			db, err = sql.Open(
				"postgres", fmt.Sprintf(
					"host=localhost port=%d user=test "+
						"password=test dbname=%s sslmode=disable", containerPort, database,
				),
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
			err = tableExists("public", "journal")
			assertions.NoError(err)
			assertions.Nil(err)
			err = tableExists("public", "snapshot")
			assertions.NoError(err)
			assertions.Nil(err)
		},
	)
}

func tableExists(schema, tableName string) error {
	var result string
	err := db.
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