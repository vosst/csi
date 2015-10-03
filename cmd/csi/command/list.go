package command

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/vosst/csi/crash"
	"io"
	"os"
)

const bullet = "\u2022"

var listFlagCrashDir = cli.StringFlag{"crash-dir", "/var/crash", "directory containing crash files", ""}

// ListingVisitor provides listing of available crash reports.
type ListingVisitor struct {
	Out io.Writer // Destination for output.
}

func (self ListingVisitor) NewReport(name string, report crash.Report) {
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

	fmt.Fprintf(self.Out, "  %s %s[%s]: %s - %s \n", bullet, name, problemType, executablePath, annotation)

}

func (self ListingVisitor) NewError(err error) {
	// We silently skip errors.
}

func actionList(c *cli.Context) {
	crashDir := c.String(listFlagCrashDir.Name)
	fmt.Fprintf(os.Stdout, "Listing crash reports in %s:\n", crashDir)
	crash.ForEachReportInDir(crashDir, ListingVisitor{os.Stdout})
}

var List = cli.Command{
	Name:   "list",
	Usage:  "lists all crash reports on the system",
	Flags:  []cli.Flag{listFlagCrashDir},
	Action: actionList,
}
