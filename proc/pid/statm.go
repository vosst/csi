package pid

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
)

// Provides information about memory usage, measured in pages.
type Statm struct {
	Size     uint // Total program size
	Resident uint // Resident set size
	Share    uint // Shared pages (i.e. backed by a file)
	Text     uint // Text (code) size
	Lib      uint // Library size
	Data     uint // Data + stack size
	Dt       uint // Dirty pages
}

// NewStatm reads /proc/%{pid}/statm into a Statm instance.
//
// Returns an error if opening /proc/%{pid}/statm or parsing an individual value fails.
func NewStatm(pid int) (*Statm, error) {
	fn := filepath.Join(Dir(pid), "statm")

	f, err := os.Open(fn)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to read %s [%s]", fn, err))
	}

	defer f.Close()

	return NewStatmFromReader(f)
}

// NewStatmFromReader parses a Statm instance from the given reader.
//
// Returns an error if parsing an individual value fails.
func NewStatmFromReader(reader io.Reader) (*Statm, error) {
	statm := Statm{}

	// We rely on reflection to step through the individual elements
	// of Statm and scan them from reader one by one. Fortunately, Fscan is
	// clever enough to do the right thing for numerical integer values.
	v := reflect.ValueOf(&statm).Elem()

	// We need the type later on to provide a rich error message in case
	// scanning an individual value fails.
	t := reflect.TypeOf(statm)

	// TODO(tvoss): Right now, we iterate over all fields in the struct, failing if any of those fails to be scanned
	// We might want to consider member annotations in the future to handle
	// optional values.
	for i := 0; i < v.NumField(); i++ {
		// We need a pointer to the field for passing it to Fscan.
		field := v.Field(i).Addr()
		if _, err := fmt.Fscan(reader, field.Interface()); err != nil {
			return nil, errors.New(fmt.Sprintf("Failed to parse field %s [%v]", t.Field(i).Name, err))
		}
	}

	return &statm, nil
}
