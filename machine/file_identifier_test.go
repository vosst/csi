package machine

import (
	"bytes"
	"encoding/binary"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testFnFileIdentifier = "/tmp/test"

func TestFileIdentifierReadsFromConfiguredFile(t *testing.T) {
	os.Remove(testFnFileIdentifier)

	f, err := os.Create(testFnFileIdentifier)
	assert.Nil(t, err)

	assert.Nil(t, binary.Write(f, binary.LittleEndian, uint64(42)))
	f.Close()

	fi := FileIdentifier{testFnFileIdentifier}
	b, err := fi.Identify()

	assert.Nil(t, err)
	value := uint64(0)
	assert.Nil(t, binary.Read(bytes.NewReader(b), binary.LittleEndian, &value))

	assert.Equal(t, uint64(42), value, "Value mismatch")
}

func TestFileIdentifierPropagatesErrorWhenTryingToAccessFile(t *testing.T) {
	os.Remove(testFnFileIdentifier)

	fi := FileIdentifier{testFnFileIdentifier}
	b, err := fi.Identify()

	assert.Nil(t, b)
	assert.NotNil(t, err)
}
