package persistencesql

import (
	"context"

	"github.com/AsynkronIT/protoactor-go/actor"
)

// NewSQLServerProvider creates an instance postgres base SQLProvider
func NewSQLServerProvider(ctx context.Context, actorSystem *actor.ActorSystem, dbConfig DBConfig, opts ...OptFunc) (*SQLProvider, error) {
	dialect, err := NewSQLServerDialect(dbConfig)
	if err != nil {
		return nil, err
	}

	return NewSQLProvider(ctx, actorSystem, dialect, opts...), nil
}

// NewSQLServerDialect creates a new instance of SQLDialect
func NewSQLServerDialect(dialectConfig DBConfig) (SQLDialect, error) {
	dialect, err := NewDialect(dialectConfig, SQLSERVER)
	if err != nil {
		return nil, err
	}
	return dialect, nil
}
