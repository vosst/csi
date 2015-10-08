package pid

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const (
	OomAdjDefault    = 0   // Default adjustment value
	OomAdjDisableOOM = -17 // Disables OOM killing for a process
	OomAdjMin        = -16 // Minimum supported value
	OomAdjMax        = 15  // Maximum supported value
)

// OomAdj describes the factor for adjusting the score used to select which
// process should be killed in an out-of-memory (OOM) situation.
type OomAdj int

// NewOomAdj reads the OomAdj value for the process identified by pid
//
// Returns an error if /proc/%{pid}/oom_adj could not be opened for reading
// or if the value read from the value exceeds the documented bounds.
func NewOomAdj(pid int) (OomAdj, error) {
	fn := filepath.Join(Dir(pid), "oom_adj")

	f, err := os.Open(fn)

	if err != nil {
		return OomAdjDefault, errors.New(fmt.Sprintf("Failed to read %s [%s]", fn, err))
	}

	defer f.Close()

	return NewOomAdjFromReader(f)
}

// NewOomAdjFromReader reads the OomAdj value from the given stream
//
// Returns an error if the value exceeds the documented bounds.
func NewOomAdjFromReader(reader io.Reader) (OomAdj, error) {
	oomAdj := OomAdj(OomAdjDefault)

	if _, err := fmt.Fscanf(reader, "%d", &oomAdj); err != nil {
		return oomAdj, errors.New(fmt.Sprintf("Failed to parse OomAdj value [%s]", err))
	}

	if oomAdj < OomAdjDisableOOM || oomAdj > OomAdjMax {
		return oomAdj, errors.New(fmt.Sprintf("Invalid OomAdj value %d", oomAdj))
	}

	return oomAdj, nil
}
