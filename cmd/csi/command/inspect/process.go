package inspect

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/vosst/csi"
	"github.com/vosst/csi/pkg/debian"
	"gopkg.in/yaml.v2"
	"os"
	"strconv"
)

var (
	processFlagPid = cli.StringFlag{"pid", "", "specify the pid of the process that should be inspected", ""}
)

func actionProcess(context *cli.Context) {
	pid := os.Getpid()

	if p := context.String(processFlagPid.Name); len(p) > 0 {
		pid, _ = strconv.Atoi(p)
	}

	pi := csi.ProcessInspector{debian.NewSystem()}
	processInfo, _ := pi.Inspect(pid)

	if b, err := yaml.Marshal(processInfo); err != nil {
		fmt.Fprintf(context.App.Writer, "Failed to query process information")
	} else {
		fmt.Fprintf(context.App.Writer, "%s\n", b)
	}
}

// Process collects process-specific information.
var Process = cli.Command{
	Name:   "process",
	Usage:  "collects process-specific information",
	Flags:  []cli.Flag{processFlagPid},
	Action: actionProcess,
}
