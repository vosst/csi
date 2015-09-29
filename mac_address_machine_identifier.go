package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

// See include/uapi/linux/if_arp.h
const ethernetDeviceType = 1

type MACAddressMachineIdentifier struct {
	sysFsDirectory string // Base directory that is explored for available network interfaces.
}

func NewMACAddressMachineIdentifier() MACAddressMachineIdentifier {
	return MACAddressMachineIdentifier{"/sys/class/net"}
}

func (self MACAddressMachineIdentifier) Identify() ([]byte, error) {
	entries, err := ioutil.ReadDir(self.sysFsDirectory)

	if err != nil {
		return nil, err
	}

	for _, e := range entries {
		if !e.IsDir() {
			continue
		}

		t, err := os.Open(fmt.Sprintf("%s/%s/type", self.sysFsDirectory, e.Name()))
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

		a, err := os.Open(fmt.Sprintf("%s/%s/address", self.sysFsDirectory, e.Name()))

		if err != nil {
			continue
		}

		defer a.Close()

		return ioutil.ReadAll(a)
	}

	return nil, ErrFailedToIdentify
}
