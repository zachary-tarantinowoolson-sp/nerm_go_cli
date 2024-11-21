/*
Copyright © 2024 Zachary Tarantino-Woolson <zachary.tarantino@sailpoint.com>
*/
package environment

import (
	"nerm/cmd/configs"
	"sort"

	"golang.org/x/exp/maps"

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
			// environments := utilities.FindEnvironments()

			envs := maps.Keys(configs.GetAllEnvironmentStrings())
			sort.Strings(envs)

			// for k := range environments {
			// 	fmt.Println(k)
			// }

			printEnvListTable(envs)

			return nil
		},
	}
}
