package command

import (
	"github.com/codegangsta/cli"
	"github.com/vosst/csi/cmd/csi/command/inspect"
)

// Inspect collects information about the entire system or a specific subsystem.
var Inspect = cli.Command{
	Name:        "inspect",
	Usage:       "inspects and collects information about the system and specific characteristics",
	Subcommands: []cli.Command{inspect.Process, inspect.System},
}
