package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"os"
)

type Version struct {
	major int
	minor int
	patch int
}

func (self Version) String() string {
	return fmt.Sprintf("%d.%d.%d", self.major, self.minor, self.patch)
}

var CurrentVersion = Version{0, 0, 1}

func main() {
	app := cli.NewApp()
	app.Name = "csi"
	app.Usage = "monitors the system and helps to investigate issues."
	app.Version = CurrentVersion.String()
	app.Authors = []cli.Author{
		cli.Author{"Thomas Voß", "thomas.voss@canonical.com"},
		cli.Author{"Evan Dandrea", "evan.dandrea@canonical.com"},
	}

	app.Commands = []cli.Command{
		CommandId,
		CommandList,
		CommandUpload,
	}

	app.Run(os.Args)
}
