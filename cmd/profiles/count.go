/*
Copyright © 2024 Zachary Tarantino-Woolson <zachary.tarantino@sailpoint.com>
*/
package profiles

import (
	"fmt"
	"log"
	"nerm/cmd/configs"
	"nerm/cmd/utilities"
	"net/http/httputil"
	"net/url"

	"github.com/spf13/cobra"
)

func newProfileCountCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "count",
		Short:   "Pulls a count of all Profiles in an environment",
		Long:    "Pulls a count of all Profiles in an environment by profile Type",
		Example: "nerm profiles count | nerm p count env_name",
		Aliases: []string{"c"},
		RunE: func(cmd *cobra.Command, args []string) error {

			type_params := url.Values{}
			type_params.Add("limit", "100")
			type_params.Add("metadata", "true")

			types_resp, types_err := utilities.MakeAPIRequests("get", "profile_types", "", type_params.Encode(), nil)
			if types_err != nil {
				log.Fatal(types_err)
			}
			dumped_res, dump_err := httputil.DumpResponse(types_resp, true)
			fmt.Println(string(dumped_res), dump_err)

			// bodyBytes, read_err := io.ReadAll(types_resp.Body)
			// if read_err != nil {
			// 	log.Fatal(read_err)
			// }
			// bodyString := string(bodyBytes)
			// fmt.Println(bodyString)

			if len(args) == 0 {

				params := url.Values{}
				// params.Add("profile_type_id", "123")
				params.Add("limit", "1")
				params.Add("metadata", "true")
				params.Add("exclude_attributes", "true")

				resp, err := utilities.MakeAPIRequests("get", "profiles", "", params.Encode(), nil)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println(resp)

			} else {
				for _, environmentName := range args {
					old_env := configs.GetCurrentEnvironment()
					configs.SetCurrentEnvironment(environmentName)

					configs.SetCurrentEnvironment(old_env)
				}

			}

			return nil
		},
	}
}
