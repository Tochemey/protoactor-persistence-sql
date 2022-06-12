package persistencesql

import (
	"context"

	"github.com/asynkron/protoactor-go/actor"
)

// NewPostgresProvider creates an instance postgres base SQLProvider
func NewPostgresProvider(ctx context.Context, actorSystem *actor.ActorSystem, dbConfig *DBConfig, opts ...OptFunc) (*SQLProvider, error) {
	dialect, err := NewPostgresDialect(dbConfig)
	if err != nil {
		return nil, err
	}

	// return the instance of the dialect
	return NewSQLProvider(ctx, actorSystem, dialect, opts...), nil
}

// NewPostgresDialect creates a new instance of SQLDialect
func NewPostgresDialect(dbConfig *DBConfig) (SQLDialect, error) {
	dialect, err := NewDialect(dbConfig, POSTGRES)
	if err != nil {
		return nil, err
	}
	return dialect, nil
}
