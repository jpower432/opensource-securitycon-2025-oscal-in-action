package main

import (
	"os"

	"github.com/jpower432/opensource-securitycon-2025-oscal-in-action/cmd/transformer-kit/cli"
)

func main() {
	command := cli.New()
	err := command.Execute()
	if err != nil {
		os.Exit(1)
	}
}
