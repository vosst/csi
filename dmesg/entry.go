package dmesg

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"syscall"
)

func facLev(v int) (Facility, Loglevel) {
	return MaskFacility(v), MaskLoglevel(v)
}

// dmesgLineRegExp parses an individual line from the kernel log buffer.
var dmesgLineRegExp = regexp.MustCompile(`<(\d+)>\[(\d+)\.(\d+)\](.*)`)

const (
	// Submatch index of the facility/level
	dmesgSmFacLev = 1
	// Submatch index of the timestamp, seconds part
	dmesgSmTsSec = 2
	// Submatch index of the timestamp, microseconds part
	dmesgSmTsUsec = 3
	// Submatch index of the actual message
	dmesgSmMsg = 4
	// Read all messages remaining in the ring buffer, placing then in the buffer pointed to  by  bufp. The
	// call reads the last len bytes from the log buffer (nondestructively), but will not read more than was
	// written into the buffer since the last "clear ring buffer" command (see command 5 below)). The call
	// returns the number of bytes read.
	sysActionReadAll int = 3
	// This command returns the total size of the kernel log buffer.
	sysActionSizeBuffer int = 10
)

// Entry models an individual log entry in the kernel ring buffer
type Entry struct {
	Level    Loglevel        // Loglevel of the entry
	Facility Facility        // Facility that the entry originated
	When     syscall.Timeval // Timestamp of the entry
	Message  string          // The actual log message
}

// ReadAll gathers all entries in the kernel log buffer nondestructively.
//
// Returns an error if a query to the underlying system facilities fails.
func ReadAll() ([]Entry, error) {
	n, err := syscall.Klogctl(10, nil)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to query size of log buffer [%s]", err))
	}

	b := make([]byte, n, n)

	m, err := syscall.Klogctl(3, b)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to read messages from log buffer [%s]", err))
	}

	entries := []Entry{}
	br := bufio.NewReader(bytes.NewReader(b[:m]))

	for line, err := br.ReadString('\n'); err == nil; line, err = br.ReadString('\n') {
		entry := Entry{}
		if matches := dmesgLineRegExp.FindStringSubmatch(line); len(matches) >= 4 {
			if fl, err := strconv.Atoi(matches[dmesgSmFacLev]); err == nil {
				entry.Facility, entry.Level = facLev(fl)
			}

			if s, err := strconv.Atoi(matches[dmesgSmTsSec]); err == nil {
				entry.When.Sec = int64(s)
			}

			if us, err := strconv.Atoi(matches[dmesgSmTsUsec]); err == nil {
				entry.When.Usec = int64(us)
			}

			entry.Message = matches[dmesgSmMsg]
			entries = append(entries, entry)
		}
	}

	return entries, nil
}
