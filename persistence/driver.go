package persistence

import (
	"errors"
	"fmt"
)

const (
	postgresSQL = `
		-- name: create-journal-table
		CREATE TABLE IF NOT EXISTS journal
		(
		    ordering        BIGSERIAL UNIQUE,
		    persistence_id  VARCHAR(255)          NOT NULL,
		    sequence_number BIGINT                NOT NULL,
		    timestamp       BIGINT                NOT NULL,
		    payload         BYTEA                 NOT NULL,
		    manifest        VARCHAR(255)          NOT NULL,
		    writer_id       VARCHAR(255)          NOT NULL,
		    deleted         BOOLEAN DEFAULT FALSE NOT NULL,
		    PRIMARY KEY (persistence_id, sequence_number)
		);
		
		-- name: create-snapshot-table
		CREATE TABLE IF NOT EXISTS snapshot
		(
		    persistence_id  VARCHAR(255) NOT NULL,
		    sequence_number BIGINT       NOT NULL,
		    timestamp       BIGINT       NOT NULL,
		    snapshot        BYTEA        NOT NULL,
		    manifest        VARCHAR(255) NOT NULL,
		    writer_id       VARCHAR(255) NOT NULL,
		    PRIMARY KEY (persistence_id, sequence_number)
		);
		
		-- name: create-journal
		INSERT INTO journal (persistence_id, sequence_number, timestamp, payload, manifest, writer_id)
		VALUES (?, ?, ?, ?, ?, ?);
		
		-- name: create-snapshot
		INSERT INTO snapshot (persistence_id, sequence_number, timestamp, snapshot, manifest, writer_id)
		VALUES (?, ?, ?, ?, ?, ?);
		
		-- name: latest-snapshot
		SELECT *
		FROM snapshot
		WHERE persistence_id = ?
		ORDER BY sequence_number DESC
		LIMIT 1

		-- name: read-journals
		SELECT * FROM journal 
		WHERE persistence_id = ? AND sequence_number >= ? AND sequence_number <= ? AND NOT deleted
		ORDER BY sequence_number ASC

		-- name: delete-journals
		DELETE FROM journal 
		WHERE persistence_id = ? AND sequence_number <= ?

		-- name: logical-delete-journals
		UPDATE journal
		SET deleted = TRUE
		WHERE persistence_id = ? AND sequence_number <= ?

		-- name: delete-snapshot
		DELETE FROM snapshots 
		WHERE persistence_id = ? AND sequence_number <= ?
	`
	mysqlSQL = `
		-- name: create-journal-table
		CREATE TABLE IF NOT EXISTS journal
		(
		    ordering        SERIAL,
		    persistence_id  VARCHAR(255)          NOT NULL,
		    sequence_number BIGINT UNSIGNED       NOT NULL,
		    timestamp       BIGINT                NOT NULL,
		    payload         BLOB                  NOT NULL,
		    manifest        VARCHAR(255)          NOT NULL,
		    writer_id       VARCHAR(255)          NOT NULL,
		    deleted         BOOLEAN DEFAULT FALSE NOT NULL,
		    PRIMARY KEY (persistence_id, sequence_number)
		);
		
		-- name: create-snapshot-table
		CREATE TABLE IF NOT EXISTS snapshot
		(
		    persistence_id  VARCHAR(255)    NOT NULL,
		    sequence_number BIGINT UNSIGNED NOT NULL,
		    timestamp       BIGINT UNSIGNED NOT NULL,
		    snapshot        BLOB            NOT NULL,
		    manifest        VARCHAR(255)    NOT NULL,
		    writer_id       VARCHAR(255)    NOT NULL,
		    PRIMARY KEY (persistence_id, sequence_number)
		);
		
		-- name: create-journal
		INSERT INTO journal (persistence_id, sequence_number, timestamp, payload, manifest, writer_id)
		VALUES (?, ?, ?, ?, ?, ?);
		
		-- name: create-snapshot
		INSERT INTO snapshot (persistence_id, sequence_number, timestamp, snapshot, manifest, writer_id)
		VALUES (?, ?, ?, ?, ?, ?);
		
		-- name: latest-snapshot
		SELECT *
		FROM snapshot
		WHERE persistence_id = ?
		ORDER BY sequence_number DESC
		LIMIT 1

		-- name: read-journals
		SELECT * FROM journal 
		WHERE persistence_id = ? AND sequence_number >= ? AND sequence_number <= ? AND deleted IS NOT TRUE
		ORDER BY sequence_number ASC

		-- name: delete-journals
		DELETE FROM journal 
		WHERE persistence_id = ? AND sequence_number <= ?

		-- name: logical-delete-journals
		UPDATE journal
		SET deleted = TRUE
		WHERE persistence_id = ? AND sequence_number <= ?

		-- name: delete-snapshots
		DELETE FROM snapshot 
		WHERE persistence_id = ? AND sequence_number <= ?
	`
	sqlServerSQL = `
		-- name: create-journal-table
		IF NOT EXISTS(SELECT 1
		              FROM sys.objects
		              WHERE object_id = OBJECT_ID(N'"journal"')
		                AND type in (N'U'))
		    begin
		        CREATE TABLE journal
		        (
		            "ordering"        BIGINT IDENTITY (1,1) NOT NULL,
		            "persistence_id"  VARCHAR(255)          NOT NULL,
		            "sequence_number" NUMERIC(10, 0)        NOT NULL,
		            "timestamp"       BIGINT                NOT NULL,
		            "payload"         VARBINARY(MAX)        NOT NULL,
		            "manifest"        VARCHAR(MAX)          NOT NULL,
		            "writer_id"       VARCHAR(255)          NOT NULL,
		            "deleted"         BIT DEFAULT 0         NOT NULL,
		            PRIMARY KEY ("persistence_id", "sequence_number")
		        )
		        CREATE UNIQUE INDEX journal_ordering_idx ON journal (ordering)
		    end;

		-- name: create-snapshot-table
		IF NOT EXISTS(SELECT 1
		              FROM sys.objects
		              WHERE object_id = OBJECT_ID(N'"snapshot"')
		                AND type in (N'U'))
		CREATE TABLE snapshot
		(
		    "persistence_id"  VARCHAR(255)   NOT NULL,
		    "sequence_number" NUMERIC(10, 0) NOT NULL,
		    "timestamp"       NUMERIC        NOT NULL,
		    "snapshot"        VARBINARY(max) NOT NULL,
		    "manifest"        VARCHAR(MAX)   NOT NULL,
		    "writer_id"       VARCHAR(255)   NOT NULL,
		    PRIMARY KEY ("persistence_id", "sequence_number")
		);
		end

		-- name: create-journal
		INSERT INTO journal(persistence_id, sequence_number, timestamp, payload, manifest, writer_id)
		VALUES (?, ?, ?, ?, ?, ?);

		-- name: create-snapshot
		INSERT INTO snapshot (persistence_id, sequence_number, timestamp, snapshot, manifest, writer_id)
		VALUES (?, ?, ?, ?, ?, ?);

		-- name: latest-snapshot
		SELECT TOP 1 *
		FROM snapshot
		WHERE persistence_id = ?
		ORDER BY sequence_number DESC

		-- name: read-journals
		SELECT * FROM journal
		WHERE persistence_id = ? AND sequence_number >= ? AND sequence_number <= ? AND deleted != 1
		ORDER BY sequence_number ASC

		-- name: delete-journals
		DELETE FROM journal
		WHERE persistence_id = ? AND sequence_number <= ?

		-- name: logical-delete-journals
		UPDATE journal
		SET deleted = 1
		WHERE persistence_id = ? AND sequence_number <= ?

		-- name: delete-snapshots
		DELETE FROM snapshot
		WHERE persistence_id = ? AND sequence_number <= ?
	`
)

