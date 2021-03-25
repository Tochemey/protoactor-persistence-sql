package persistence

import "testing"

func TestDriver_ConnStr(t *testing.T) {
	type args struct {
		dbHost     string
		dbPort     int
		dbName     string
		dbUser     string
		dbPassword string
		dbSchema   string
	}
	testCases := map[string]struct {
		d        Driver
		args     args
		expected string
	}{
		"postgres connection string": {
			d: POSTGRES,
			args: args{
				dbHost:     "localhost",
				dbPort:     5432,
				dbName:     "postgres",
				dbUser:     "test",
				dbPassword: "test",
				dbSchema:   "public",
			},
			expected: "host=localhost port=5432 user=test dbName=postgres sslmode=disable search_path=public password=test",
		},

		"mysql connection string": {
			d: MYSQL,
			args: args{
				dbHost:     "localhost",
				dbPort:     3306,
				dbName:     "db",
				dbUser:     "root",
				dbPassword: "test",
			},
			expected: "root:test@tcp(localhost:3306)/db",
		},

		"sqlserver connection string": {
			d: SQLSERVER,
			args: args{
				dbHost:     "localhost",
				dbPort:     1433,
				dbName:     "tests",
				dbUser:     "sa",
				dbPassword: "test",
			},
			expected: "server=localhost;user id=sa;password=test;port=1433;database=tests;",
		},

		"not yet supported driver": {
			d:        ORACLE,
			args:     args{},
			expected: "",
		},
	}
	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			if got := testCase.d.ConnStr(testCase.args.dbHost, testCase.args.dbPort, testCase.args.dbName, testCase.args.dbUser, testCase.args.dbPassword, testCase.args.dbSchema); got != testCase.expected {
				t.Errorf("ConnStr() = %v, expected %v", got, testCase.expected)
			}
		})
	}
}

func TestDriver_IsValid(t *testing.T) {
	testCases := map[string]struct {
		d         Driver
		expectErr bool
	}{
		"postgres": {
			d:         POSTGRES,
			expectErr: false,
		},

		"mysql": {
			d:         MYSQL,
			expectErr: false,
		},

		"sqlserver": {
			d:         SQLSERVER,
			expectErr: false,
		},

		"oracle": {
			d:         ORACLE,
			expectErr: false,
		},

		"unknown": {
			d:         "DB2",
			expectErr: true,
		},
	}
	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			if err := testCase.d.IsValid(); (err != nil) != testCase.expectErr {
				t.Errorf("IsValid() error = %v, expectErr %v", err, testCase.expectErr)
			}
		})
	}
}

func TestDriver_SchemaFile(t *testing.T) {
	testCases := map[string]struct {
		d        Driver
		expected string
	}{
		"postgres": {
			d:        POSTGRES,
			expected: postgresSQL,
		},

		"mysql": {
			d:        MYSQL,
			expected: mysqlSQL,
		},

		"sqlserver": {
			d:        SQLSERVER,
			expected: sqlServerSQL,
		},

		"unknown": {
			d:        "DB2",
			expected: "",
		},
	}
	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			if got := testCase.d.SchemaFile(); got != testCase.expected {
				t.Errorf("SchemaFile() = %v, expected %v", got, testCase.expected)
			}
		})
	}
}

func TestDriver_String(t *testing.T) {
	testCases := map[string]struct {
		d        Driver
		expected string
	}{
		"postgres": {
			d:        POSTGRES,
			expected: "postgres",
		},

		"mysql": {
			d:        MYSQL,
			expected: "mysql",
		},

		"sqlserver": {
			d:        SQLSERVER,
			expected: "sqlserver",
		},

		"oracle": {
			d:        ORACLE,
			expected: "oracle",
		},

		"unknown": {
			d:        "DB2",
			expected: "",
		},
	}
	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			if got := testCase.d.String(); got != testCase.expected {
				t.Errorf("String() = %v, expected %v", got, testCase.expected)
			}
		})
	}
}
