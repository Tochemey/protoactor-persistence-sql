package persistencesql

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"

	pb "github.com/tochemey/protoactor-persistence-sql/gen"
)

func TestJournal(t *testing.T) {
	// get instance of assert
	assertions := assert.New(t)

	// create an event to wrap into the journal
	event := &pb.AccountOpened{
		AccountNumber: "1234555",
		Balance:       2000,
	}

	journal := NewJournal("some-persistence-id", event, 1, "some-writer-id")

	assertions.Equal(journal.EventManifest, proto.MessageName(event))
	assertions.True(proto.Equal(journal.message(), event))
}
