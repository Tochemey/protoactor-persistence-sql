package persistencesql

import (
	"log"
	"time"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
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
	Timestamp int64
	// This persistent message's payload (the event).
	Payload []byte
	// A type hint for the event. This will be the proto message name of the event
	EventManifest Manifest
	// Unique identifier of the writing persistent actor.
	WriterID string
	// Flag to indicate the event has been deleted when logical deletion is set.
	Deleted bool
}

// NewJournal creates a new instance of Snapshot
func NewJournal(persistenceID string, message proto.Message, sequenceNumber int, writerID string) *Journal {
	manifest := proto.MessageName(message)
	bytes, err := proto.Marshal(message)
	if err != nil {
		log.Fatal(err)
	}

	return &Journal{
		PersistenceID:  persistenceID,
		SequenceNumber: sequenceNumber,
		Timestamp:      time.Now().UTC().Unix(),
		Payload:        bytes,
		EventManifest:  Manifest(manifest),
		WriterID:       writerID,
	}
}

func (journal *Journal) message() proto.Message {
	mt, err := protoregistry.GlobalTypes.FindMessageByName(protoreflect.FullName(journal.EventManifest))
	if err != nil {
		log.Fatal(err)
	}

	pm := mt.New().Interface()
	err = proto.Unmarshal(journal.Payload, pm)
	if err != nil {
		log.Fatal(err)
	}
	return pm
}
