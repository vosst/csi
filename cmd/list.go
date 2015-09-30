package cmd

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/vosst/csi/crash"
	"io/ioutil"
	"os"
	"path/filepath"
)

const bullet = "\u2022"

var listFlagCrashDir = cli.StringFlag{"crash-dir", "/var/crash", "directory containing crash files", ""}

func actionList(c *cli.Context) {
	crashDir := c.String(listFlagCrashDir.Name)
	entries, err := ioutil.ReadDir(crashDir)

	fmt.Fprintf(os.Stdout, "Listing crash reports in %s:\n", crashDir)

	if err != nil {
		fmt.Fprintf(os.Stdout, "[x]\n")
		return
	}

	for _, entry := range entries {
		fn := filepath.Join(crashDir, entry.Name())

		if matched, _ := filepath.Match("*.crash", entry.Name()); !matched {
			// Silently skip over the file.
			continue
		}

		out := os.Stdout

		f, err := os.Open(fn)
		if err != nil {
			fmt.Fprintf(out, "  %s %s: Failed to open file - %s\n", bullet, entry.Name(), err)
			continue
		}

		defer f.Close()
		report, err := crash.ParseReport(crash.NewLineReader{f})

		if err != nil {
			fmt.Fprintf(out, "  %s %s: Failed to parse crash report - %s\n", bullet, entry.Name(), err)
			continue
		}

		// fmt.Print(report)
		problemType := "unknown"
		if pt := report["Problemtype"]; pt != nil && len(pt) > 0 {
			problemType = pt[0]
		}

		executablePath := "unknown executable"
		if ep := report["Executablepath"]; ep != nil && len(ep) > 0 {
			executablePath = ep[0]
		}

		annotation := "no further details"
		if a := report["Annotation"]; a != nil && len(a) > 0 {
			annotation = a[0]
		}

		fmt.Fprintf(out, "  %s %s[%s]: %s - %s \n", bullet, entry.Name(), problemType, executablePath, annotation)
	}
}

var List = cli.Command{
	Name:   "list",
	Usage:  "lists all crash reports on the system",
	Flags:  []cli.Flag{listFlagCrashDir},
	Action: actionList,
}
