package cli

import (
	"github.com/apenella/ransidble/internal/configuration"
	"github.com/apenella/ransidble/internal/handler/cli/serve"
	"github.com/spf13/cobra"
)

func NewCommand(config *configuration.Configuration) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ransidble",
		Short: "Ransidble is a tool to execute Ansible commands on a remote host",
		Long:  "Ransidble is a tool to execute Ansible commands on a remote host",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := cmd.Help()
			if err != nil {
				return err
			}
			return nil
		},
	}

	cmd.AddCommand(serve.NewCommand(config))

	return cmd
}
