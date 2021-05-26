## ProtoActor Persistence SQL

![Codacy grade](https://img.shields.io/codacy/grade/3e0d5b0d52cd4ef4943a9045375f216d?style=for-the-badge)
![GitHub Workflow Status (branch)](https://img.shields.io/github/workflow/status/tochemey/protoactor-persistence-sql/CI/master?style=for-the-badge)
![Codecov](https://img.shields.io/codecov/c/github/tochemey/protoactor-persistence-sql?style=for-the-badge)

An implementation of the ProtoActor persistence plugin APIs using RDBMS. It writes journal and snapshot to a configured
SQL datastore. At the moment the following data stores are supported out of the box:

- [MySQL](https://www.mysql.com/)
- [Postgres](https://www.postgresql.org/)

The events and state snapshots are protocol buffer bytes array persisted respectively in the journal and snapshot
tables.

Note: _The developer does not need to create the database tables. They are created by default by the library._
One can have a look at them in the _constants.go_ code.