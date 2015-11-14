package main

import (
	"github.com/docker/machine/libmachine/drivers/plugin"
	"github.com/yamamoto-febc/docker-machine-sakuracloud"
)

func main() {
	plugin.RegisterDriver(sakuracloud.NewDriver("", ""))
}
