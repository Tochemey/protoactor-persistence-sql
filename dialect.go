package persistence

import (
	"context"
	"database/sql"
	"log"

	"github.com/gchaincl/dotsql"
	"github.com/hashicorp/go-multierror"
)

const (
	createJournalLabel       = "create-journal-table"
	createSnapshotLabel      = "create-snapshot-table"
	createJournalQueryLabel  = "create-journal"
	createSnapshotQuery      = "create-snapshot"
	latestSnapshotQueryLabel = "select-snapshot"
)

type Config struct {
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBSchema   string
	DBName     string
}

// SqlDialect will be implemented any database dialect
type SqlDialect interface {
	CreateSchemasIfNotExist() error

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
	dotSql *dotsql.DotSql
}

// NewDialect creates a new instance of SqlDialect
func NewDialect(config Config, driver Driver) SqlDialect {
	// validates driver
	if err := driver.IsValid(); err != nil {
		log.Fatalf("error: %v", err)
	}

	return &dialect{
		config: config,
		driver: driver,
	}
}

func (d *dialect) CreateSchemasIfNotExist() error {
	var result error
	// Create the various tables
	if _, err := d.dotSql.Exec(d.db, createJournalLabel); err != nil {
		result = multierror.Append(result, err)
	}

	if _, err := d.dotSql.Exec(d.db, createSnapshotLabel); err != nil {
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

	err = d.db.PingContext(ctx)
	if err != nil {
		return err
	}

	// Loads queries from file
	dot, err := dotsql.LoadFromFile(d.driver.SchemaSql())
	if err != nil {
		return err
	}

	d.db = db
	d.dotSql = dot
	return nil
}

// Close closes the underlying database connection
func (d *dialect) Close() error {
	return d.db.Close()
}

func (d *dialect) PersistJournal(ctx context.Context, journal *Journal) error {
	_, err := d.dotSql.ExecContext(
		ctx,
		d.db, createJournalLabel, journal.PersistenceID, journal.SequenceNumber, journal.Timestamp, journal.Payload,
		journal.EventManifest, journal.WriterID,
	)

	return err
}

func (d *dialect) PersistSnapshot(ctx context.Context, snapshot *Snapshot) error {
	_, err := d.dotSql.ExecContext(
		ctx,
		d.db, createSnapshotLabel, snapshot.PersistenceID, snapshot.SequenceNumber, snapshot.Timestamp,
		snapshot.Snapshot,
		snapshot.SnapshotManifest, snapshot.WriterID,
	)

	return err
}

func (d *dialect) GetLatestSnapshot(ctx context.Context, persistenceID string) (*Snapshot, error) {
	panic("")
}

func (d *dialect) GetJournals(
	ctx context.Context, persistenceID string, fromSequenceNumber int, toSequenceNumber int,
) ([]*Journal, error) {
	panic("implement me")
}

func (d *dialect) DeleteSnapshots(ctx context.Context, persistenceID string, toSequenceNumber int) error {
	panic("implement me")
}

func (d *dialect) DeleteJournals(ctx context.Context, persistenceID string, toSequenceNumber int, logical bool) error {
	panic("implement me")
}
