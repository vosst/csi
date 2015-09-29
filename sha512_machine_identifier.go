package main

import (
	"crypto/sha512"
)

// SHA512MachineIdentifier computes the SHA512 hash of the result
// of a MachineIdentifier.Identify call.
type SHA512MachineIdentifier struct {
	Identifier MachineIdentifier
}

// Identify calls into the contained MachineIdentifier instance and hashes
// the result if the call to the inner MachineIdentifier succeeds.
func (self SHA512MachineIdentifier) Identify() ([]byte, error) {
	b, err := self.Identifier.Identify()
	if err != nil {
		return nil, err
	}

	hash := sha512.New()

	if _, err = hash.Write(b); err != nil {
		return nil, err
	}

	return hash.Sum(nil), nil
}
