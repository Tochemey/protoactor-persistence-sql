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