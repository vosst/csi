package pid

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"syscall"
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

	stat, err := f.Stat()

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to stat %s [%s]", fn, err))
	}

	b, err := syscall.Mmap(int(f.Fd()), 0, int(stat.Size()), syscall.PROT_READ, syscall.MAP_PRIVATE)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to mmap %s [%s]", fn, err))
	}

	defer syscall.Munmap(b)

	return NewEnvironFromReader(bytes.NewReader(b)), nil
}

// NewEnvironFromReader reads the environment from reader.
func NewEnvironFromReader(reader io.Reader) Environ {
	br := bufio.NewReader(reader)

	env := Environ{}

	for arg, err := br.ReadString('\x00'); err != nil; arg, err = br.ReadString('\x00') {
		kv := strings.Split(arg, ":")
		env[kv[0]] = kv[1]
	}

	return env
}
