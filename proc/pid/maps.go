package pid

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const (
	mapsAddressBegin  = 2
	mapsAddressEnd    = 3
	mapsReadPerm      = 5
	mapsWritePerm     = 6
	mapsExecPerm      = 7
	mapsPrivatePerm   = 8
	mapsOffset        = 9
	mapsDevMajor      = 11
	mapsDevMinor      = 12
	mapsInode         = 13
	mapsPath          = 14
	mapsSubmatchCount = 15
)

// mapsRegExp parses an individual line from /proc/%{pid}/maps. Please see
// https://regex101.com/r/cD2tN2/1 for a more readable overview together with
// an example. Index constants are easily verifiable over there, too.
var mapsRegExp = regexp.MustCompile(`(([[:xdigit:]]+)\-([[:xdigit:]]+))\s+((\-|r)(\-|w)(\-|x)(\-|p))\s+([[:xdigit:]]+)\s+(([[:digit:]]+)\:([[:digit:]]+))\s+([[:digit:]]+)\s+(.+)`)

// MemoryRegion describes a mapped memory region and its access permissions
type MemoryRegion struct {
	Address struct { // The addresses bounding the memory region
		Begin int64 // Begin of the memory region
		End   int64 // End of the memroy region
	}

	Permissions struct { // Access permissions on the memory region
		Read    bool // Region can be read
		Write   bool // Region can be written
		Exec    bool // Region can be executed
		Private bool // Region is private (copy on write)
	}
	Offset int64    // Offset into the file/device/whatever backing the memory region
	Device struct { // Device backing the memory region
		Major int // Major identifier
		Minor int // Minor identifier
	}
	Inode int    // Inode on the device backing the memory region
	Path  string // Usually the file backing the memory mapping
}

// Maps is the set of all mapped memory regions of a process
type Maps []MemoryRegion

// NewMaps reads the memory mappings from /proc/pid/maps, returning a Maps instance
// and an error in case of issues.
func NewMaps(pid int) (Maps, error) {
	fn := filepath.Join(Dir(pid), "maps")

	f, err := os.Open(fn)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to read %s [%s]", fn, err))
	}

	defer f.Close()

	return NewMapsFromReader(f)
}

// NewMapsFromReader parses memory mappings from reader.
//
// Returns all parsed memory mappings and an error. If the error
// is nil, all memory mappings were successfully parsed.
func NewMapsFromReader(reader io.Reader) (Maps, error) {
	maps := Maps{}
	br := bufio.NewReader(reader)

	for line, err := br.ReadString('\n'); err == nil; line, err = br.ReadString('\n') {
		tokens := mapsRegExp.FindStringSubmatch(strings.TrimRight(line, "\n"))

		if len(tokens) >= mapsSubmatchCount {
			mr := MemoryRegion{}

			mr.Address.Begin, _ = strconv.ParseInt(tokens[mapsAddressBegin], 16, 64)
			mr.Address.End, _ = strconv.ParseInt(tokens[mapsAddressEnd], 16, 64)

			mr.Permissions.Read = tokens[mapsReadPerm] == "r"
			mr.Permissions.Write = tokens[mapsWritePerm] == "w"
			mr.Permissions.Exec = tokens[mapsExecPerm] == "x"
			mr.Permissions.Private = tokens[mapsPrivatePerm] == "p"

			mr.Offset, _ = strconv.ParseInt(tokens[mapsOffset], 16, 64)

			mr.Device.Major, _ = strconv.Atoi(tokens[mapsDevMajor])
			mr.Device.Minor, _ = strconv.Atoi(tokens[mapsDevMinor])

			mr.Inode, _ = strconv.Atoi(tokens[mapsInode])
			mr.Path = tokens[mapsPath]

			maps = append(maps, mr)
		}
	}

	return maps, nil
}
