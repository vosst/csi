package csi

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"syscall"

	"github.com/vosst/csi/log"
	"github.com/vosst/csi/pkg"
)

// Poor man's version of StatFs, just exposing the values we are actually interested in
type FSStats struct {
	BlockSize      int64  // Optimal transfer block size.
	BlockCount     uint64 // Total number of data blocks in a file system.
	BlockFree      uint64 // Free blocks in a file system.
	BlockAvailable uint64 // Free blocks available to unprivileged users.

}

// Mount describes a mounted filesytem. Please see man fstab for further details.
type Mount struct {
	Spec          string   // The field describes the block special device or remote filesystem to be mounted.
	File          string   // Describes the mount point for the filesytem.
	Type          string   // Describes the type of the filesystem.
	MntOps        string   // Describes the mount options associated with the filesystem.
	DumpFrequency int      // Dump frequency in days.
	PassNo        int      // Pass number on parallel fsck.
	FSStats       *FSStats // Filesystem data, may be nil.
}

// ParseMounts reads all mounted file systems from reader, expecting a line format
// as specified in man fstab. For every mounted filesystem, ParseMounts tries to
// query size information.
// The function is quite robust and tries to keep on processing for as long as possible,
// reporting partial results together with errors.
func parseMounts(reader io.Reader) []Mount {
	mounts := []Mount{}

	br := bufio.NewReader(reader)
	for s, err := br.ReadString('\n'); err == nil; s, err = br.ReadString('\n') {
		mnt := Mount{}
		if _, err := fmt.Sscanf(s, "%s %s %s %s %d %d", &mnt.Spec, &mnt.File, &mnt.Type, &mnt.MntOps, &mnt.DumpFrequency, &mnt.PassNo); err != nil {
			continue
		}

		statfs := syscall.Statfs_t{}
		if err = syscall.Statfs(mnt.File, &statfs); err == nil {
			fsStats := FSStats{statfs.Bsize, statfs.Blocks, statfs.Bfree, statfs.Bavail}
			mnt.FSStats = &fsStats
		}
		mounts = append(mounts, mnt)
	}

	return mounts
}

// OSReport summarizes information about the operating system
type OSReport struct {
	Name    string   // Name of the OS
	Release string   // Relase of the OS
	Logs    struct { // Central logs documenting the OS operations
		Dmesg  []byte // Contents of the kernel log buffer
		Syslog []byte // Contents of syslog
	}
	Memory struct { // Information about total available/free memory
		Total uint64 // Total usable RAM
		Free  uint64 // Amound of memory currently unused
	}
	Swap struct { // Information about total available/free swap
		Total uint64 // Total amount of swap space available
		Free  uint64 // Amount of swap space that is currently unused
	}
	Mounts []Mount // All mounted filesystems
}

// OSInspector provides means to gather information about the operating system
type OSInspector struct {
	DmesgCollector  log.Collector
	SyslogCollector log.Collector
	ReleaseFile     string
	MemInfo         string
	MTab            string
}

func (self OSInspector) Inspect() (OSReport, error) {
	osi := OSReport{}
	{
		f, err := os.Open(self.ReleaseFile)
		if err != nil {
			return osi, errors.New(fmt.Sprintf("Failed to open release file %s [%s]", self.ReleaseFile, err))
		}

		defer f.Close()

		br := bufio.NewReader(f)
		for s, err := br.ReadString('\n'); err == nil; s, err = br.ReadString('\n') {
			tokens := strings.Split(s, "=")

			switch tokens[0] {
			case "DISTRIB_ID":
				osi.Name = strings.Trim(tokens[1], "\t \n")
			case "DISTRIB_RELEASE":
				osi.Release = strings.Trim(tokens[1], "\t \n")
			}
		}
	}

	{
		f, err := os.Open(self.MemInfo)
		if err != nil {
			return osi, errors.New(fmt.Sprintf("Failed to open meminfo file %s [%s]", self.MemInfo, err))
		}

		defer f.Close()

		br := bufio.NewReader(f)
		for s, err := br.ReadString('\n'); err == nil; s, err = br.ReadString('\n') {
			tokens := strings.Split(s, ":")
			switch tokens[0] {
			case "MemTotal":
				fmt.Sscanf(strings.Trim(tokens[1], "\t \n"), "%d kB", &(osi.Memory.Total))
			case "MemFree":
				fmt.Sscanf(strings.Trim(tokens[1], "\t \n"), "%d kB", &osi.Memory.Free)
			case "SwapTotal":
				fmt.Sscanf(strings.Trim(tokens[1], "\t \n"), "%d kB", &osi.Swap.Total)
			case "SwapFree":
				fmt.Sscanf(strings.Trim(tokens[1], "\t \n"), "%d kB", &osi.Swap.Free)
			}
		}
	}

	{
		f, err := os.Open(self.MTab)
		if err != nil {
			return osi, errors.New(fmt.Sprintf("Failed to open mtab file %s [%s]", self.MTab, err))
		}

		defer f.Close()

		osi.Mounts = parseMounts(f)

	}

	if b, err := self.DmesgCollector.Collect(); err == nil {
		osi.Logs.Dmesg = b
	}

	if b, err := self.SyslogCollector.Collect(); err == nil {
		osi.Logs.Syslog = b
	}

	return osi, nil
}

// SystemReport bundles system-specific information relevant
// in reporting and tracking down issues.
type SystemReport struct {
	HostName     string   // HostName of this machine.
	Architecture pkg.Arch // Host architecture.
	OS           OSReport // Information about the OS.
}

// SystemInspector inspects core properties of the current system.
type SystemInspector struct {
	PkgSystem pkg.System // Retrievs information from the packaging system.
}

// Inspect gathers information about the current system and encodes
// it to encoder.
func (self SystemInspector) Inspect() (si SystemReport, err error) {
	var hn string
	if hn, err = os.Hostname(); err != nil {
		return
	}

	si.HostName = hn
	si.Architecture, _ = self.PkgSystem.Arch()

	os := OSInspector{log.NewDmesgCollector(), log.NewSyslogCollector(), "/etc/lsb-release", "/proc/meminfo", "/etc/mtab"}
	si.OS, err = os.Inspect()

	if err != nil {
		err = errors.New(fmt.Sprintf("Failed to inspect the OS [%s]", err))
		return
	}

	return
}
