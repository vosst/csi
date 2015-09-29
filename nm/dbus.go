package nm

import "launchpad.net/go-dbus/v1"

var (
	Service   = "org.freedesktop.NetworkManager"
	Object    = dbus.ObjectPath("/org/freedesktop/NetworkManager")
	Interface = "org.freedesktop.NetworkManager"
)
