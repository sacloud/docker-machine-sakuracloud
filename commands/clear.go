package commands

import (
	"fmt"
	"github.com/codegangsta/cli"
	// TODO use shell detect "github.com/docker/machine/libmachine/shell"
	sakura_cli "github.com/yamamoto-febc/docker-machine-sakuracloud/lib/cli"
	"os"
)

const (
	clearUsageFormat    = "Clear parameter for docker-machine create option with SAKURA CLOUD driver. \nUsage: clear [parameter name] %s"
	clearUsageParamName = ""
	clearEnvFormat      = `unset %s
# Run this command to configure your shell:
# eval $(docker-machine-driver-sakuracloud clear %s)
`
)

var clearCommand = cli.Command{
	Name:  "clear",
	Usage: fmt.Sprintf(clearUsageFormat, clearUsageParamName),
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "all, a",
			Usage: "clear all config",
		},
	},
	Action: func(c *cli.Context) {
		cnt := len(c.Args())
		client := sakura_cli.NewClient()
		targetNames := []string{}

		if c.Bool("all") {
			opts := client.ListOptions()
			for _, opt := range opts {
				targetNames = append(targetNames, opt.KeyName())
			}
		} else {
			if cnt == 0 {
				fmt.Println(c.Command.Usage)
				return
			}

			optionName := c.Args()[0]
			option := client.GetOption(optionName)
			if option == nil {
				fmt.Println(c.Command.Usage)
				return
			}
			targetNames = append(targetNames, optionName)
		}

		for _, name := range targetNames {
			config, err := client.GetConfigValue(name)
			if err != nil {
				fmt.Printf("error:%v", err)
				os.Exit(1)
			}

			if config.IsFromEnv() {
				fmt.Printf(clearEnvFormat, config.EnvName, config.KeyName)
			}

			client.ClearConfigValue(name)
		}

	},
	BashComplete: func(c *cli.Context) {
		// This will complete if apply over 2 args (arg1:name , arg2:value)
		cnt := len(c.Args())
		if cnt > 0 {
			return
		}

		client := sakura_cli.NewClient()
		opts := client.ListOptions()

		// List All Parameter Names
		for _, opt := range opts {
			fmt.Println(opt.KeyName())
		}

	},
}
