package persistencesql

import (
	"testing"

	"github.com/stretchr/testify/assert"
	pb "github.com/tochemey/protoactor-persistence-sql/gen"
	"google.golang.org/protobuf/proto"
)

func TestJournal(t *testing.T) {
	// get instance of assert
	assertions := assert.New(t)

	// create an event to wrap into the journal
	event := &pb.AccountDebited{
		AccountNumber: "1234555",
		Balance:       2000,
	}

	journal := NewJournal("some-persistence-id", event, 1, "some-writer-id")

	assertions.Equal(journal.EventManifest, Manifest(proto.MessageName(event)))
	assertions.True(proto.Equal(journal.message(), event))
}
