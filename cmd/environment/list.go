/*
Copyright © 2024 Zachary Tarantino-Woolson <zachary.tarantino@sailpoint.com>
*/
package environment

import (
	"fmt"

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
			for _, env := range environments {
				fmt.Println(env)
			}

			// if len(environments) != 0 {
			// 	log.Warn("You are about to Print out the list of Environments")
			// 	res := terminal.InputPrompt("Press Enter to continue")
			// 	log.Info("Response", "res", res)
			// 	if res == "" {
			// 		fmt.Println(util.PrettyPrint(environments))
			// 	}
			// } else {
			// 	log.Warn("No environments configured")
			// 	return nil
			// }
			return nil
		},
	}
}
