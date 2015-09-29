package main

import (
	"errors"
	"fmt"
	"os"
)

const dataDir = "/var/lib/csi"
const dataFile = "csi"

// ensureDataDir tries to create the dataDir if it does not exist.
func ensureDataDir() error {
	return os.MkdirAll(dataDir, os.ModeDir)
}

// DefaultMachineIdentifier sets up the default machine identifier that is used
// for tagging uploads.
func DefaultMachineIdentifier() (MachineIdentifier, error) {
	if err := ensureDataDir(); err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to create data dir: %s", err))
	}

	ami := FileMachineIdentifier{"/sys/class/android_usb/android0/iSerial"}
	sumi := FileMachineIdentifier{"/sys/class/dmi/id/product_uuid"}
	mmi := NewMACAddressMachineIdentifier()
	smi := SHA512MachineIdentifier{
		ChainingMachineIdentifier{
			mmi,
			ChainingMachineIdentifier{
				sumi,
				ami}}}
	cmi := CachingMachineIdentifier{smi, "/var/lib/whoopsie", "whoopsie"}

	return cmi, nil
}
