package persistencesql

import (
	"context"
	"database/sql"
	"log"
	"time"

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
	latestSnapshotQueryStmt    = "latest-snapshot"
	readJournalQueryStmt       = "read-journals"
	logicalJournalDeletionStmt = "logical-delete-journals"
	journalDeletionStmt        = "delete-journals"
	snapshotDeletionStmt       = "delete-snapshots"
)

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
	config *DBConfig
	db     *sql.DB

	driver Driver
	dotSQL *dotsql.DotSql
}

// NewDialect creates a new instance of SQLDialect
func NewDialect(config *DBConfig, driver Driver) (SQLDialect, error) {
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
	// create the journal table
	if _, err := d.dotSQL.ExecContext(ctx, d.db, createJournalTableStmt); err != nil {
		result = multierror.Append(result, err)
	}

	// create the snapshot table
	if _, err := d.dotSQL.ExecContext(ctx, d.db, createSnapshotTableStmt); err != nil {
		result = multierror.Append(result, err)
	}

	return result
}

// Connect connects to the database
func (d *dialect) Connect(ctx context.Context) error {
	// get the connection string provided by the driver
	connStr := d.driver.ConnStr(
		d.config.dbHost, d.config.dbPort, d.config.dbName, d.config.dbUser, d.config.dbPassword, d.config.dbSchema,
	)

	// Open the database connection
	db, err := sql.Open(d.driver.String(), connStr)
	if err != nil {
		log.Fatalf("error opening database connection: %v", err)
	}

	// set some critical database settings
	db.SetConnMaxLifetime(time.Duration(d.config.dbConnectionMaxLife) * time.Second)
	db.SetMaxIdleConns(d.config.dbMaxIdleConnections)
	db.SetMaxOpenConns(d.config.dbMaxOpenConnections)

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
	_, err := d.dotSQL.ExecContext(
		ctx,
		d.db, createSnapshotQueryStmt, snapshot.PersistenceID, snapshot.SequenceNumber, snapshot.Timestamp,
		snapshot.Snapshot,
		snapshot.SnapshotManifest, snapshot.WriterID,
	)
	return err
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
	_, err := d.dotSQL.ExecContext(ctx, d.db, snapshotDeletionStmt, persistenceID, toSequenceNumber)
	return err
}

// DeleteJournals removes some events from the journal. All events which sequence numbers are less than
// the given sequence number will be either soft deleted or hard-deleted
func (d *dialect) DeleteJournals(ctx context.Context, persistenceID string, toSequenceNumber int, logical bool) error {
	stmt := journalDeletionStmt
	if logical {
		stmt = logicalJournalDeletionStmt
	}

	// execute the query against the database
	_, err := d.dotSQL.ExecContext(ctx, d.db, stmt, persistenceID, toSequenceNumber)
	return err
}
