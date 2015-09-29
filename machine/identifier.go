package machine

import "errors"

var ErrFailedToIdentify = errors.New("Failed to establish a unique ID of the device/machine.")

// Identifier abstracts ID generation for a device/machine.
type Identifier interface {
	// Identify returns a byte-slice containing the globally-unique
	// ID of the machine.
	Identify() ([]byte, error)
}
