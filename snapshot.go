package persistencesql

import (
	"log"
	"time"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

// Snapshot defines the snapshot row
type Snapshot struct {
	// Persistent ID that journals a persistent message.
	PersistenceID string
	// This persistent message's sequence number
	SequenceNumber int
	// The `timestamp` is the time the event was stored, in milliseconds since midnight, January 1, 1970 UTC.
	Timestamp int64
	// This snapshot message's payload.
	Snapshot []byte
	// A type hint for the snapshot. This will be the proto message name of the snapshot
	SnapshotManifest Manifest
	// Unique identifier of the writing persistent actor.
	WriterID string
}

// NewSnapshot creates a new instance of Snapshot
func NewSnapshot(persistenceID string, message proto.Message, sequenceNumber int, writerID string) *Snapshot {
	manifest := proto.MessageName(message)
	bytes, err := proto.Marshal(message)
	if err != nil {
		log.Fatal(err)
	}

	return &Snapshot{
		PersistenceID:    persistenceID,
		SequenceNumber:   sequenceNumber,
		Timestamp:        time.Now().UTC().Unix(),
		Snapshot:         bytes,
		SnapshotManifest: Manifest(manifest),
		WriterID:         writerID,
	}
}

func (snapshot *Snapshot) message() proto.Message {
	mt, err := protoregistry.GlobalTypes.FindMessageByName(protoreflect.FullName(snapshot.SnapshotManifest))
	if err != nil {
		log.Fatal(err)
	}

	pm := mt.New().Interface()
	err = proto.Unmarshal(snapshot.Snapshot, pm)
	if err != nil {
		log.Fatal(err)
	}
	return pm
}
