package main

import (
	"os"
	"path"

	"github.com/codegangsta/cli"
	"github.com/docker/machine/libmachine/drivers/plugin"
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
`

func main() {

	cli.AppHelpTemplate = appHelpTemplate
	app := cli.NewApp()
	app.Name = path.Base(os.Args[0])
	app.Usage = "This is a Docker Machine plugin binary. Please use it through the main 'docker-machine' binary."
	app.Author = "Kazumichi Yamamoto(yamamoto.febc@gmail.com)"
	app.Email = "yamamoto.febc@gmail.com"
	app.Version = version.FullVersion()
	app.Action = func(c *cli.Context) {
		plugin.RegisterDriver(driver.NewDriver("", ""))
	}
	app.Run(os.Args)
}
