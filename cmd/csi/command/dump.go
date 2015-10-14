package command

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/golang/snappy"
	"github.com/vosst/csi"
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
	dumpFlagCompress = cli.BoolFlag{"compress", "compress the core dump with snappy", ""}
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
	verbose := c.Bool(dumpFlagVerbose.Name)
	compress := c.Bool(dumpFlagCompress.Name)
	cd := c.String(dumpFlagCrashDir.Name)
	pid := c.Int(dumpFlagPid.Name)
	when := time.Unix(int64(c.Int(dumpFlagTime.Name)), 0)
	sig := c.Int(dumpFlagSig.Name)
	exe := c.String(dumpFlagExe.Name)

	// Make sure that we are using a sensible umask.
	oldMask := syscall.Umask(0000)
	if verbose {
		dumpLogger.Printf("Old umask %#o", oldMask)
	}
	defer syscall.Umask(oldMask)

	cd = filepath.Join(cd, exe, when.Format(time.RFC3339), fmt.Sprint(pid))
	if err := os.MkdirAll(cd, 0755); err != nil {
		fmt.Fprintf(c.App.Writer, "Failed to create crash directory %s [%s]\n", cd, err)
	}

	ci := csi.CrashInspector{}
	if cr, err := ci.Inspect(pid, syscall.Signal(sig)); err != nil {
		fmt.Fprintf(c.App.Writer, "Failed to gather crash meta data [%s]\n", err)
	} else {
		if b, err := yaml.Marshal(cr); err != nil {
			fmt.Fprintf(c.App.Writer, "Failed to write crash meta data [%s]\n", err)
		} else {
			ry := filepath.Join(cd, "report.yaml")
			if f, err := os.Create(ry); err != nil {
				fmt.Fprintf(c.App.Writer, "Failed to dump crash report to %s [%s]\n", ry, err)
			} else {
				defer f.Close()
				fmt.Fprintf(f, "%s", b)
			}
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
		var dest io.Writer = f
		if compress {
			dest = snappy.NewWriter(f)
		}

		start := time.Now()
		n, _ := io.Copy(dest, os.Stdin)
		elapsed := time.Since(start)

		if verbose {
			dumpLogger.Printf("Wrote %d bytes of core dump to %s in %f seconds", n, df, elapsed.Seconds())
		}
	}
}

var Dump = cli.Command{
	Name:        "dump",
	Usage:       "dumps information about a crashed process",
	Description: `Usually used as the default core dump handler. Install in your system with 'csi install' (requires elevated privileges).`,
	Action:      actionDump,
	Flags:       []cli.Flag{dumpFlagVerbose, dumpFlagCrashDir, dumpFlagPid, dumpFlagUid, dumpFlagGid, dumpFlagSig, dumpFlagTime, dumpFlagHost, dumpFlagExe, dumpFlagSize},
}
