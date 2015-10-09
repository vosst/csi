package csi

import (
	"fmt"
	"github.com/vosst/csi/pkg"
	"github.com/vosst/csi/proc/pid"
)

// ProcessReport bundles information about an individual process.
type ProcessReport struct {
	Bundle pkg.Bundle // The package/bundle the executable executed in the process belongs to

	Cmdline     pid.Cmdline     // Command line
	Cwd         pid.Cwd         // Current working directory
	Env         pid.Environ     // Runtime environment
	Exe         pid.Exe         // Path to executed command
	Fd          pid.Fd          // All open fds
	IO          pid.IO          // IO statistics
	Limits      pid.Limits      // Resource limits
	Maps        pid.Maps        // Mapped memory regions of the process
	OomAdj      pid.OomAdj      // OomAdj factor for altering the kernel's badness heuristic
	OomScore    pid.OomScore    // Badness score of the process for OOM selection
	OomScoreAdj pid.OomScoreAdj // New style adjustment factor for altering the kernel's badness heuristic
	Root        pid.Root        // Filesystem root of a process
	Stat        pid.Stat        // Statistics about a process
	Statm       pid.Statm       // Statistics about a process's memory usage
}

// ProcessInspector inspects an individual process
type ProcessInspector struct {
	PackagingSystem pkg.System // Queries into the underlying packaging system
}

// Inspect inspects an individual process, returning a report on sucess and nil in case of an error.
//
// Returns an error in case of assembling required information fails.
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

	if fd, err := pid.NewFd(id); err != nil {
		return nil, err
	} else {
		pr.Fd = fd
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

	if oomScore, err := pid.NewOomScore(id); err != nil {
		return nil, err
	} else {
		pr.OomScore = oomScore
	}

	if root, err := pid.NewRoot(id); err != nil {
		return nil, err
	} else {
		pr.Root = root
	}

	if stat, err := pid.NewStat(id); err != nil {
		return nil, err
	} else {
		pr.Stat = *stat
	}

	if statm, err := pid.NewStatm(id); err != nil {
		return nil, err
	} else {
		pr.Statm = *statm
	}

	if bundles, err := self.PackagingSystem.Resolve(string(pr.Exe)); err != nil {
		return nil, err
	} else if len(bundles) > 0 {
		pr.Bundle = bundles[0]
	}

	return &pr, nil
}
