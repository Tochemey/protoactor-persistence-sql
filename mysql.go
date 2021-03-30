package persistencesql

import (
	"context"

	"github.com/AsynkronIT/protoactor-go/actor"
)

// NewMySQLProvider creates an instance postgres base SQLProvider
func NewMySQLProvider(ctx context.Context, actorSystem *actor.ActorSystem, dbConfig DBConfig, opts ...OptFunc) (*SQLProvider, error) {
	dialect, err := NewMySQLDialect(dbConfig)
	if err != nil {
		return nil, err
	}

	return NewSQLProvider(ctx, actorSystem, dialect, opts...), nil
}

// NewMySQLDialect creates a new instance of SQLDialect
func NewMySQLDialect(dbConfig DBConfig) (SQLDialect, error) {
	dialect, err := NewDialect(dbConfig, MYSQL)
	if err != nil {
		return nil, err
	}
	return dialect, nil
}
