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
	LineRegExp        = `([[:print:]]+?)\s\s+([[:digit:]]+|unlimited)\s\s+([[:digit:]]+|unlimited)\s\s+([[:alpha:]]+)`
	SubmatchLimit     = 1
	SubmatchSoftLimit = 2
	SubmatchHardLimit = 3
	SubmatchUnits     = 4
)

var (
	keys = map[string]ResourceLimit{
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
	lineRegExp = regexp.MustCompile(`([[:print:]]+?)\s\s+([[:digit:]]+|unlimited)\s\s+([[:digit:]]+|unlimited)\s\s+([[:alpha:]]+)?`)
)

// Limit describes a limit on a resource
type Limit struct {
	Soft *int // Soft limit, nil indicates unlimited
	Hard *int // hard limit, nil indicates unlimited
	Unit Unit // Unit of measurement
}

type Limits map[ResourceLimit]*Limit

func NewLimits(pid int) (Limits, error) {
	fn := filepath.Join(Dir(pid), "limits")

	f, err := os.Open(fn)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to read %s [%s]", fn, err))
	}

	defer f.Close()

	return NewLimitsFromReader(f)
}

func NewLimitsFromReader(reader io.Reader) (Limits, error) {
	limits := Limits{}
	br := bufio.NewReader(reader)

	for line, err := br.ReadString('\n'); err == nil; line, err = br.ReadString('\n') {
		tokens := lineRegExp.FindStringSubmatch(line)
		if len(tokens) < 4 {
			// An error occured while parsing the line
			continue
		}

		if v, present := keys[tokens[SubmatchLimit]]; present {
			limit := Limit{nil, nil, Unknown}

			if tokens[SubmatchSoftLimit] != "unlimited" {
				if i, err := strconv.Atoi(tokens[SubmatchSoftLimit]); err == nil {
					limit.Soft = &i
				}
			}

			if tokens[SubmatchHardLimit] != "unlimited" {
				if i, err := strconv.Atoi(tokens[SubmatchHardLimit]); err == nil {
					limit.Hard = &i
				}
			}

			if len(tokens) == 5 {
				limit.Unit = Unit(tokens[SubmatchUnits])
			}

			limits[v] = &limit
		}
	}

	return limits, nil
}
