package persistence

import (
	"context"
	"log"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/persistence"
)

type SqlProvider struct {
	writer  *actor.PID
	dialect SqlDialect

	ctx context.Context
	cfg Option
}

// NewSqlProvider creates a new instance of the SqlProvider
func NewSqlProvider(ctx context.Context, actorSystem *actor.ActorSystem, dialect SqlDialect, option Option) *SqlProvider {
	// let us get the sql dialect connected
	if err := dialect.Connect(ctx); err != nil {
		log.Fatalf("error connecting: %v", err)
	}

	// let us create the various schemas required for the persistence provider
	// to be up running
	if err := dialect.CreateSchemasIfNotExist(ctx); err != nil {
		log.Fatalf("error creating schemas: %v", err)
	}

	pid := actorSystem.Root.Spawn(actor.PropsFromFunc(newWriter()))

	// create a new instance of the SqlProvider and returns it
	return &SqlProvider{
		writer:  pid,
		dialect: dialect,
		ctx:     ctx,
		cfg:     option,
	}
}

func (p *SqlProvider) GetState() persistence.ProviderState {
	return &SqlProviderState{
		SqlProvider: p,
	}
}
