package crash

// ReportPersister abstracts persisting of Report instances.
type ReportPersister interface {
	// Persist stores report for future processing
	Persist(report Report) error
}
