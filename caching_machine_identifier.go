package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

// CachineMachineIdentifier stores an id in a file,
// trying to load first prior to dispatching to a
// subsequent MachineIdentifier
type CachingMachineIdentifier struct {
	Identifier MachineIdentifier
	Dir        string
	File       string
}

// Identify tries to load the id from the configured file. If no id
// is stored, yet, it dispatches to a subsequent identifier, storing its
// result.
func (self CachingMachineIdentifier) Identify() ([]byte, error) {
	path := fmt.Sprintf("%s/%s", self.Dir, self.File)
	f, err := os.Open(path)

	if err == nil {
		defer f.Close()
		return ioutil.ReadAll(f)
	}

	b, err := self.Identifier.Identify()

	if err != nil {
		return nil, err
	}

	tmpPath := fmt.Sprintf("%s/whoopsie-temp", self.Dir)

	f, err = os.Create(tmpPath)

	if err != nil {
		return nil, err
	}

	_, err = f.Write(b)

	if err != nil {
		return nil, err
	}

	if err = os.Rename(tmpPath, path); err != nil {
		return nil, err
	}

	return b, nil
}
