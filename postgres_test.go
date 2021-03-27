package persistence

import (
	"context"
	"strconv"
	"testing"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"

	"github.com/stretchr/testify/suite"
)

type PostgresTestSuite struct {
	suite.Suite

	ctx           context.Context
	dockerPool    *dockertest.Pool
	dockerRes     *dockertest.Resource
	containerPort int
}

func (s *PostgresTestSuite) SetupSuite() {
	s.ctx = context.TODO()
	var err error

	// create a new docker pool
	s.dockerPool, err = dockertest.NewPool("")
	s.Assert().NoError(err)

	// create the postgres docker resource that will be used in the suite
	s.dockerRes, err = s.dockerPool.RunWithOptions(
		&dockertest.RunOptions{
			Repository: "postgres",
			Tag:        "11",
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
	s.Assert().NoError(err)

	s.containerPort, err = strconv.Atoi(s.dockerRes.GetPort("5432/tcp"))
	s.Assert().NoError(err)
}

func (s *PostgresTestSuite) TearDownSuite() {
	err := s.dockerPool.Purge(s.dockerRes)
	s.Assert().NoError(err)
}

func TestPostgresTestSuite(t *testing.T) {
	suite.Run(t, new(PostgresTestSuite))
}

func (s *PostgresTestSuite) TestConnect() {
	testCases := map[string]struct {
		config Config
		driver Driver
		err    error
	}{
		"happy path": {
			config: Config{
				DBHost:     "localhost",
				DBPort:     s.containerPort,
				DBUser:     "test",
				DBPassword: "test",
				DBSchema:   "public",
				DBName:     "testdb",
			},
			driver: POSTGRES,
			err:    nil,
		},
	}

	for name, testCase := range testCases {
		s.Run(
			name, func() {
				// create the dialect instance
				sqlDialect, err := NewDialect(testCase.config, testCase.driver)
				s.Assert().NotNil(sqlDialect)
				s.Assert().NoError(err)

				err = s.dockerPool.Retry(
					func() error {
						// create the various database schemas for the given dialect
						return sqlDialect.Connect(s.ctx)
					},
				)

				s.Assert().Equal(testCase.err, err)
			},
		)
	}
}
