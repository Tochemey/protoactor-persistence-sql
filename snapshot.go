package persistence

import (
	"log"
	"time"

	"github.com/golang/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/dynamicpb"
)

// Snapshot defines the snapshot row
type Snapshot struct {
	// Persistent ID that journals a persistent message.
	PersistenceID string
	// This persistent message's sequence number
	SequenceNumber int
	// The `timestamp` is the time the event was stored, in milliseconds since midnight, January 1, 1970 UTC.
	Timestamp time.Time
	// This snapshot message's payload.
	Snapshot []byte
	// A type hint for the snapshot. This will be the proto message name of the snapshot
	SnapshotManifest string
	// Unique identifier of the writing persistent actor.
	WriterID string
}

// NewSnapshot creates a new instance of Snapshot
func NewSnapshot(persistenceID string, message proto.Message, sequenceNumber int, writerID string) *Snapshot {
	manifest := protoreflect.MessageDescriptor.FullName(message)
	bytes, err := proto.Marshal(message)
	if err != nil {
		log.Fatal(err)
	}

	return &Snapshot{
		PersistenceID:    persistenceID,
		SequenceNumber:   sequenceNumber,
		Timestamp:        time.Now().UTC(),
		Snapshot:         bytes,
		SnapshotManifest: string(manifest),
		WriterID:         writerID,
	}
}

func (snapshot *Snapshot) message() proto.Message {
	t, err := protoregistry.GlobalTypes.FindMessageByName(protoreflect.FullName(snapshot.SnapshotManifest))
	if err != nil {
		log.Fatal(err)
	}

	message := dynamicpb.NewMessage(t.Descriptor())
	err = proto.Unmarshal(snapshot.Snapshot, message)
	if err != nil {
		log.Fatal(err)
	}
	return message
}
