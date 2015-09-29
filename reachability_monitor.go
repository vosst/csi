package main

const (
	// Indicates that a host is not reachable via the current networking setup
	NotReachable = 0
	// Indicates that a host is reachable via the current networking setup
	IsReachable = 1 << 0
	// Indicates that the route to a host goes via a wwan connection
	IsWWAN = 1 << 1
)

// Reachability is a flag field.
type Reachability uint32

// ReachabilityMonitor helps in monitoring the reachability of a specifc host
type ReachabilityMonitor interface {
	// CheckHostReachability checks whether the given host is
	// reachable via the current networking setup of the machine/device.
	CheckHostReachability(host string) Reachability
}
