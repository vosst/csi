package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testDir = "/tmp"
const testFile = "testId"

var testFn = filepath.Join(testDir, testFile)

func TestCachingMachineIdentifierPrefersValueStoredInFile(t *testing.T) {
	os.Remove(testFn)

	f, err := os.Create(testFn)
	assert.Nil(t, err)

	binary.Write(f, binary.LittleEndian, uint64(42))
	f.Close()

	cmi := CachingMachineIdentifier{nil, testDir, testFile}

	id, err := cmi.Identify()

	assert.Nil(t, err)

	value := uint64(1)
	binary.Read(bytes.NewReader(id), binary.LittleEndian, &value)
	assert.Equal(t, uint64(42), value, "IDs do not match")
}

func TestCachingMachineIdentifierCallsIntoNext(t *testing.T) {
	os.Remove(testFn)

	mmi := MockMachineIdentifier{}
	mmi.On("Identify").Return([]byte{42}, nil)

	cmi := CachingMachineIdentifier{&mmi, testDir, testFile}

	cmi.Identify()

	mmi.AssertExpectations(t)
}

func TestCachingMachineIdentifierStoresResultOfCallToNext(t *testing.T) {
	os.Remove(testFn)

	mmi := MockMachineIdentifier{}
	mmi.On("Identify").Return([]byte{42, 42, 42}, nil)

	cmi := CachingMachineIdentifier{&mmi, testDir, testFile}

	cmi.Identify()

	mmi.On("Identiy").Return(nil, errors.New("Dummy"))

	id, _ := cmi.Identify()

	assert.Equal(t, []byte{42, 42, 42}, id, "Id mismatch")
}
