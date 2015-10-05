package command

import (
	"encoding/json"
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/vosst/csi"
)

func actionInspect(context *cli.Context) {
	si := csi.SystemInspector{}
	sysInfo, _ := si.Inspect()

	if b, err := json.MarshalIndent(sysInfo, "", "  "); err != nil {
		fmt.Fprintf(context.App.Writer, "Failed to query system information")
	} else {
		fmt.Fprintf(context.App.Writer, "%s\n", b)
	}
}

// Inspect collects information about the entire system or a specific subsystem.
var Inspect = cli.Command{
	Name:   "inspect",
	Usage:  "inspects and collects information about the (sub)system",
	Action: actionInspect,
}
