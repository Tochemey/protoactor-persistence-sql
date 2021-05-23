package persistencesql

import (
	"testing"
)

func TestConnectionString(t *testing.T) {
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
			expected: "host=localhost port=5432 user=test dbname=postgres sslmode=disable search_path=public password=test",
		},
		// asserting mysql connection string
		"mysql connection string": {
			driver: MYSQL,
			config: dbConfig{
				dbHost:     "localhost",
				dbPort:     3306,
				dbName:     "pg",
				dbUser:     "root",
				dbPassword: "test",
			},
			expected: "root:test@tcp(localhost:3306)/pg",
		},
		// asserting that unknown driver type will return empty string
		"not yet supported driver": {
			driver:   "ORACLE",
			config:   dbConfig{},
			expected: "",
		},
	}

	// run the test cases
	for name, testCase := range testCases {
		t.Run(
			name, func(t *testing.T) {
				if got := testCase.driver.ConnStr(
					testCase.config.dbHost, testCase.config.dbPort, testCase.config.dbName, testCase.config.dbUser,
					testCase.config.dbPassword, testCase.config.dbSchema,
				); got != testCase.expected {
					t.Errorf("ConnStr() = %v, expected %v", got, testCase.expected)
				}
			},
		)
	}
}

func TestIsValid(t *testing.T) {
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
		// asserting that oracle is driver type
		"oracle": {
			driver:    "ORACLE",
			expectErr: true,
		},
		"unknown": {
			driver:    "DB2",
			expectErr: true,
		},
	}
	for name, testCase := range testCases {
		t.Run(
			name, func(t *testing.T) {
				if err := testCase.driver.IsValid(); (err != nil) != testCase.expectErr {
					t.Errorf("IsValid() error = %v, expectErr %v", err, testCase.expectErr)
				}
			},
		)
	}
}

func TestSchemaFile(t *testing.T) {
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
		// asserting that an unknown driver type will return an empty string as schema file
		"unknown": {
			driver:   "DB2",
			expected: "",
		},
	}
	for name, testCase := range testCases {
		t.Run(
			name, func(t *testing.T) {
				if got := testCase.driver.SQLFile(); got != testCase.expected {
					t.Errorf("SQLFile() = %v, expected %v", got, testCase.expected)
				}
			},
		)
	}

}

func TestDriverString(t *testing.T) {
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
		"oracle": {
			driver:   "ORACLE",
			expected: "",
		},
		"unknown": {
			driver:   "DB2",
			expected: "",
		},
	}
	for name, testCase := range testCases {
		t.Run(
			name, func(t *testing.T) {
				if got := testCase.driver.String(); got != testCase.expected {
					t.Errorf("String() = %v, expected %v", got, testCase.expected)
				}
			},
		)
	}
}
