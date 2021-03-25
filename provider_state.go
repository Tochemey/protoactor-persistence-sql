package persistence

import (
	"log"
	"sync"

	"github.com/golang/protobuf/proto"
)

type SQLProviderState struct {
	*SQLProvider
	wg sync.WaitGroup
}

// GetSnapshot fetches the latest snapshot of a given persistenceID represented by the actorName
// actorName is the persistenceID
func (s *SQLProviderState) GetSnapshot(actorName string) (snapshot interface{}, eventIndex int, ok bool) {
	record, err := s.dialect.GetLatestSnapshot(s.ctx, actorName)
	if err != nil {
		log.Fatalf("error fetching snapshot: %v", err)
		return nil, 0, false
	}

	return record.message(), record.SequenceNumber, true
}

// PersistSnapshot saves the snapshot of a given persistenceID.
// actorName is the persistenceID
// snapshotIndex is the sequenceNumber of the snapshot data
// snapshot is the payload to persist
func (s *SQLProviderState) PersistSnapshot(actorName string, snapshotIndex int, snapshot proto.Message) {
	// let us convert the v1 proto to a v2 proto message

	newSnapshot := NewSnapshot(actorName, proto.MessageV2(snapshot), snapshotIndex, s.writer.Id)
	if err := s.dialect.PersistSnapshot(s.ctx, newSnapshot); err != nil {
		log.Fatalf(
			"error: %v persisting snapshot: %s for persistenceID: %s", err, newSnapshot.SnapshotManifest, actorName,
		)
	}
}

// DeleteSnapshots deletes snapshots for a given persistenceID from the store to a given sequenceNumber.
// actorName is the persistenceID
// inclusiveToIndex is the sequenceNumber
func (s *SQLProviderState) DeleteSnapshots(actorName string, inclusiveToIndex int) {
	if err := s.dialect.DeleteSnapshots(s.ctx, actorName, inclusiveToIndex); err != nil {
		log.Fatalf("error deleting snapshots: %v for persistenceID: %s", err, actorName)
	}
}

// GetEvents list events from the journal store within a range of sequenceNumber for a given persistence ID
// actorName is the persistenceID
// eventIndexStart is the from sequenceNumber
// eventIndexEnd is the to sequenceNumber
func (s *SQLProviderState) GetEvents(
	actorName string, eventIndexStart int, eventIndexEnd int, callback func(e interface{}),
) {
	events, err := s.dialect.GetJournals(s.ctx, actorName, eventIndexStart, eventIndexEnd)
	if err != nil {
		log.Fatalf("error fetching events: %v", err)
	}

	for _, e := range events {
		callback(e)
	}
}

// PersistEvent persists an event for a given persistence ID
// actorName is the persistenceID
// eventIndex is the event to persist sequenceNumber
// event is the event payload
func (s *SQLProviderState) PersistEvent(actorName string, eventIndex int, event proto.Message) {
	journal := NewJournal(actorName, proto.MessageV2(event), eventIndex, s.writer.Id)
	if err := s.dialect.PersistJournal(s.ctx, journal); err != nil {
		log.Fatalf("error: %v persisting event: %s for persistenceID: %s", err, journal.EventManifest, actorName)
	}
}

// DeleteEvents deletes events from journal to a given index
// actorName is the persistenceID
// inclusiveToIndex is the sequence Number
func (s *SQLProviderState) DeleteEvents(actorName string, inclusiveToIndex int) {
	if err := s.dialect.DeleteJournals(s.ctx, actorName, inclusiveToIndex, s.cfg.LogicalDeletion); err != nil {
		log.Fatalf("error deleting events: %v for persistenceID: %s", err, actorName)
	}
}

func (s *SQLProviderState) Restart() {
	// let us wait for any pending  writes to complete
	s.wg.Wait()
}

// GetSnapshotInterval return the snapshot interval
func (s *SQLProviderState) GetSnapshotInterval() int {
	return s.cfg.SnapshotInterval
}
