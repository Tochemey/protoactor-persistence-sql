package persistence

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/denisenkom/go-mssqldb" // load the mssql driver
	"github.com/gchaincl/dotsql"
	_ "github.com/go-sql-driver/mysql" // load the mysql driver
	"github.com/hashicorp/go-multierror"
	_ "github.com/lib/pq" // loads the Postgres driver
)

const (
	createJournalTableStmt     = "create-journal-table"
	createSnapshotTableStmt    = "create-snapshot-table"
	createJournalQueryStmt     = "create-journal"
	createSnapshotQueryStmt    = "create-snapshot"
	latestSnapshotQueryStmt    = "select-snapshot"
	readJournalQueryStmt       = "read-journals"
	logicalJournalDeletionStmt = "logical-delete-journals"
	journalDeletionStmt        = "delete-journals"
	snapshotDeletionStmt       = "delete-snapshots"
)

// Config represents the database configuration
type Config struct {
	DBHost     string // the datastore host name. it can be an ip address as well
	DBPort     int    // the datastore port number
	DBUser     string // the database user name
	DBPassword string // the database password
	DBSchema   string // the schema when required, particular for postgres
	DBName     string // the database name
}

// SQLDialect will be implemented any database dialect
type SQLDialect interface {
	CreateSchemasIfNotExist(ctx context.Context) error

	Connect(ctx context.Context) error
	Close() error

	PersistJournal(ctx context.Context, journal *Journal) error
	PersistSnapshot(ctx context.Context, snapshot *Snapshot) error

	GetLatestSnapshot(ctx context.Context, persistenceID string) (*Snapshot, error)
	GetJournals(ctx context.Context, persistenceID string, fromSequenceNumber int, toSequenceNumber int) (
		[]*Journal, error,
	)

	DeleteSnapshots(ctx context.Context, persistenceID string, toSequenceNumber int) error
	DeleteJournals(ctx context.Context, persistenceID string, toSequenceNumber int, logical bool) error
}

type dialect struct {
	config Config
	db     *sql.DB

	driver Driver
	dotSQL *dotsql.DotSql
}

// NewDialect creates a new instance of SQLDialect
func NewDialect(config Config, driver Driver) (SQLDialect, error) {
	// validates driver
	if err := driver.IsValid(); err != nil {
		return nil, err
	}

	return &dialect{
		config: config,
		driver: driver,
	}, nil
}

// CreateSchemasIfNotExist creates the database tables required
func (d *dialect) CreateSchemasIfNotExist(ctx context.Context) error {
	var result error
	// Create the various tables
	if _, err := d.dotSQL.ExecContext(ctx, d.db, createJournalTableStmt); err != nil {
		result = multierror.Append(result, err)
	}

	if _, err := d.dotSQL.ExecContext(ctx, d.db, createSnapshotTableStmt); err != nil {
		result = multierror.Append(result, err)
	}

	return result
}

// Connect connects to the database
func (d *dialect) Connect(ctx context.Context) error {
	// get the connection string provided by the driver
	connStr := d.driver.ConnStr(
		d.config.DBHost, d.config.DBPort, d.config.DBName, d.config.DBUser, d.config.DBPassword, d.config.DBSchema,
	)

	// Open the database connection
	db, err := sql.Open(d.driver.String(), connStr)
	if err != nil {
		log.Fatalf("error opening database connection: %v", err)
	}

	err = db.PingContext(ctx)
	if err != nil {
		return err
	}

	// Loads sql statements
	dot, err := dotsql.LoadFromString(d.driver.SQLFile())
	if err != nil {
		return err
	}

	d.db = db
	d.dotSQL = dot
	return nil
}

// Close closes the underlying database connection
func (d *dialect) Close() error {
	return d.db.Close()
}

// PersistJournal persists a journal entry into the datastore
func (d *dialect) PersistJournal(ctx context.Context, journal *Journal) error {
	_, err := d.dotSQL.ExecContext(
		ctx,
		d.db, createJournalQueryStmt, journal.PersistenceID, journal.SequenceNumber, journal.Timestamp, journal.Payload,
		journal.EventManifest, journal.WriterID,
	)

	return err
}

// PersistSnapshot persists a snapshot entry into the snapshot data store
func (d *dialect) PersistSnapshot(ctx context.Context, snapshot *Snapshot) error {
	if _, err := d.dotSQL.ExecContext(
		ctx,
		d.db, createSnapshotQueryStmt, snapshot.PersistenceID, snapshot.SequenceNumber, snapshot.Timestamp,
		snapshot.Snapshot,
		snapshot.SnapshotManifest, snapshot.WriterID,
	); err != nil {
		return err
	}

	return nil
}

// GetLatestSnapshot fetch the latest snapshot for a given persistenceID
func (d *dialect) GetLatestSnapshot(ctx context.Context, persistenceID string) (*Snapshot, error) {
	// execute the query against the database
	row, err := d.dotSQL.QueryRowContext(ctx, d.db, latestSnapshotQueryStmt, persistenceID)
	if err != nil {
		return nil, err
	}

	// let us read the data that has been returned by the query
	var snapshot Snapshot
	err = row.Scan(
		&snapshot.PersistenceID, &snapshot.SequenceNumber, &snapshot.Timestamp,
		&snapshot.Snapshot, &snapshot.SnapshotManifest, &snapshot.WriterID,
	)

	if err != nil {
		return nil, err
	}

	return &snapshot, nil
}

// GetJournals fetch some events from the journal store
func (d *dialect) GetJournals(
	ctx context.Context, persistenceID string, fromSequenceNumber int, toSequenceNumber int,
) ([]*Journal, error) {
	events := make([]*Journal, 0)
	// execute the query against the database
	rows, err := d.dotSQL.QueryContext(
		ctx, d.db, readJournalQueryStmt, persistenceID, fromSequenceNumber, toSequenceNumber,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var journal Journal
	for rows.Next() {
		// read the row data
		if err = rows.Scan(
			&journal.Ordering, &journal.PersistenceID, &journal.SequenceNumber, &journal.Timestamp,
			&journal.Payload, &journal.EventManifest, &journal.WriterID, &journal.Deleted,
		); err != nil {
			return nil, err
		}

		// append the read row into the event slice
		events = append(events, &journal)
	}
	// get any error encountered during iteration
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

// DeleteSnapshots removes some events from the journal. All snapshots which sequence numbers are less than
// the given sequence number will be either soft deleted or hard-deleted
func (d *dialect) DeleteSnapshots(ctx context.Context, persistenceID string, toSequenceNumber int) error {
	// execute the query against the database
	if _, err := d.dotSQL.ExecContext(ctx, d.db, snapshotDeletionStmt, persistenceID, toSequenceNumber); err != nil {
		return err
	}

	return nil
}

// DeleteJournals removes some events from the journal. All events which sequence numbers are less than
// the given sequence number will be either soft deleted or hard-deleted
func (d *dialect) DeleteJournals(ctx context.Context, persistenceID string, toSequenceNumber int, logical bool) error {
	stmt := journalDeletionStmt
	if logical {
		stmt = logicalJournalDeletionStmt
	}

	// execute the query against the database
	if _, err := d.dotSQL.ExecContext(ctx, d.db, stmt, persistenceID, toSequenceNumber); err != nil {
		return err
	}

	return nil
}
