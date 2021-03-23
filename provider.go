package persistence

import (
	"context"
	"log"
	"time"

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
func NewSqlProvider(actorSystem *actor.ActorSystem, dialect SqlDialect, cfg Option) *SqlProvider {
	// create a background context to pass it on
	ctx := context.Background()

	// let us get the sql dialect connected
	if err := dialect.Connect(); err != nil {
		log.Fatalf("error connecting: %v", err)
	}

	// let us create the various schemas required for the persistence provider
	// to be up running
	if err := dialect.CreateSchemasIfNotExist(); err != nil {
		log.Fatalf("error creating schemas: %v", err)
	}

	pid := actorSystem.Root.Spawn(actor.PropsFromFunc(newWriter(time.Second / 10000)))

	// create a new instance of the SqlProvider and returns it
	return &SqlProvider{
		writer:  pid,
		dialect: dialect,
		ctx:     ctx,
		cfg:     cfg,
	}
}

func (p *SqlProvider) GetState() persistence.ProviderState {
	return &SqlProviderState{
		SqlProvider: p,
	}
}
