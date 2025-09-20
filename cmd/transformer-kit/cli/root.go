package cli

import (
	"github.com/complytime/baseline-demo/cmd/cli"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	command := &cobra.Command{
		Use:   "transform",
		Short: "transform CLI",
	}
	command.AddCommand(cli.NewComponentCommand(), cli.NewPlanCommand())
	return command
}
