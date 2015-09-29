package main

// ReachbilityMonitorCreator abstracts creation of ReachabilityMonitor instances.
type ReachabilityMonitorCreator interface {
	// Create returns a new ReachabilityMonitor instance monitoring reachability
	// of the given host or IP. Returns an error if setup of the monitor fails.
	Create(hostOrIp string) (ReachabilityMonitor, error)
}
