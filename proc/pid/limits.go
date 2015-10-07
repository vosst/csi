package pid

// #include <sys/resource.h>
import "C"

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

// ResourceLimit describes a limit on a specific resource.
// Please see man limits
type ResourceLimit int

const (
	CpuTime        ResourceLimit = C.RLIMIT_CPU
	FileSize                     = C.RLIMIT_FSIZE
	DataSize                     = C.RLIMIT_DATA
	StackSize                    = C.RLIMIT_STACK
	CoreFileSize                 = C.RLIMIT_CORE
	ResidentSize                 = C.RLIMIT_RSS
	NoProcesses                  = C.RLIMIT_NPROC
	OpenFiles                    = C.RLIMIT_NOFILE
	LockedMemory                 = C.RLIMIT_MEMLOCK
	AddressSpace                 = C.RLIMIT_AS
	FileLocks                    = C.RLIMIT_LOCKS
	PendingSignals               = C.RLIMIT_SIGPENDING
	MsgQueueSize                 = C.RLIMIT_MSGQUEUE
	NicePrio                     = C.RLIMIT_NICE
	RealtimePrio                 = C.RLIMIT_RTTIME
)

func (self ResourceLimit) String() string {
	switch self {
	case CpuTime:
		return "CpuTime"
	case FileSize:
		return "FileSize"
	case DataSize:
		return "DataSize"
	case StackSize:
		return "StackSize"
	case CoreFileSize:
		return "CoreFileSize"
	case ResidentSize:
		return "ResidentSize"
	case NoProcesses:
		return "NoProcesses"
	case OpenFiles:
		return "OpenFiles"
	case LockedMemory:
		return "LockedMemory"
	case AddressSpace:
		return "AddressSpace"
	case FileLocks:
		return "FileLocks"
	case PendingSignals:
		return "PendingSignals"
	case MsgQueueSize:
		return "MsgQueueSize"
	case NicePrio:
		return "NicePrio"
	case RealtimePrio:
		return "RealtimePrio"
	}

	return fmt.Sprintf("Unknown [%d]", self)
}

// Unit describes the unit of a resouce limit.
type Unit string

const (
	Bytes        Unit = "bytes"
	Seconds           = "seconds"
	Processes         = "processes"
	Files             = "files"
	Signals           = "signals"
	Locks             = "locks"
	Microseconds      = "us"
	Unknown           = "unknown"

	// Parses a single line from /proc/pid/limits
	limitsSubmatchLimit     = 1
	limitsSubmatchSoftLimit = 2
	limitsSubmatchHardLimit = 3
	limitsSubmatchUnits     = 4
)

var (
	// See https://git.kernel.org/cgit/linux/kernel/git/torvalds/linux.git/tree/fs/proc/base.c#n599
	limitsKeys = map[string]ResourceLimit{
		"Max cpu time":          CpuTime,
		"Max file size":         FileSize,
		"Max data size":         DataSize,
		"Max stack size":        StackSize,
		"Max core file size":    CoreFileSize,
		"Max resident size":     ResidentSize,
		"Max processes":         NoProcesses,
		"Max open files":        OpenFiles,
		"Max locked memory":     LockedMemory,
		"Max address space":     AddressSpace,
		"Max file locks":        FileLocks,
		"Max pending signals":   PendingSignals,
		"Max msgqueue size":     MsgQueueSize,
		"Max nice priority":     NicePrio,
		"Max realtime priority": RealtimePrio,
	}
	// lineRegexp parses a single line from /proc/pid/limits.
	// Please see https://regex101.com/r/sK7nJ0/1 for details and an overview of submatches.
	limitsRegExp = regexp.MustCompile(`([[:print:]]+?)\s\s+([[:digit:]]+|unlimited)\s\s+([[:digit:]]+|unlimited)\s\s+([[:alpha:]]+)?`)
)

// Limit describes a limit on a resource
type Limit struct {
	Soft *int // Soft limit, nil indicates unlimited
	Hard *int // hard limit, nil indicates unlimited
	Unit Unit // Unit of measurement
}

// Limits describes all resource limits for a process
type Limits map[ResourceLimit]*Limit

// NewLimits parses the resource limits for the process identified by id
//
// Returns a Limits instance or an error if /proc/%{pid}/file could not be opened for reading
func NewLimits(pid int) (Limits, error) {
	fn := filepath.Join(Dir(pid), "limits")

	f, err := os.Open(fn)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to read %s [%s]", fn, err))
	}

	defer f.Close()

	return NewLimitsFromReader(f)
}

// NewLimitsFromReader parses resource limits from reader
//
// Returns a valid Limits instance or an error if parsing failed
func NewLimitsFromReader(reader io.Reader) (Limits, error) {
	limits := Limits{}
	br := bufio.NewReader(reader)

	for line, err := br.ReadString('\n'); err == nil; line, err = br.ReadString('\n') {
		tokens := limitsRegExp.FindStringSubmatch(line)
		if len(tokens) < 4 {
			// An error occured while parsing the line
			continue
		}

		if v, present := limitsKeys[tokens[limitsSubmatchLimit]]; present {
			limit := Limit{nil, nil, Unknown}

			if tokens[limitsSubmatchSoftLimit] != "unlimited" {
				if i, err := strconv.Atoi(tokens[limitsSubmatchSoftLimit]); err == nil {
					limit.Soft = &i
				}
			}

			if tokens[limitsSubmatchHardLimit] != "unlimited" {
				if i, err := strconv.Atoi(tokens[limitsSubmatchHardLimit]); err == nil {
					limit.Hard = &i
				}
			}

			if len(tokens) == 5 {
				limit.Unit = Unit(tokens[limitsSubmatchUnits])
			}

			limits[v] = &limit
		}
	}

	return limits, nil
}
