package persistencesql

import (
	"testing"

	"github.com/stretchr/testify/assert"
	pb "github.com/tochemey/protoactor-persistence-sql/gen"
	"google.golang.org/protobuf/proto"
)

func TestSnapshot(t *testing.T) {
	// get instance of assert
	assertions := assert.New(t)

	// create an event to wrap into the journal
	event := &pb.Account{
		AccountNumber: "1234555",
		ActualBalance: 2000,
	}

	snapshot := NewSnapshot("some-persistence-id", event, 1, "some-writer-id")

	assertions.Equal(snapshot.SnapshotManifest, proto.MessageName(event))
	assertions.True(proto.Equal(snapshot.message(), event))
}
