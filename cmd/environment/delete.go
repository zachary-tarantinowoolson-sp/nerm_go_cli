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
	"github.com/spf13/viper"
)

func newDeleteCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "delete",
		Short:   "Delete a specific environment",
		Long:    "Delete environment details from the config.yaml file",
		Example: "nerm env delete | nerm env create tenant_name",
		Aliases: []string{"d"},
		RunE: func(cmd *cobra.Command, args []string) error {
			environments := utilities.FindEnvironments()

			if len(args) == 0 {
				current_environment := configs.GetCurrentEnvironment()

				if current_environment != "" {
					log.Print("You are about to delete the environment " + current_environment)
					var s string
					r := bufio.NewReader(os.Stdin)
					for {
						fmt.Fprint(os.Stderr, "Please press enter to continue:")
						s, _ = r.ReadString('\n')
						if s != "" {
							break
						}
					}
					configs.SetCurrentEnvironment("")
					delete(environments, current_environment)
					viper.Set("ALL_ENVIRONMENTS", environments)

				}
			} else {
				for _, tenant := range args {
					if environments[tenant] != nil {

						log.Print("You are about to delete the environment ", tenant)
						var s string
						r := bufio.NewReader(os.Stdin)
						for {
							fmt.Fprint(os.Stderr, "Please press enter to continue:")
							s, _ = r.ReadString('\n')
							if s != "" {
								break
							}
						}
						configs.SetCurrentEnvironment("")
						delete(environments, tenant)
						viper.Set("ALL_ENVIRONMENTS", environments)

						return nil
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
