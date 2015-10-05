package inspect

import (
	"encoding/json"
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/vosst/csi"
)

func actionSystem(context *cli.Context) {
	si := csi.SystemInspector{}
	sysInfo, _ := si.Inspect()

	if b, err := json.MarshalIndent(sysInfo, "", "  "); err != nil {
		fmt.Fprintf(context.App.Writer, "Failed to query system information")
	} else {
		fmt.Fprintf(context.App.Writer, "%s\n", b)
	}

}

// System collects system/OS-specific information.
var System = cli.Command{
	Name:   "system",
	Usage:  "collects system/OS-specific information",
	Action: actionSystem,
}
