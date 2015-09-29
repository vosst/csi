package machine

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

// DefaultIdentifier sets up the default machine identifier that is used
// for tagging uploads.
func DefaultIdentifier() (Identifier, error) {
	if err := ensureDataDir(); err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to create data dir: %s", err))
	}

	ami := FileIdentifier{"/sys/class/android_usb/android0/iSerial"}
	sumi := FileIdentifier{"/sys/class/dmi/id/product_uuid"}
	mmi := NewMACAddressIdentifier()
	smi := SHA512Identifier{
		ChainingIdentifier{
			mmi,
			ChainingIdentifier{
				sumi,
				ami}}}
	cmi := CachingIdentifier{smi, "/var/lib/whoopsie", "whoopsie"}

	return cmi, nil
}
