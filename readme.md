# ProtoActor Persistence SQL

![Codacy grade](https://img.shields.io/codacy/grade/3e0d5b0d52cd4ef4943a9045375f216d?style=for-the-badge)
![GitHub Workflow Status (branch)](https://img.shields.io/github/workflow/status/tochemey/protoactor-persistence-sql/CI/master?style=for-the-badge)
![Codecov](https://img.shields.io/codecov/c/github/tochemey/protoactor-persistence-sql?style=for-the-badge)

An implementation of the ProtoActor persistence plugin APIs using RDBMS. It writes journal and snapshot to a configured
SQL datastore. At the moment the following data stores are supported out of the box:

- [MySQL](https://www.mysql.com/)
- [Postgres](https://www.postgresql.org/)

The events and state snapshots are protocol buffer bytes array persisted respectively in the journal and snapshot
tables.

## Journal schemas

### Postgres

```postgresql
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
```

### MySQL

```mysql
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
```

## Snapshot schemas

### Postgres

```postgresql
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
```

### MySQL

```mysql
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
```

