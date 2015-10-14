package command

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/golang/snappy"
	"github.com/vosst/csi"
	"github.com/vosst/csi/pkg/debian"
	"gopkg.in/yaml.v2"
	"io"
	"log"
	"log/syslog"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

var (
	dumpLogger, _ = syslog.NewLogger(syslog.LOG_INFO|syslog.LOG_USER, log.LstdFlags)

	dumpFlagVerbose  = cli.BoolFlag{"verbose", "enables verbose output to syslog for debugging purposes", ""}
	dumpFlagCrashDir = cli.StringFlag{"crash-dir", "/var/crash", "destination directory for crash reports", ""}
	dumpFlagPid      = cli.IntFlag{"pid", -1, "pid of the crashed process", ""}
	dumpFlagUid      = cli.IntFlag{"uid", -1, "real UID of dumped process", ""}
	dumpFlagGid      = cli.IntFlag{"gid", -1, "real GID of dumped process", ""}
	dumpFlagSig      = cli.IntFlag{"sig", -1, "number of signal causing dump", ""}
	dumpFlagTime     = cli.IntFlag{"time", 0, "time of dump in seconds since epoch", ""}
	dumpFlagHost     = cli.StringFlag{"host", "", "hostname (same as nodename returned by uname)", ""}
	dumpFlagExe      = cli.StringFlag{"exe", "", "executable filename (without path prefix)", ""}
	dumpFlagSize     = cli.StringFlag{"size", "", "core file size soft resource limit", ""}
)

func actionDump(c *cli.Context) {
	// Make sure that we are using a predicable umask
	oldMask := syscall.Umask(0000)
	defer syscall.Umask(oldMask)

	verbose := c.Bool(dumpFlagVerbose.Name)
	cd := c.String(dumpFlagCrashDir.Name)
	pid := c.Int(dumpFlagPid.Name)
	// uid := c.Int(dumpFlagUid.Name)
	// gid := c.Int(dumpFlagGid.Name)
	// sig := c.Int(dumpFlagSig.Name)
	when := time.Unix(int64(c.Int(dumpFlagTime.Name)), 0)
	// size := c.String(dumpFlagSize.Name)
	// host := c.String(dumpFlagHost.Name)
	exe := c.String(dumpFlagExe.Name)

	cd = filepath.Join(cd, exe, when.Format(time.RFC3339), fmt.Sprint(pid))
	if err := os.MkdirAll(cd, 0777); err != nil {
		fmt.Fprintf(c.App.Writer, "Failed to create crash directory %s [%s]\n", cd, err)
	}

	// We start over by dumping information about the system to system.yaml
	si := csi.SystemInspector{debian.NewSystem()}
	if sysInfo, err := si.Inspect(); err != nil {
		fmt.Fprintf(c.App.Writer, "Failed to gather system information [%s]\n", err)
	} else {
		if b, err := yaml.Marshal(sysInfo); err != nil {
			fmt.Fprintf(c.App.Writer, "Failed to write system information [%s]\n", err)
		} else {
			sy := filepath.Join(cd, "system.yaml")
			if f, err := os.Create(sy); err != nil {
				fmt.Fprintf(c.App.Writer, "Failed to dump system information to %s [%s]\n", sy, err)
			} else {
				defer f.Close()
				fmt.Fprintf(f, "%s", b)
			}
		}
	}

	// Next we dump information about the crashed process
	pi := csi.ProcessInspector{debian.NewSystem()}
	processInfo, _ := pi.Inspect(pid)

	if b, err := yaml.Marshal(processInfo); err != nil {
		fmt.Fprintf(c.App.Writer, "Failed to query process information [%s] \n", err)
	} else {
		py := filepath.Join(cd, "process.yaml")
		if f, err := os.Create(py); err != nil {
			fmt.Fprintf(c.App.Writer, "Failed to dump process information to %s [%s]\n", py, err)
		} else {
			defer f.Close()
			fmt.Fprintf(f, "%s", b)
		}
	}

	// And we finally dump the actual core file
	df := filepath.Join(cd, "core")
	if f, err := os.Create(df); err != nil {
		fmt.Fprintf(c.App.Writer, "Failed to dump core to %s [%s]", df, err)
	} else {
		defer f.Close()
		// TODO(tvoss): Investigate into syscall.Sendfile and figure out a way
		// to avoid copying of data to userspace.
		start := time.Now()
		n, _ := io.Copy(snappy.NewWriter(f), os.Stdin)
		elapsed := time.Since(start)

		if verbose {
			dumpLogger.Printf("Wrote %d bytes of compressed core to %s in %f seconds", n, df, elapsed.Seconds())
		}
	}
}

var Dump = cli.Command{
	Name:   "dump",
	Usage:  "dumps information about a crashed process",
	Action: actionDump,
	Flags:  []cli.Flag{dumpFlagVerbose, dumpFlagCrashDir, dumpFlagPid, dumpFlagUid, dumpFlagGid, dumpFlagSig, dumpFlagTime, dumpFlagHost, dumpFlagExe, dumpFlagSize},
}
