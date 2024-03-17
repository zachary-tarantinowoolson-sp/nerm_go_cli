/*
Copyright © 2024 Zachary Tarantino-Woolson <zachary.tarantino@sailpoint.com>
*/
package environment

import (
	"log"
	"nerm/cmd/configs"

	"github.com/spf13/cobra"
)

func newUseCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "use",
		Short:   "Select a tenant to be used as the current environment",
		Long:    "Select a tenant to be used as the current environment in the CLI",
		Example: "nerm env use tenant_name",
		Aliases: []string{"u"},
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, environmentName := range args {
				configs.SetCurrentEnvironment(environmentName)
			}

			if len(args) == 0 {
				log.Print("Please specify a tenant_name to use")
			}

			return nil
		},
	}
}
