package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
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

	persister := HttpCrashReportPersister{*u, AlwaysReachableReachabilityMonitor{}, &http.Client{}}

	for _, entry := range entries {
		fn := fmt.Sprintf("%s/%s", crashDir, entry.Name())

		fmt.Fprintf(os.Stdout, "Processing %s ...", fn)

		f, err := os.Open(fn)
		if err != nil {
			fmt.Fprintf(os.Stderr, "  Failed to open crash report.")
			continue
		}

		defer f.Close()
		report, err := ParseCrashReport(NewLineReader{f})

		if err != nil {
			fmt.Fprintf(os.Stderr, "  Failed to parse crash report.")
			continue
		}

		err = persister.Persist(report)

		if err != nil {
			fmt.Fprintf(os.Stderr, "  Failed to upload crash report.")
			continue
		}

		if c.Bool(uploadFlagCleanup.Name) {
			f.Close()
			os.Remove(fn)
		}

		fmt.Fprintf(os.Stdout, " [\u2713]\n")
	}
}

// CommandUpload uploads crash reports to the server infrastructure
var CommandUpload = cli.Command{
	Name:   "upload",
	Usage:  "uploads crash reports to the server infrastructure",
	Action: actionUpload,
	Flags:  []cli.Flag{uploadFlagDest, uploadFlagCrashDir, uploadFlagCleanup},
}
