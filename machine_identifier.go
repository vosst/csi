package main

import "errors"

var ErrFailedToIdentify = errors.New("Failed to establish a unique ID of the device/machine.")

// MachineIdentifier abstracts ID generation for a device/machine.
type MachineIdentifier interface {
	// Identify returns a byte-slice containing the globally-unique
	// ID of the machine.
	Identify() ([]byte, error)
}
