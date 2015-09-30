package cmd

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/vosst/csi/crash"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

var (
	uploadFlagDest     = cli.StringFlag{"dest", "https://daisy.ubuntu.com", "the upload destination", ""}
	uploadFlagCrashDir = cli.StringFlag{"crash-dir", "/var/crash", "directory containing crash files", ""}
	uploadFlagCleanup  = cli.BoolFlag{"cleanup", "deletes crash reports after successful upload", ""}
)

func actionUpload(c *cli.Context) {
	u, err := url.Parse(c.String(uploadFlagDest.Name))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Destination needs to be a valid url")
		return
	}

	crashDir := c.String(uploadFlagCrashDir.Name)
	entries, err := ioutil.ReadDir(crashDir)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to list files in %s (%s)\n", crashDir, err)
		return
	}

	fmt.Fprintf(os.Stdout, "Uploading crash reports from %s to %s:\n", crashDir, u.String())

	persister := crash.HttpReportPersister{*u, "0.0.1", &http.Client{}}

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

		err = persister.Persist(report)

		if err != nil {
			fmt.Fprintf(out, "  %s %s: Failed to upload crash report - %s\n", bullet, entry.Name(), err)
			continue
		}

		if c.Bool(uploadFlagCleanup.Name) {
			f.Close()
			os.Remove(fn)
		}

		fmt.Fprintf(out, "  %s %s: Successfully uploaded\n", bullet, entry.Name())
	}
}

// Command upload uploads crash reports to the server infrastructure
var Upload = cli.Command{
	Name:   "upload",
	Usage:  "uploads crash reports to the server infrastructure",
	Action: actionUpload,
	Flags:  []cli.Flag{uploadFlagDest, uploadFlagCrashDir, uploadFlagCleanup},
}
