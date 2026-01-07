package main

import (
	"github.com/stainton/service-center/cmd/app/cli"
	cmd "github.com/stainton/service-center/cmd/app/server"
)

func main() {
	command := cmd.NewCommand()
	cli.Run(command)
}
