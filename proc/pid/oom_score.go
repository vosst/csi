package pid

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// OomScore describes the current score that the kernel gives to this process for the
// purpose of selecting a process for the OOM-killer.
type OomScore int

// NewOomScore reads the OomScore value for the process identified by pid
//
// Returns an error if /proc/%{pid}/oom_score could not be opened for reading
func NewOomScore(pid int) (OomScore, error) {
	fn := filepath.Join(Dir(pid), "oom_score")

	f, err := os.Open(fn)

	if err != nil {
		return OomAdjDefault, errors.New(fmt.Sprintf("Failed to read %s [%s]", fn, err))
	}

	defer f.Close()

	return NewOomScoreFromReader(f)
}

// NewOomScoreFromReader reads the OomScore value from the given stream
//
// Returns an error if reading from the stream fails
func NewOomScoreFromReader(reader io.Reader) (OomScore, error) {
	oomScore := OomScore(0)

	if _, err := fmt.Fscanf(reader, "%d", &oomScore); err != nil {
		return oomScore, errors.New(fmt.Sprintf("Failed to parse OomScore value [%s]", err))
	}

	return oomScore, nil
}
