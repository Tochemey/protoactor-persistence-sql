package persistencesql

// DBConfig represents the database configuration
type DBConfig struct {
	dbHost               string // the datastore host name. it can be an ip address as well
	dbPort               int    // the datastore dbPort number
	dbUser               string // the database user name
	dbPassword           string // the database password
	dbSchema             string // the schema when required, particular for postgres
	dbName               string // the database name
	dbConnectionMaxLife  int    // to ensure connections are closed by the driver safely
	dbMaxOpenConnections int
	dbMaxIdleConnections int
}

// PoolOpt defines the connection pool options
type PoolOpt = func(*DBConfig)

// NewDBConfig creates an instance of DBConfig
func NewDBConfig(dbUser, dbPassword, dbName, dbSchema, dbHost string, dbPort int, opts ...PoolOpt) *DBConfig {
	dbConfig := &DBConfig{
		dbHost:               dbHost,
		dbPort:               dbPort,
		dbUser:               dbUser,
		dbPassword:           dbPassword,
		dbSchema:             dbSchema,
		dbName:               dbName,
		dbMaxIdleConnections: 10,
		dbMaxOpenConnections: 10,
		dbConnectionMaxLife:  3,
	}

	// set pool settings if defined
	for _, opt := range opts {
		opt(dbConfig)
	}

	return dbConfig
}

// WithConnectionMaxLife sets the database connection max life
func WithConnectionMaxLife(connectionMaxLife int) PoolOpt {
	return func(config *DBConfig) {
		config.dbConnectionMaxLife = connectionMaxLife
	}
}

// WithMaxOpenConnections sets the max open connections
func WithMaxOpenConnections(maxOpenConnection int) PoolOpt {
	return func(config *DBConfig) {
		config.dbMaxOpenConnections = maxOpenConnection
	}
}

// WithMaxIdleConnections sets the max idle connections
func WithMaxIdleConnections(maxIdleConnections int) PoolOpt {
	return func(config *DBConfig) {
		config.dbMaxIdleConnections = maxIdleConnections
	}
}
