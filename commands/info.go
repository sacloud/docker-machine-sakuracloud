package commands

import (
	"fmt"
	"github.com/codegangsta/cli"
	sakura_cli "github.com/yamamoto-febc/docker-machine-sakuracloud/lib/cli"
	"os"
	"text/template"
)

const (
	infoUsageFormat    = "Infomation of parameter for docker-machine create option with SAKURA CLOUD driver. \nUsage: info [parameter name] %s"
	infoUsageParamName = ""
	infoTemplate       = `
===== Config for SAKURA CLOUD =====

Name : {{.option.KeyName}}
Description : {{.option.Description}}
Current : {{.config.FormatedCurrentValue}}

`
)

var infoCommand = cli.Command{
	Name:  "info",
	Usage: fmt.Sprintf(infoUsageFormat, infoUsageParamName),
	Action: func(c *cli.Context) {
		cnt := len(c.Args())
		client := sakura_cli.NewClient()

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

		config, err := client.GetConfigValue(optionName)
		if err != nil {
			fmt.Printf("error:%v", err)
			os.Exit(1)
		}

		t := template.New("infoTemplate")
		tmpl, err := t.Parse(infoTemplate)
		if err != nil {
			fmt.Printf("error:%v", err)
			os.Exit(1)
		}

		configValue := map[string]interface{}{
			"option": option,
			"config": config,
		}

		tmpl.Execute(os.Stdout, configValue)
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
