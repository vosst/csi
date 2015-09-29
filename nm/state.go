package nm

const (
	// Networking state is unknown
	StateUnknown = iota
	// Networking is inactive and all devices are disabled.
	StateAsleep
	// There is no active network connection.
	StateDisconnected
	// Network connections are being cleaned up.
	StateDisconnecting
	// A network device is connecting to a network and there is no other available network connection.
	StateConnecting
	// A network device is connected, but there is only link-local connectivity.
	StateConnectedLocal
	// A network device is connected, but there is only site-local connectivity.
	StateConnectedSite
	// A network device is connected, with global network connectivity.
	StateConnectedGlobal
)

type State uint32
