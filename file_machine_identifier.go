package main

import (
	"io/ioutil"
	"os"
)

// FileMachineIdentifier reads an ID from a file.
type FileMachineIdentifier struct {
	Path string // Path to the file containing the ID.
}

// Identify tries to read the machine/device ID from the file
// under Path.
func (self FileMachineIdentifier) Identify() ([]byte, error) {
	f, err := os.Open(self.Path)

	if err != nil {
		return nil, err
	}

	defer f.Close()

	return ioutil.ReadAll(f)

}
