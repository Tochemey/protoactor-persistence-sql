package persistence

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"
)

type DialectTestSuite struct {
	suite.Suite

	ctx context.Context
}

func (s *DialectTestSuite) SetupTest() {
	s.ctx = context.TODO()
}

func TestDialectTestSuite(t *testing.T) {
	suite.Run(t, new(DialectTestSuite))
}

func (s *DialectTestSuite) TestNewDialect() {
	testCases := map[string]struct {
		config Config
		driver Driver
		err    error
	}{
		// asserting the creation of Postgres SQLDialect
		"postgres": {
			config: Config{},
			driver: POSTGRES,
			err:    nil,
		},
		// asserting the creation of MySQL SQLDialect
		"mysql": {
			config: Config{},
			driver: MYSQL,
			err:    nil,
		},
		// asserting the creation of SQL Server SQLDialect
		"sqlserver": {
			config: Config{},
			driver: SQLSERVER,
			err:    nil,
		},
		// asserting the creation of Oracle SQLDialect
		"oracle": {
			config: Config{},
			driver: ORACLE,
			err:    nil,
		},
		// asserting unknown driver type
		"unknown": {
			config: Config{},
			driver: "unknown",
			err:    errors.New("invalid driver type"),
		},
	}

	for name, testCase := range testCases {
		s.Run(
			name, func() {
				sqlDialect, err := NewDialect(testCase.config, testCase.driver)

				if testCase.err != nil {
					s.Assert().Equal(testCase.err, err)
				} else {
					s.Assert().NotNil(sqlDialect)
					_, ok := sqlDialect.(SQLDialect)
					s.Assert().True(ok)
				}
			},
		)
	}
}
