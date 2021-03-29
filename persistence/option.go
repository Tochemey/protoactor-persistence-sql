package persistence

// Option represents the various option to set
type Option struct {
	// states whether the events deletion should be soft or not.
	// when this value is set to true the deleted flag will be set during event deletion
	// otherwise the event will be erased from the journal
	LogicalDeletion bool

	// Set the snapshot interval
	SnapshotInterval int
}
