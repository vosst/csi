package main

// CrashReportPersister abstracts persisting of CrashReport instances.
type CrashReportPersister interface {
	// Persist stores report for future processing
	Persist(report *CrashReport) error
}
