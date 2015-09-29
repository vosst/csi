package machine

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// See include/uapi/linux/if_arp.h
const ethernetDeviceType = 1

type MACAddressIdentifier struct {
	sysFsDirectory string // Base directory that is explored for available network interfaces.
}

func NewMACAddressIdentifier() MACAddressIdentifier {
	return MACAddressIdentifier{"/sys/class/net"}
}

func (self MACAddressIdentifier) Identify() ([]byte, error) {
	entries, err := ioutil.ReadDir(self.sysFsDirectory)

	if err != nil {
		return nil, err
	}

	for _, e := range entries {
		if !e.IsDir() {
			continue
		}

		t, err := os.Open(filepath.Join(self.sysFsDirectory, e.Name(), "type"))
		if err != nil {
			continue
		}

		defer t.Close()

		deviceType := int32(-1)
		_, err = fmt.Fscan(t, deviceType)

		if err != nil {
			continue
		}

		if deviceType != ethernetDeviceType {
			continue
		}

		a, err := os.Open(filepath.Join(self.sysFsDirectory, e.Name(), "address"))

		if err != nil {
			continue
		}

		defer a.Close()

		return ioutil.ReadAll(a)
	}

	return nil, ErrFailedToIdentify
}
