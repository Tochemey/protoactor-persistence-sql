package persistencesql

import (
	"context"
	"log"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/persistence"
)

type OptFunc = func(provider *SQLProvider)

// SQLProvider defines a generic persistence provider.
// The type of provider is determined by the type of SQLDialect defined
type SQLProvider struct {
	writer  *actor.PID
	dialect SQLDialect

	ctx context.Context
	// states whether the events deletion should be soft or not.
	// when this value is set to true the deleted flag will be set during event deletion
	// otherwise the event will be erased from the journal
	logicalDeletion bool

	// Set the snapshot interval
	snapshotInterval int
}

// NewSQLProvider creates a new instance of the SQLProvider
func NewSQLProvider(
	ctx context.Context, actorSystem *actor.ActorSystem, dialect SQLDialect, opts ...OptFunc,
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

	// create a new instance of SQLProvider
	provider := new(SQLProvider)

	// call option functions on instance to set options on it
	for _, opt := range opts {
		opt(provider)
	}

	// set the provider
	provider.writer = pid
	provider.dialect = dialect
	provider.ctx = ctx

	// create a new instance of the SqlProvider and returns it
	return provider
}

// GetState returns an instance of the ProviderState
func (p *SQLProvider) GetState() persistence.ProviderState {
	return &SQLProviderState{
		SQLProvider: p,
	}
}

// WithLogicalDeletion enables logical deletion
func WithLogicalDeletion() OptFunc {
	return func(provider *SQLProvider) {
		provider.logicalDeletion = true
	}
}

// WithSnapshotInterval sets the snapshot interval
func WithSnapshotInterval(interval int) OptFunc {
	return func(provider *SQLProvider) {
		provider.snapshotInterval = interval
	}
}
