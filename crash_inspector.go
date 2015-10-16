package csi

import (
	"errors"
	"fmt"
	"github.com/vosst/csi/pkg/debian"
	"os"
)

// Crash report bundles all meta-data about a crashed process.
type CrashReport struct {
	Signal  os.Signal      // Signal that caused the crash
	System  *SystemReport  // Information about the overall system
	Process *ProcessReport // Information about the crashed process
}

// CrashInspector gathers information about a crash.
type CrashInspector struct {
}

// Inspect gathers information for a crashed process identfied by pid, recording the signal that caused the crash.
//
// Returns an error if either gathering system info or process-specific info fails.
func (self CrashInspector) Inspect(pid int, signal os.Signal) (*CrashReport, error) {
	si := SystemInspector{debian.NewSystem()}
	sr, err := si.Inspect()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to gather system information [%s]\n", err))
	}

	pi := ProcessInspector{debian.NewSystem()}
	pr, err := pi.Inspect(pid)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to gather process information [%s]\n", err))
	}

	return &CrashReport{signal, &sr, pr}, nil
}
