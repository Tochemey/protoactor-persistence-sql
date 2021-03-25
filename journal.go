package persistence

import (
	"log"
	"time"

	"github.com/golang/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/dynamicpb"
)

// Journal defines the journal row
type Journal struct {
	// the unique id of the journal row
	Ordering int64
	// Persistent ID that journals a persistent message.
	PersistenceID string
	// This persistent message's sequence number
	SequenceNumber int
	// The `timestamp` is the time the event was stored, in milliseconds since midnight, January 1, 1970 UTC.
	Timestamp time.Time
	// This persistent message's payload (the event).
	Payload []byte
	// A type hint for the event. This will be the proto message name of the event
	EventManifest string
	// Unique identifier of the writing persistent actor.
	WriterID string
	// Flag to indicate the event has been deleted when logical deletion is set.
	Deleted bool
}

// NewJournal creates a new instance of Snapshot
func NewJournal(persistenceID string, message proto.Message, sequenceNumber int, writerID string) *Journal {
	manifest := protoreflect.MessageDescriptor.FullName(message)
	bytes, err := proto.Marshal(message)
	if err != nil {
		log.Fatal(err)
	}

	return &Journal{
		Ordering:       0,
		PersistenceID:  persistenceID,
		SequenceNumber: sequenceNumber,
		Timestamp:      time.Now().UTC(),
		Payload:        bytes,
		EventManifest:  string(manifest),
		WriterID:       writerID,
	}
}

func (journal *Journal) message() proto.Message {
	t, err := protoregistry.GlobalTypes.FindMessageByName(protoreflect.FullName(journal.EventManifest))
	if err != nil {
		log.Fatal(err)
	}

	message := dynamicpb.NewMessage(t.Descriptor())
	err = proto.Unmarshal(journal.Payload, message)
	if err != nil {
		log.Fatal(err)
	}
	return message
}
