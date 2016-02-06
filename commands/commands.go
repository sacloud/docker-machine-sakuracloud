package commands

import (
	"github.com/codegangsta/cli"
)

// Commands cli commands
var Commands = []cli.Command{
	listCommand,
	setCommand,
	clearCommand,
	infoCommand,
}

// Flags cli flags
var Flags = []cli.Flag{}
