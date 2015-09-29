package main

import (
	"io/ioutil"
	"os"
)

// SystemUUIDMachineIdentifier reads the product id of the device/machine.
type SystemUUIDMachineIdentifier struct {
	sysFsFile string // Path to file in sysfs.
}

func NewSystemUUIDMachineIdentifier() SystemUUIDMachineIdentifier {
	return SystemUUIDMachineIdentifier{"/sys/class/dmi/id/product_uuid"}
}

// Identify tries to read the machine's uuid from /sys/class/dmi/id/product_uuid
func (self SystemUUIDMachineIdentifier) Identify() ([]byte, error) {
	f, err := os.Open(self.sysFsFile)

	if err != nil {
		return nil, err
	}

	defer f.Close()

	return ioutil.ReadAll(f)
}
