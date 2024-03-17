/*
Copyright © 2024 Zachary Tarantino-Woolson <zachary.tarantino@sailpoint.com>
*/
package environment

import (
	"bufio"
	"fmt"
	"log"
	"nerm/cmd/configs"
	"nerm/cmd/utilities"
	"os"

	"github.com/spf13/cobra"
)

func newUpdateCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "update",
		Short:   "Update an existing environment",
		Long:    "Update an existing environment's Tenant and/or Token for use in the CLI",
		Example: "nerm env update | nerm env update env_name",
		Aliases: []string{"c"},
		RunE: func(cmd *cobra.Command, args []string) error {
			environments := utilities.FindEnvironments()

			if len(args) == 0 {
				current_environment := configs.GetCurrentEnvironment()

				if current_environment != "" {
					log.Print("You are about to update the environment " + current_environment)
					var s string
					r := bufio.NewReader(os.Stdin)
					for {
						fmt.Fprint(os.Stderr, "Please press enter to continue:")
						s, _ = r.ReadString('\n')
						if s != "" {
							break
						}
					}
					utilities.UpdateEnvironment(current_environment)
				}
			} else {
				for _, tenant := range args {
					if environments[tenant] != nil {

						log.Print("You are about to update the environment ", tenant)
						var s string
						r := bufio.NewReader(os.Stdin)
						for {
							fmt.Fprint(os.Stderr, "Please press enter to continue:")
							s, _ = r.ReadString('\n')
							if s != "" {
								break
							}
						}
						utilities.UpdateEnvironment(tenant)

					} else {
						log.Fatal("Environment does not exist.")
						return nil
					}

				}
			}

			return nil
		},
	}
}