// Driver defines a type of SQL driver accepted.
// This will be used by the golang sql library to load a specific driver
type Driver string

const (
	// POSTGRES driver type
	POSTGRES Driver = "postgres"
	// MYSQL driver type
	MYSQL Driver = "mysql"
	// ORACLE driver type
	ORACLE Driver = "oracle"
	// SQLSERVER driver type
	SQLSERVER Driver = "sqlserver"
)

// IsValid checks whether the given driver is valid or not
func (d Driver) IsValid() error {
	switch d {
	case POSTGRES, MYSQL, SQLSERVER, ORACLE:
		return nil
	}
	return errors.New("invalid driver type")
}

// String returns the actual value
func (d Driver) String() string {
	if err := d.IsValid(); err != nil {
		return ""
	}
	return string(d)
}

// ConnStr returns the connection string provided by the driver
func (d Driver) ConnStr(dbHost string, dbPort int, dbName, dbUser, dbPassword, dbSchema string) string {
	var connectionInfo string
	switch d {
	case POSTGRES:
		connectionInfo = fmt.Sprintf(
			"host=%s dbPort=%d user=%s dbname=%s sslmode=disable search_path=%s", dbHost, dbPort, dbUser, dbName,
			dbSchema,
		)
		// The POSTGRES driver gets confused in cases where the user has no password
		// set but a password is passed, so only set password if its non-empty
		if dbPassword != "" {
			connectionInfo += fmt.Sprintf(" password=%s", dbPassword)
		}
	case MYSQL:
		connectionInfo = fmt.Sprintf(
			"%s:%s@tcp(%s:%v)/%s", dbUser, dbPassword, dbHost, dbPort, dbName,
		)
	case SQLSERVER:
		connectionInfo = fmt.Sprintf(
			"server=%s;user id=%s;password=%s;dbPort=%d;database=%s;",
			dbHost, dbUser, dbPassword, dbPort, dbName,
		)
	}

	return connectionInfo
}

// SQLFile returns the sql file to create schema for a given driver
func (d Driver) SQLFile() string {
	switch d {
	case POSTGRES:
		return postgresSQL
	case MYSQL:
		return mysqlSQL
	case SQLSERVER:
		return sqlServerSQL
	}

	return ""
}
