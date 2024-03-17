/*
Copyright © 2024 Zachary Tarantino-Woolson <zachary.tarantino@sailpoint.com>
*/
package health_check

import (
	"fmt"
	"io/ioutil"
	"log"
	"nerm/cmd/configs"
	"net/http"

	"github.com/spf13/cobra"
)

func NewHealthCheckCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "health_check",
		Short:   "Pings an environment to perform a health check",
		Long:    "Pings an environment to perform a health check",
		Example: "nerm hc | nerm hc env_name",
		Aliases: []string{"hc"},
		RunE: func(cmd *cobra.Command, args []string) error {
			baseurl := configs.GetBaseURL()

			if len(args) == 0 {
				tenant := configs.GetTenant()

				response, err := http.Get("http://" + tenant + "." + baseurl + "/health_check")
				if err != nil {
					log.Fatal(err)
				}
				data, err := ioutil.ReadAll(response.Body)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println(string(data))
			} else {
				for _, environmentName := range args {
					old_env := configs.GetCurrentEnvironment()
					configs.SetCurrentEnvironment(environmentName)
					tenant := configs.GetTenant()

					response, err := http.Get("http://" + tenant + "." + baseurl + "/health_check")
					if err != nil {
						log.Fatal(err)
					}
					data, err := ioutil.ReadAll(response.Body)
					if err != nil {
						log.Fatal(err)
					}
					fmt.Println(string(data))
					configs.SetCurrentEnvironment(old_env)
				}

			}

			return nil
		},
	}
}
