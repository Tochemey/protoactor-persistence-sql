package persistencesql

import (
	"context"
	"math"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	pb "github.com/tochemey/protoactor-persistence-sql/gen"
	"google.golang.org/protobuf/proto"
)

func TestPostgresConnection(t *testing.T) {
	ctx := context.TODO()
	testCases := map[string]struct {
		config      *DBConfig
		expectError bool
	}{
		"valid connection settings": {
			config: NewDBConfig(
				"test",
				"test",
				"testdb",
				"public",
				"localhost",
				postgresContainerPort,
				WithConnectionMaxLife(maxConnectionLifetime),
			),
			expectError: false,
		},
		"database does not exist": {
			config: NewDBConfig(
				"test",
				"test",
				"test",
				"public",
				"localhost",
				postgresContainerPort,
				WithConnectionMaxLife(maxConnectionLifetime),
			),
			expectError: true,
		},
		"authentication failed": {
			config: NewDBConfig(
				"some-username",
				"some-password",
				"testdb",
				"public",
				"localhost",
				postgresContainerPort,
				WithConnectionMaxLife(maxConnectionLifetime),
			),
			expectError: true,
		},
	}

	for testName, testCase := range testCases {
		t.Run(
			testName, func(t *testing.T) {
				// get instance of assert
				assertions := assert.New(t)

				// create the dialect instance
				dialect, err := NewPostgresDialect(testCase.config)
				assertions.NoError(err)
				assertions.NotNil(dialect)

				err = dialect.Connect(ctx)
				switch testCase.expectError {
				case false:
					assertions.NoError(err)
				default:
					assertions.Error(err)
				}
			},
		)
	}
}

func TestPostgresSQLDialect(t *testing.T) {
	ctx := context.TODO()
	numEvents := 10
	numSnapshots := 3
	persistenceID := uuid.New().String()

	// set the database config
	config := NewDBConfig(
		"test",
		"test",
		"testdb",
		"public",
		"localhost",
		postgresContainerPort,
		WithConnectionMaxLife(maxConnectionLifetime),
	)

	// get instance of assert
	assertions := assert.New(t)
	// create the postgresDialect instance
	postgresDialect, err := NewPostgresDialect(config)
	assertions.NoError(err)
	assertions.NotNil(postgresDialect)

	// connect to the database
	err = postgresDialect.Connect(ctx)
	assertions.NoError(err)

	// create the journal and snapshot table successfully
	err = postgresDialect.CreateSchemasIfNotExist(ctx)
	assertions.NoError(err)

	// check whether both tables have been created
	err = tableExist(postgresHandle, POSTGRES, "public", "journal")
	assertions.NoError(err)
	assertions.Nil(err)
	err = tableExist(postgresHandle, POSTGRES, "public", "snapshot")
	assertions.NoError(err)
	assertions.Nil(err)

	// insert events into the journal store
	for i := 0; i < numEvents; i++ {
		persistenceID := uuid.New().String()
		journal := NewJournal(persistenceID, &pb.AccountDebited{
			AccountNumber: persistenceID,
			Balance:       float32(i * 100),
		}, i+1, "some-actor-pid")

		err = postgresDialect.PersistJournal(ctx, journal)
		assertions.NoError(err)
		assertions.Nil(err)
	}

	// insert some data into the snapshot store
	for i := 0; i < numSnapshots; i++ {
		snapshot := NewSnapshot(persistenceID, &pb.Account{
			AccountNumber: persistenceID,
			ActualBalance: float32(i * 100),
		}, i+1, "some-actor-pid")

		err = postgresDialect.PersistSnapshot(ctx, snapshot)
		assertions.NoError(err)
		assertions.Nil(err)
	}

	// let us count the number of elements in the journal and snapshot
	count := countJournal(postgresHandle)
	assertions.Equal(numEvents, count)
	count = countSnapshot(postgresHandle)
	assertions.Equal(numSnapshots, count)

	// let us fetch the latest snapshot for the given persistenceId
	// and perform some assertions
	latest, err := postgresDialect.GetLatestSnapshot(ctx, persistenceID)
	assertions.NoError(err)
	assertions.Equal(latest.SequenceNumber, 3)
	assertions.Equal(string(latest.SnapshotManifest), string(proto.MessageName(&pb.Account{})))
	snapshot, ok := (latest.message()).(*pb.Account)
	assertions.True(ok)
	assertions.Equal(snapshot.ActualBalance, float32(200))

	// let fetch some events from the journal store
	for i := 0; i < numEvents; i++ {
		journal := NewJournal(persistenceID, &pb.AccountDebited{
			AccountNumber: persistenceID,
			Balance:       float32(i * 100),
		}, i+1, "some-actor-pid")

		err = postgresDialect.PersistJournal(ctx, journal)
		assertions.NoError(err)
		assertions.Nil(err)
	}

	journals, err := postgresDialect.GetJournals(ctx, persistenceID, 2, 6)
	assertions.NoError(err)
	assertions.NotNil(journals)
	assertions.Equal(len(journals), 5)

	// delete some events from the journal
	err = postgresDialect.DeleteJournals(ctx, persistenceID, 2, true)
	assertions.NoError(err)

	// check the number of events remaining for the given persistence ID
	journals, err = postgresDialect.GetJournals(ctx, persistenceID, 1, math.MaxInt32)
	assertions.NoError(err)
	assertions.NotNil(journals)
	assertions.Equal(len(journals), 8)

	// delete some snapshots
	err = postgresDialect.DeleteSnapshots(ctx, persistenceID, 2)
	assertions.NoError(err)

}
