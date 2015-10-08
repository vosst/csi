package pid

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// IO summarizes I/O statistics for a process
type IO struct {
	RChar               int // Characters read
	WChar               int // Characters written
	Syscr               int // Syscalls to read
	Syscw               int // Syscalls to write
	ReadBytes           int // Bytes read
	WriteBytes          int // Bytes written
	CancelledWriteBytes int // See man proc
}

// NewIO determines the io statistics for the process identified by pid.
func NewIO(pid int) (*IO, error) {
	fn := filepath.Join(Dir(pid), "io")

	f, err := os.Open(fn)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to read %s [%s]", fn, err))
	}

	defer f.Close()
	io := NewIOFromReader(f)
	return io, nil
}

// NewIOFromReader reads I/O statistics from reader.
func NewIOFromReader(reader io.Reader) *IO {
	result := IO{}
	br := bufio.NewReader(reader)

	for line, err := br.ReadString('\n'); err == nil; line, err = br.ReadString('\n') {
		if kv := strings.Split(line, ":"); len(kv) == 2 {
			switch strings.TrimSpace(kv[0]) {
			case "rchar":
				result.RChar, _ = strconv.Atoi(strings.TrimSpace(kv[1]))
			case "wchar":
				result.WChar, _ = strconv.Atoi(strings.TrimSpace(kv[1]))
			case "syscr":
				result.Syscr, _ = strconv.Atoi(strings.TrimSpace(kv[1]))
			case "syscw":
				result.Syscw, _ = strconv.Atoi(strings.TrimSpace(kv[1]))
			case "read_bytes":
				result.ReadBytes, _ = strconv.Atoi(strings.TrimSpace(kv[1]))
			case "write_bytes":
				result.WriteBytes, _ = strconv.Atoi(strings.TrimSpace(kv[1]))
			case "cancelled_write_bytes":
				result.CancelledWriteBytes, _ = strconv.Atoi(strings.TrimSpace(kv[1]))
			}
		}
	}

	return &result
}
