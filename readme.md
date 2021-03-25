# ProtoActor Persistence SQL

[![Go Report Card](https://goreportcard.com/badge/github.com/Tochemey/protoactor-persistence-sql)](https://goreportcard.com/report/github.com/Tochemey/protoactor-persistence-sql)
[![Build Status](https://www.travis-ci.com/Tochemey/protoactor-persistence-sql.svg?branch=master)](https://www.travis-ci.com/Tochemey/protoactor-persistence-sql)
[![codecov](https://codecov.io/gh/Tochemey/protoactor-persistence-sql/branch/master/graph/badge.svg?token=HVCXK21FQU)](https://codecov.io/gh/Tochemey/protoactor-persistence-sql)

An implementation of the ProtoActor persistence plugin APIs using RDBMS. It writes journal and snapshot to a configured
sql datastore.

## Features

The following sql-based datastore are supported out-of-the-box:

- MySQL
- Microsoft SQLServer
- Postgres

## TODOs

- [ ] Documentation
- [ ] More unit tests
- [ ] Code cleanup
- [ ] Cut an official release tag