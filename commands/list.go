package commands

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/olekukonko/tablewriter"
	sakura_cli "github.com/yamamoto-febc/docker-machine-sakuracloud/lib/cli"
	"os"
)

const (
	listCommandHeaderText = "===== Setting List for SAKURA CLOUD ====="
)

var listCommand = cli.Command{
	Name:  "list",
	Usage: "List parameters for docker-machine create option with SAKURA CLOUD driver",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "detail, d",
			Usage: "output detail columns",
		},
	},
	Action: func(c *cli.Context) {
		client := sakura_cli.NewClient()
		options, err := client.ListConfigValue()
		if err != nil {
			fmt.Printf("Error : %v", err)
			os.Exit(1)
		}

		isDetail := c.Bool("detail")
		table := tablewriter.NewWriter(os.Stdout)
		// set table format
		if isDetail {
			table.SetRowLine(true)
			table.SetRowSeparator("-")
		}
		table.SetAlignment(tablewriter.ALIGN_LEFT)

		//set values (with header)
		table.SetHeader(sakura_cli.GetPrintHeader(isDetail))
		for _, opt := range options {
			table.Append(opt.GetPrintInfo(isDetail))
		}

		//output
		// command header
		fmt.Fprintf(os.Stdout, "\n%s\n\n", listCommandHeaderText)
		// command body
		table.Render()
	},
}
