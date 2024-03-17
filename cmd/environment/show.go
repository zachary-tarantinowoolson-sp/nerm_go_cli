/*
Copyright © 2024 Zachary Tarantino-Woolson <zachary.tarantino@sailpoint.com>
*/
package environment

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"nerm/cmd/configs"
	"nerm/cmd/utilities"
	"os"

	"github.com/spf13/cobra"
)

func newShowCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "show",
		Short:   "Show the details of an environment",
		Long:    "Show the details of the current or a specified environment",
		Example: "nerm env show | nerm env show tenant-name",
		Aliases: []string{"s"},
		RunE: func(cmd *cobra.Command, args []string) error {
			environments := utilities.FindEnvironments()

			if len(args) == 0 {
				current_environment := configs.GetCurrentEnvironment()

				if current_environment != "" {
					log.Print("You are about to print the detials of the environment details of " + current_environment)
					var s string
					r := bufio.NewReader(os.Stdin)
					for {
						fmt.Fprint(os.Stderr, "Please press enter to continue:")
						s, _ = r.ReadString('\n')
						if s != "" {
							break
						}
					}
					formatted, err := json.MarshalIndent(environments[current_environment], "", "  ")
					if err != nil {
						log.Fatal(err)
					}
					fmt.Println(string(formatted))
					return nil
				}
			} else {
				for _, tenant := range args {
					if environments[tenant] != nil {

						log.Print("You are about to Print out the environment details of ", "env", tenant)
						var s string
						r := bufio.NewReader(os.Stdin)
						for {
							fmt.Fprint(os.Stderr, "Please press enter to continue:")
							s, _ = r.ReadString('\n')
							if s != "" {
								break
							}
						}
						formatted, err := json.MarshalIndent(environments[tenant], "", "  ")
						if err != nil {
							log.Fatal(err)
						}
						fmt.Println(string(formatted))

						return nil
					} else {
						log.Fatal("Environment does not exist\nUse `nerm env create " + tenant + "` to create it.")
						return nil
					}

				}
			}

			return nil
		},
	}
}
