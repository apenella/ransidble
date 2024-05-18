package main

import (
	"fmt"
	"os"

	"github.com/apenella/ransidble/internal/configuration"
	"github.com/apenella/ransidble/internal/handler/cli"
)

func main() {

	config, err := configuration.LoadConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = cli.NewCommand(config).Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
