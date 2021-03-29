package postgres

import (
	"context"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/tochemey/protoactor-persistence-sql/persistence"
)

// NewPostgresProvider creates an instance postgres base SQLProvider
func NewPostgresProvider(
	ctx context.Context, actorSystem *actor.ActorSystem, dialectConfig persistence.DialectConfig,
	option persistence.Option,
) (*persistence.SQLProvider, error) {

	dialect, err := persistence.NewDialect(dialectConfig, persistence.POSTGRES)
	if err != nil {
		return nil, err
	}

	return persistence.NewSQLProvider(ctx, actorSystem, dialect, option), nil
}
