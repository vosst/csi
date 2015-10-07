package pid

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Environ describes the environment of a process
type Environ map[string]string

// NewEnviron loads the environment for the process with the given pid.
func NewEnviron(pid int) (Environ, error) {
	fn := filepath.Join(Dir(pid), "environ")

	f, err := os.Open(fn)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to read %s [%s]", fn, err))
	}

	defer f.Close()

	return NewEnvironFromReader(f), nil
}

// NewEnvironFromReader reads the environment from reader.
func NewEnvironFromReader(reader io.Reader) Environ {
	br := bufio.NewReader(reader)

	env := Environ{}

	for arg, err := br.ReadString('\x00'); err == nil; arg, err = br.ReadString('\x00') {
		if kv := strings.Split(arg, "="); len(kv) == 2 {
			env[kv[0]] = strings.TrimRight(kv[1], "\x00")
		}
	}

	return env
}
