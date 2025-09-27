package cli

import (
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	command := &cobra.Command{
		Use:   "transform",
		Short: "transform CLI",
	}
	command.AddCommand(NewPlanCommand())
	return command
}
