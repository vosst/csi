package cmd

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/vosst/csi/machine"
	"os"
)

func actionId(c *cli.Context) {
	mi, err := machine.DefaultIdentifier()
	if err == nil {
		id, err := mi.Identify()
		if err == nil {
			fmt.Fprintf(c.App.Writer, "%x\n", id)
		} else {
			fmt.Fprintf(os.Stderr, "Error determining machine id: %s\n", err)
		}
	} else {
		fmt.Fprintf(os.Stderr, "Error initializing identification infrastructure: %s\n", err)
	}

}

// Command id prints the sha512 hash of the machine/device id.
var Id = cli.Command{
	Name:   "id",
	Usage:  "prints the sha512 hash of the system id",
	Action: actionId,
}
