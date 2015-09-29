package machine

import (
	"io/ioutil"
	"os"
)

// FileIdentifier reads an ID from a file.
type FileIdentifier struct {
	Path string // Path to the file containing the ID.
}

// Identify tries to read the machine/device ID from the file
// under Path.
func (self FileIdentifier) Identify() ([]byte, error) {
	f, err := os.Open(self.Path)

	if err != nil {
		return nil, err
	}

	defer f.Close()

	return ioutil.ReadAll(f)

}
