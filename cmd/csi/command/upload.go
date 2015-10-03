package command

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/vosst/csi/crash"
	"github.com/vosst/csi/machine"
	"io"
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

type UploadingVisitor struct {
	CrashDir  string                // The directory containing crashes
	Out       io.Writer             // Destination for output
	Persister crash.ReportPersister // Persister provides persistence of crash reports.
	Cleanup   bool                  // If true, successfully uploaded crash reports are deleted.
}

func (self UploadingVisitor) NewReport(name string, report crash.Report) {
	err := self.Persister.Persist(report)

	if err != nil {
		fmt.Fprintf(self.Out, "  %s %s: Failed to upload crash report - %s\n", bullet, name, err)
		return
	}

	fmt.Fprintf(self.Out, "  %s %s: Successfully uploaded\n", bullet, name)

	if self.Cleanup {
		os.Remove(filepath.Join(self.CrashDir, name))
	}
}

func (self UploadingVisitor) NewError(err error) {
}

func actionUpload(c *cli.Context) {
	mi, err := machine.DefaultIdentifier()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not create default machine identifier")
		return
	}

	u, err := url.Parse(c.String(uploadFlagDest.Name))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Destination needs to be a valid url")
		return
	}

	crashDir := c.String(uploadFlagCrashDir.Name)
	persister := crash.HttpReportPersister{*u, mi, &http.Client{}}

	crash.ForEachReportInDir(crashDir, UploadingVisitor{crashDir, os.Stdout, persister, c.Bool(uploadFlagCleanup.Name)})
}

// Command upload uploads crash reports to the server infrastructure
var Upload = cli.Command{
	Name:   "upload",
	Usage:  "uploads crash reports to the server infrastructure",
	Action: actionUpload,
	Flags:  []cli.Flag{uploadFlagDest, uploadFlagCrashDir, uploadFlagCleanup},
}
