package pid

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// OomScoreAdj models the value used to adjust the badness heuristic
type OomScoreAdj int

const (
	OomScoreAdjDefault = 0     // Default value
	OomScoreAdjMin     = -1000 // Minimum supported value
	OomScoreAdjMax     = 1000  // Maximum supported value
)

// NewOomScoreAdj reads the OomScoreAdj value for the process identified by pid
//
// Returns an error if /proc/%{pid}/oom_score_adj could not be opened for reading
// or if the value read from the value exceeds the documented bounds.
func NewOomScoreAdj(pid int) (OomScoreAdj, error) {
	fn := filepath.Join(Dir(pid), "oom_score_adj")

	f, err := os.Open(fn)

	if err != nil {
		return OomScoreAdjDefault, errors.New(fmt.Sprintf("Failed to read %s [%s]", fn, err))
	}

	defer f.Close()

	return NewOomScoreAdjFromReader(f)
}

// NewOomScoreAdjFromReader reads the OomScoreAdj value from the given reader
//
// Returns an error if the value exceeds the documented bounds.
func NewOomScoreAdjFromReader(reader io.Reader) (OomScoreAdj, error) {
	oomScoreAdj := OomScoreAdj(OomScoreAdjDefault)

	if _, err := fmt.Fscanf(reader, "%d", &oomScoreAdj); err != nil {
		return oomScoreAdj, errors.New(fmt.Sprintf("Failed to parse OomScoreAdj value [%s]", err))
	}

	if oomScoreAdj < OomScoreAdjMin || oomScoreAdj > OomScoreAdjMax {
		return oomScoreAdj, errors.New(fmt.Sprintf("Invalid OomScoreAdj value %d", oomScoreAdj))
	}

	return oomScoreAdj, nil
}
