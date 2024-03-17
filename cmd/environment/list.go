/*
Copyright © 2024 Zachary Tarantino-Woolson <zachary.tarantino@sailpoint.com>
*/
package environment

import (
	"fmt"
	"nerm/cmd/utilities"

	"github.com/spf13/cobra"
)

func newListCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "List all environments",
		Long:    "List all environments currently configured for use in the CLI",
		Example: "nerm env list",
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			environments := utilities.FindEnvironments()
			// for _, env := range environments {
			// }

			for k := range environments {
				fmt.Println(k)
			}

			return nil
		},
	}
}
