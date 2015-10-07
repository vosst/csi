package csi

import (
	"fmt"
	"github.com/vosst/csi/proc/pid"
)

// ProcessReport bundles information about an individual process.
type ProcessReport struct {
	Cmdline pid.Cmdline // Command line
	Cwd     pid.Cwd     // Current working directory
	Env     pid.Environ // Runtime environment
	Exe     pid.Exe     // Path to executed command
	IO      pid.IO      // IO statistics
	Limits  pid.Limits  // Resource limits
	Maps    pid.Maps    // Mapped memory regions of the process
	OomAdj  pid.OomAdj  // OomAdj factor for altering the kernel's badness heuristic
}

type ProcessInspector struct {
}

func (self ProcessInspector) Inspect(id int) (*ProcessReport, error) {
	pr := ProcessReport{}

	if cl, err := pid.NewCmdline(id); err != nil {
		return nil, err
	} else {
		pr.Cmdline = cl
	}

	if cwd, err := pid.NewCwd(id); err != nil {
		return nil, err
	} else {
		pr.Cwd = cwd
	}

	if env, err := pid.NewEnviron(id); err != nil {
		return nil, err
	} else {
		pr.Env = env
	}

	if exe, err := pid.NewExe(id); err != nil {
		return nil, err
	} else {
		pr.Exe = exe
	}

	if io, err := pid.NewIO(id); err != nil {
		return nil, err
	} else {
		pr.IO = *io
	}

	if limits, err := pid.NewLimits(id); err != nil {
		fmt.Println(err)
		return nil, err
	} else {
		pr.Limits = limits
	}

	if maps, err := pid.NewMaps(id); err != nil {
		return nil, err
	} else {
		pr.Maps = maps
	}

	if oomAdj, err := pid.NewOomAdj(id); err != nil {
		return nil, err
	} else {
		pr.OomAdj = oomAdj
	}

	return &pr, nil
}
