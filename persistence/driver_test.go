package persistence

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type DriverTestSuite struct {
	suite.Suite
}

func (s *DriverTestSuite) SetupTest() {}

func TestDriverTestSuite(t *testing.T) {
	suite.Run(t, new(DriverTestSuite))
}

func (s *DriverTestSuite) TestConnectionString() {
	type dbConfig struct {
		dbHost     string
		dbPort     int
		dbName     string
		dbUser     string
		dbPassword string
		dbSchema   string
	}
	testCases := map[string]struct {
		driver   Driver
		config   dbConfig
		expected string
	}{
		// asserting postgres connection string
		"postgres connection string": {
			driver: POSTGRES,
			config: dbConfig{
				dbHost:     "localhost",
				dbPort:     5432,
				dbName:     "postgres",
				dbUser:     "test",
				dbPassword: "test",
				dbSchema:   "public",
			},
			expected: "host=localhost dbPort=5432 user=test dbname=postgres sslmode=disable search_path=public password=test",
		},
		// asserting mysql connection string
		"mysql connection string": {
			driver: MYSQL,
			config: dbConfig{
				dbHost:     "localhost",
				dbPort:     3306,
				dbName:     "db",
				dbUser:     "root",
				dbPassword: "test",
			},
			expected: "root:test@tcp(localhost:3306)/db",
		},
		// asserting sql server connection string
		"sqlserver connection string": {
			driver: SQLSERVER,
			config: dbConfig{
				dbHost:     "localhost",
				dbPort:     1433,
				dbName:     "tests",
				dbUser:     "sa",
				dbPassword: "test",
			},
			expected: "server=localhost;user id=sa;password=test;dbPort=1433;database=tests;",
		},
		// asserting that unknown driver type will return empty string
		"not yet supported driver": {
			driver:   ORACLE,
			config:   dbConfig{},
			expected: "",
		},
	}

	// run the test cases
	for name, testCase := range testCases {
		s.Run(
			name, func() {
				if got := testCase.driver.ConnStr(
					testCase.config.dbHost, testCase.config.dbPort, testCase.config.dbName, testCase.config.dbUser,
					testCase.config.dbPassword, testCase.config.dbSchema,
				); got != testCase.expected {
					s.T().Errorf("ConnStr() = %v, expected %v", got, testCase.expected)
				}
			},
		)
	}
}

func (s *DriverTestSuite) TestIsValid() {
	testCases := map[string]struct {
		driver    Driver
		expectErr bool
	}{
		// asserting that postgres is driver type
		"postgres": {
			driver:    POSTGRES,
			expectErr: false,
		},
		// asserting that mysql is driver type
		"mysql": {
			driver:    MYSQL,
			expectErr: false,
		},
		// asserting that sqlserver is driver type
		"sqlserver": {
			driver:    SQLSERVER,
			expectErr: false,
		},
		// asserting that oracle is driver type
		"oracle": {
			driver:    ORACLE,
			expectErr: false,
		},
		"unknown": {
			driver:    "DB2",
			expectErr: true,
		},
	}
	for name, testCase := range testCases {
		s.Run(
			name, func() {
				if err := testCase.driver.IsValid(); (err != nil) != testCase.expectErr {
					s.T().Errorf("IsValid() error = %v, expectErr %v", err, testCase.expectErr)
				}
			},
		)
	}
}

func (s *DriverTestSuite) TestSchemaFile() {
	testCases := map[string]struct {
		driver   Driver
		expected string
	}{
		// asserting that the postgres driver will return the correct schema file
		"postgres": {
			driver:   POSTGRES,
			expected: postgresSQL,
		},
		// asserting that the mysql driver will return the correct schema file
		"mysql": {
			driver:   MYSQL,
			expected: mysqlSQL,
		},
		// asserting that the sqlserver driver will return the correct schema file
		"sqlserver": {
			driver:   SQLSERVER,
			expected: sqlServerSQL,
		},
		// asserting that an unknown driver type will return an empty string as schema file
		"unknown": {
			driver:   "DB2",
			expected: "",
		},
	}
	for name, testCase := range testCases {
		s.Run(
			name, func() {
				if got := testCase.driver.SQLFile(); got != testCase.expected {
					s.T().Errorf("SQLFile() = %v, expected %v", got, testCase.expected)
				}
			},
		)
	}

}

func (s *DriverTestSuite) TestDriverString() {
	testCases := map[string]struct {
		driver   Driver
		expected string
	}{
		"postgres": {
			driver:   POSTGRES,
			expected: "postgres",
		},
		"mysql": {
			driver:   MYSQL,
			expected: "mysql",
		},
		"sqlserver": {
			driver:   SQLSERVER,
			expected: "sqlserver",
		},
		"oracle": {
			driver:   ORACLE,
			expected: "oracle",
		},
		"unknown": {
			driver:   "DB2",
			expected: "",
		},
	}
	for name, testCase := range testCases {
		s.Run(
			name, func() {
				if got := testCase.driver.String(); got != testCase.expected {
					s.T().Errorf("String() = %v, expected %v", got, testCase.expected)
				}
			},
		)
	}
}
