package commands

import (
	"fmt"
	"github.com/codegangsta/cli"
	sakura_cli "github.com/yamamoto-febc/docker-machine-sakuracloud/lib/cli"
	"strings"
)

const (
	setUsageFormat    = "Set parameter for docker-machine create option with SAKURA CLOUD driver. \nUsage: set %s %s"
	setUsageParamName = "[parameter name]"
	setUsageValue     = "[value]"
)

var setCommand = cli.Command{
	Name:  "set",
	Usage: fmt.Sprintf(setUsageFormat, setUsageParamName, setUsageValue),
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

		sakuraAPIClient := client.GetClient()
		if cnt == 1 {
			msg := c.Command.Usage
			// has extra usage?
			if option.UsageStringsFunc != nil {
				usage := option.UsageStringsFunc(sakuraAPIClient)
				msg = strings.Replace(msg, setUsageParamName, optionName, 1)
				msg = strings.Replace(msg, setUsageValue, usage, 1)
			}
			fmt.Println(msg)
			return
		}

		value := c.Args()[1]
		//TODO varidation
		client.SetConfigValue(optionName, value)

	},
	BashComplete: func(c *cli.Context) {
		// This will complete if apply over 2 args (arg1:name , arg2:value)
		cnt := len(c.Args())
		if cnt > 1 {
			return
		}

		client := sakura_cli.NewClient()
		opts := client.ListOptions()

		if cnt == 0 {
			// List All Parameter Names
			for _, opt := range opts {
				fmt.Println(opt.KeyName())
			}
		} else if cnt == 1 {
			//パラメータごとに設定可能値を判定する。
			paramName := c.Args()[0]
			for _, opt := range opts {
				if opt.KeyName() == paramName && opt.CandidateFunc != nil {
					sakuraAPIClient := client.GetClient()
					if sakuraAPIClient != nil {
						values := opt.CandidateFunc(sakuraAPIClient)
						for _, v := range values {
							fmt.Println(v)
						}
					}
				}
			}

		}
	},
}
