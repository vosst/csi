package pid

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Cmdline describes the complete command line of a process, unless the process is a zombie.
type Cmdline []string

// NewCmdline reads the complete command line of the process identified by pid and returns
// the original command line or an error in case of issues.
func NewCmdline(pid int) (Cmdline, error) {
	fn := filepath.Join(Dir(pid), "cmdline")

	f, err := os.Open(fn)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to read %s [%s]", fn, err))
	}

	defer f.Close()

	return NewCmdlineFromReader(f), nil
}

// NewCmdlineFromReader parses all command line arguments from the given reader.
func NewCmdlineFromReader(reader io.Reader) Cmdline {
	br := bufio.NewReader(reader)
	cmdline := Cmdline{}

	for arg, err := br.ReadString('\x00'); err != nil; arg, err = br.ReadString('\x00') {
		cmdline = append(cmdline, arg)
	}

	return cmdline
}
