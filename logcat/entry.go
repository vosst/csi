// Package logcat provides types and functions to read from
// android's logcat facilities.
package logcat

import (
	"errors"
	"fmt"
	"io"
	"net"
	"path/filepath"
	"syscall"
	"time"
)

// A Logger refers to a specific log stream.
type Logger string

const (
	readTimeout         = 500 * time.Millisecond // Timeout when reading from the underlying Unix socket
	maxEntrySize        = 5 * 1024               // Maximum size of a single entry
	RadioLogger  Logger = "radio"                // radio-related messages
	EventsLogger        = "events"               // system/hardware events
	SystemLogger        = "system"               // system/framework messages
	MainLogger          = "main"                 // everything else
)

// An Entry is the userspace structure for version 1 of the logger_entry ABI.
type Entry struct {
	Pid     int             // Id of the process that generated the entry.
	Tid     int             // Id of the thread that generated the entry.
	When    syscall.Timeval // Timestamp of the entry.
	Message string          // Actual message payload.
}

// ReadOne reads the next logcat entry from reader.
//
// Returns an error if reading from reader fails.
func ReadOne(reader io.Reader) ([]byte, error) {
	buf := make([]byte, 5*1024, 5*1024)
	n, err := reader.Read(buf)
	if err != nil {
		return nil, err
	}

	return buf[:n], nil
}

// ReadAll returns all entries buffered in logger.
//
// Returns an error if connecting to logger fails.
func ReadAll(logger Logger) ([]byte, error) {
	fn := filepath.Join("/dev", "alog", string(logger))

	addr, err := net.ResolveUnixAddr("unixpacket", fn)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to resolve unix address %s [%s]", fn, err))
	}

	s, err := net.DialUnix("unixpacket", nil, addr)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to connect to %s [%s]", fn, err))
	}

	defer s.Close()

	s.SetReadDeadline(time.Now().Add(readTimeout))

	all := []byte{}

	for b, err := ReadOne(s); err == nil; b, err = ReadOne(s) {
		all = append(all, b...)
	}

	return all, nil
}
