package pid

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
)

// State describes the state of a process
//
// Known values are presented in constants, taken from ${KERNELSRC}/fs/proc/array.c
type State string

// String pretty prints a process State instance
func (self State) String() string {
	switch self {
	case Running:
		return "Running"
	case Sleeping:
		return "Sleeping"
	case DiskSleep:
		return "DiskSleep"
	case Zombie:
		return "Zombie"
	case Stopped:
		return "Stopped"
	case TracingStop:
		return "TracingStop"
	case Dead:
		return "Dead"
	}

	return "Unknown"
}

const (
	Running     State = "R"
	Sleeping          = "S"
	DiskSleep         = "D"
	Zombie            = "Z"
	Stopped           = "T"
	TracingStop       = "t"
	Dead              = "X"
)

// Stat provides status information about the process as used by ps(1)
type Stat struct {
	Pid                 int    // The process id
	Comm                string // Filename of the executable in parentheses
	State               State  // State of the process
	Ppid                int    // The PID of the parent process
	Pgrp                int    // The process group ID of the process
	Session             int    // The session ID of the process
	TtyNr               int    // Controlling terminal of the process
	Tpgid               int    // ID of the foreground process of the controlling tmerinal of the process
	Flags               uint   // Kernel flags word of the process
	Minflt              uint   // Number of minor faults the process has made which have not required loading a memory page from disk
	Cminflt             uint   // Number of minor faults that the process's waited-for children have made
	Majflt              uint   // Number of major faults that the process's waited-for children have made
	Cmajflt             uint   // Number of major faults that process's waited-for children have made
	Utime               uint   // Amount of time that this process has been scheduled in user mode, in clock ticks
	Stime               uint   // Amount of time that this process has been scheduled in kernel mode, in clock ticks
	Cutime              int    // Amount of time that this process's waited-for children have been scheduled in kernel mode, in clock ticks
	Cstime              int    // Amount of time that this process's waited-for children have been scheduled in kernel mode, in clock ticks
	Priority            int    // Raw nice value as represented in the kernel or negated scheduling priority, minus one, for processes running a real-time scheduling policy
	Nice                int    // The nice value, in the range 19 (low priority) to -20 (high priority)
	NumThreads          int    // Number of threads in this process
	Itrealvalue         uint   // Time in jiffies before the next SIGALRM is sent to the process due to an interval timer
	StartTime           uint   // Time the process started after system boot.
	Vsize               uint   // Virtual memory size in bytes
	Rss                 uint   // Resident set size, number of pages the process has in real memory
	RssLim              uint   // Current soft limit in bytes on the rss of the process
	StartCode           uint   // Address above which program text can run
	EndCode             uint   // Address below which program text can run
	StartStack          uint   // Address of the start (i.e. bottom) of the stack
	Kstkesp             uint   // Current value of ESP (stack pointer), as found in the kernel stack page for the process
	Kstkeip             uint   // The current EIP (instruction pointer)
	Signal              uint   // The bitmap of pending signals, displayed as a decimal number.
	Blocked             uint   // The bitmap of blocked signals, displayed as a decimal number.
	SigIgnore           uint   // The bitmap of ignored signals, displayed as a decimal number.
	SigCatch            uint   // The bitmap of caught signals, displayed as a decimal number.
	Wchan               uint   // The channel in which the process is waiting, where channel is the address of a system call.
	Nswap               uint   // Number of pages swapped.
	Cnswap              uint   // Cumulative nswap for child processes
	ExitSignal          int    // Signal to be sent to parent when we die
	Processor           int    // CPU number last executed on
	RtPriority          uint   // Real-time scheduling priority, in the range 1-99 for processes scheduled under a real-time policy, or 0, for non-real-time processes
	Policy              uint   // Scheduling policy
	DelayacctBlkioTicks uint   // Aggregated block I/O delays, measured in ticks
	GuestTime           uint   // Guest time of the process (time spent running a virtual CPU for a gust OS), in clock ticks
	CguestTime          int    // Guest time of the process's children, measured in clock ticks
}

// NewStat reads /proc/%{pid}/stat into a Stat instance.
//
// Returns an error if opening /proc/%{pid}/stat or parsing an individual value fails.
func NewStat(pid int) (*Stat, error) {
	fn := filepath.Join(Dir(pid), "stat")

	f, err := os.Open(fn)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to read %s [%s]", fn, err))
	}

	defer f.Close()

	return NewStatFromReader(f)
}

// NewStatFromReader parses a Stat instance from the given reader.
//
// Returns an error if parsing an individual value fails.
func NewStatFromReader(reader io.Reader) (*Stat, error) {
	stat := Stat{}

	// We rely on reflection to step through the individual elements
	// of Stat and scan them from reader one by one. Fortunately, Fscan is
	// clever enough to do the right thing for numerical integer values.
	v := reflect.ValueOf(&stat).Elem()

	// We need the type later on to provide a rich error message in case
	// scanning an individual value fails.
	t := reflect.TypeOf(stat)

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

	return &stat, nil
}
