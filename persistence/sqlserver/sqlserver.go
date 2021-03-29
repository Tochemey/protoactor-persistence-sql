package sqlserver

import (
	"context"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/tochemey/protoactor-persistence-sql/persistence"
)

// NewSQLServerProvider creates an instance postgres base SQLProvider
func NewSQLServerProvider(ctx context.Context, actorSystem *actor.ActorSystem, dbConfig persistence.DBConfig, opts ...persistence.OptFunc) (*persistence.SQLProvider, error) {
	dialect, err := NewSQLServerDialect(dbConfig)
	if err != nil {
		return nil, err
	}

	return persistence.NewSQLProvider(ctx, actorSystem, dialect, opts...), nil
}

// NewSQLServerDialect creates a new instance of SQLDialect
func NewSQLServerDialect(dialectConfig persistence.DBConfig) (persistence.SQLDialect, error) {
	dialect, err := persistence.NewDialect(dialectConfig, persistence.SQLSERVER)
	if err != nil {
		return nil, err
	}
	return dialect, nil
}
