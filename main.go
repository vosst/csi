package main

import (
	"github.com/codegangsta/cli"
	"github.com/vosst/csi/cmd"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "csi"
	app.Usage = "monitors the system and helps to investigate issues."
	app.Version = "0.0.1"
	app.Authors = []cli.Author{
		cli.Author{"Thomas Voß", "thomas.voss@canonical.com"},
		cli.Author{"Evan Dandrea", "evan.dandrea@canonical.com"},
	}

	app.Commands = []cli.Command{
		cmd.Id,
		cmd.List,
		cmd.Upload,
	}

	app.Run(os.Args)
}
