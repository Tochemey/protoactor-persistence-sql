package postgres

import (
	"context"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/tochemey/protoactor-persistence-sql/persistence"
)

// NewPostgresProvider creates an instance postgres base SQLProvider
func NewPostgresProvider(ctx context.Context, actorSystem *actor.ActorSystem, dbConfig persistence.DBConfig, opts ...persistence.OptFunc) (*persistence.SQLProvider, error) {
	dialect, err := NewPostgresDialect(dbConfig)
	if err != nil {
		return nil, err
	}

	return persistence.NewSQLProvider(ctx, actorSystem, dialect, opts...), nil
}

// NewPostgresDialect creates a new instance of SQLDialect
func NewPostgresDialect(dbConfig persistence.DBConfig) (persistence.SQLDialect, error) {
	dialect, err := persistence.NewDialect(dbConfig, persistence.POSTGRES)
	if err != nil {
		return nil, err
	}
	return dialect, nil
}
