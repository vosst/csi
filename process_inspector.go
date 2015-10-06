package csi

import (
	"github.com/vosst/csi/proc/pid"
)

// ProcessReport bundles information about an individual process.
type ProcessReport struct {
	Cmdline pid.Cmdline // Command line
	Cwd     pid.Cwd     // Current working directory
	Env     pid.Environ // Runtime environment
}
