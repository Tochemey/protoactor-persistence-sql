package persistence

const (
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
)
