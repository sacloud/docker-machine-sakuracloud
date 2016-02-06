package main

import (
	"os"
	"path"

	"github.com/codegangsta/cli"
	"github.com/docker/machine/libmachine/drivers/plugin"
	"github.com/docker/machine/libmachine/log"
	"github.com/yamamoto-febc/docker-machine-sakuracloud/commands"
	"github.com/yamamoto-febc/docker-machine-sakuracloud/driver"
	"github.com/yamamoto-febc/docker-machine-sakuracloud/version"
)

var appHelpTemplate = `This is a Docker Machine plugin for SAKURA CLOUD.
Plugin binaries are not intended to be invoked directly.
Please use this plugin through the main 'docker-machine' binary.

Version: {{.Version}}{{if or .Author .Email}}

Author:{{if .Author}}
  {{.Author}}{{if .Email}} - <{{.Email}}>{{end}}{{else}}
  {{.Email}}{{end}}{{end}}
{{if .Flags}}
Options:
  {{range .Flags}}{{.}}
  {{end}}{{end}}
Commands:
  {{range .Commands}}{{.Name}}{{with .ShortName}}, {{.}}{{end}}{{ "\t" }}{{.Usage}}
  {{end}}
`

func main() {

	cli.AppHelpTemplate = appHelpTemplate
	app := cli.NewApp()
	app.Name = path.Base(os.Args[0])
	app.Usage = "This is a Docker Machine plugin binary. Please use it through the main 'docker-machine' binary."
	app.Author = "Kazumichi Yamamoto(yamamoto.febc@gmail.com)"
	app.Email = "https://github.com/yamamoto-febc/docker-machine-sakuracloud/"
	app.Version = version.FullVersion()
	app.Action = func(c *cli.Context) {
		plugin.RegisterDriver(driver.NewDriver("", ""))
	}
	app.CommandNotFound = cmdNotFound
	app.Commands = commands.Commands
	app.Flags = commands.Flags
	app.EnableBashCompletion = true
	app.Run(os.Args)
}

func cmdNotFound(c *cli.Context, command string) {
	log.Errorf(
		"%s: '%s' is not a %s command. See '%s --help'.",
		c.App.Name,
		command,
		c.App.Name,
		os.Args[0],
	)
	os.Exit(1)
}
