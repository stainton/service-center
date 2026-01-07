package main

import (
	"github.com/stainton/service-center/cmd/app/cli"
	"github.com/stainton/service-center/cmd/app/client"
)

func main() {
	command := client.NewCommand()
	cli.Run(command)
}
