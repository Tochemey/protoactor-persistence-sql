package persistence

import (
	"context"
	"log"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/persistence"
)

type SQLProvider struct {
	writer  *actor.PID
	dialect SQLDialect

	ctx context.Context
	cfg Option
}

// NewSQLProvider creates a new instance of the SQLProvider
func NewSQLProvider(
	ctx context.Context, actorSystem *actor.ActorSystem, dialect SQLDialect, option Option,
) *SQLProvider {
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
	return &SQLProvider{
		writer:  pid,
		dialect: dialect,
		ctx:     ctx,
		cfg:     option,
	}
}

// GetState returns an instance of the ProviderState
func (p *SQLProvider) GetState() persistence.ProviderState {
	return &SQLProviderState{
		SQLProvider: p,
	}
}
