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