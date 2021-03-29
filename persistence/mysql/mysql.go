package mysql

import (
	"context"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/tochemey/protoactor-persistence-sql/persistence"
)

// NewMySQLProvider creates an instance postgres base SQLProvider
func NewMySQLProvider(ctx context.Context, actorSystem *actor.ActorSystem, dbConfig persistence.DBConfig, opts ...persistence.OptFunc) (*persistence.SQLProvider, error) {
	dialect, err := NewMySQLDialect(dbConfig)
	if err != nil {
		return nil, err
	}

	return persistence.NewSQLProvider(ctx, actorSystem, dialect, opts...), nil
}

// NewMySQLDialect creates a new instance of SQLDialect
func NewMySQLDialect(dbConfig persistence.DBConfig) (persistence.SQLDialect, error) {
	dialect, err := persistence.NewDialect(dbConfig, persistence.MYSQL)
	if err != nil {
		return nil, err
	}
	return dialect, nil
}
