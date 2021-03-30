package persistencesql

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"testing"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"

	"github.com/stretchr/testify/assert"
)

var mysqlHandle *sql.DB
var mysqlContainerResource *dockertest.Resource
var mysqlContainerPool *dockertest.Pool
var mysqlContainerPort int

var postgresHandle *sql.DB
var postgresContainerPort int
var postgresContainerResource *dockertest.Resource
var postgresContainerPool *dockertest.Pool

var database = "testdb"
var maxConnectionLifetime = 120

func TestMain(m *testing.M) {
	// start the postgres test container
	postgresHandle = startPostgres()

	// start the mysql test container
	mysqlHandle = startMySQL()

	// run the tests
	code := m.Run()

	// free resources
	freeResources()

	os.Exit(code)
}

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

func startPostgres() *sql.DB {
	var db *sql.DB
	var err error
	postgresContainerPool, err = dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// create the postgres docker resource that will be used in the suite
	postgresContainerResource, err = postgresContainerPool.RunWithOptions(
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
	postgresContainerPort, _ = strconv.Atoi(postgresContainerResource.GetPort("5432/tcp"))

	if err = postgresContainerPool.Retry(
		func() error {
			var err error
			db, err = sql.Open(
				"postgres", fmt.Sprintf(
					"host=localhost port=%d user=test "+
						"password=test dbname=%s sslmode=disable", postgresContainerPort, database,
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

	return db
}

func startMySQL() *sql.DB {
	var db *sql.DB
	var err error

	mysqlContainerPool, err = dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// create the postgres docker resource that will be used in the suite
	mysqlContainerResource, err = mysqlContainerPool.RunWithOptions(
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
	mysqlContainerPort, _ = strconv.Atoi(mysqlContainerResource.GetPort("3306/tcp"))

	if err = mysqlContainerPool.Retry(
		func() error {
			var err error
			db, err = sql.Open(
				"mysql", fmt.Sprintf("test:test@(localhost:%d)/%s?parseTime=true", mysqlContainerPort, database),
			)
			if err != nil {
				return err
			}
			return db.Ping()
		},
	); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	return db
}

func freeResources() {
	if err := postgresContainerPool.Purge(postgresContainerResource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	if err := mysqlContainerPool.Purge(mysqlContainerResource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
}
