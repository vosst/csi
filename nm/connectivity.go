package nm

const (
	// Network connectivity is unknown.
	ConnectivityUnknown = iota
	// The host is not connected to any network.
	ConnectivityNone
	// The host is behind a captive portal and cannot reach the full Internet.
	ConnectivityPortal
	// The host is connected to a network, but does not appear to be able to reach the full Internet.
	ConnectivityLimited
	// The host is connected to a network, and appears to be able to reach the full Internet
	ConnectivityFull
)

type Connectivity uint32
