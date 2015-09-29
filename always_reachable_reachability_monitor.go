package main

// AlwaysReachableReachabilityMonitor is a dummy implementation always reporting IsReachable
// for every possible host.
type AlwaysReachableReachabilityMonitor struct {
}

// CheckHostReachability always returns IsReachable.
func (self AlwaysReachableReachabilityMonitor) CheckHostReachability(host string) Reachability {
	return IsReachable
}
