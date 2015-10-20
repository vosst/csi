package log

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/vosst/csi/dmesg"
)

// A Collector handles gathering of all contents of a log facility.
type Collector interface {
	// Collect gathers a blob of bytes representing the contents of a specific log facility.
	//
	// Returns an error if snapshotting the underlying log facility fails.
	Collect() ([]byte, error)
}

// A DmesgCollector gathers the contents of the kernel log buffer.
type DmesgCollector struct {
}

func NewDmesgCollector() DmesgCollector {
	return DmesgCollector{}
}

// Collect returns the contents of the kernel log buffer.
//
// Returns an error if querying the kernel log buffer fails due to a lag of permissions.
func (d DmesgCollector) Collect() ([]byte, error) {
	blob, err := dmesg.ReadAll()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to collect contents of the kernel log buffer [%s]", err))
	}

	return blob, nil
}

// A SyslogCollector gathers the contents of the syslog.
type SyslogCollector struct {
	fn string // File containing the syslog
}

func NewSyslogCollector() SyslogCollector {
	return SyslogCollector{"/var/log/syslog"}
}

func (s SyslogCollector) Collect() ([]byte, error) {
	blob, err := ioutil.ReadFile(s.fn)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to collect syslog from %s [%s]", s.fn, err))
	}

	return blob, nil
}
