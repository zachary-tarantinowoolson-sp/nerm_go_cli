/*
Copyright © 2024 Zachary Tarantino-Woolson <zachary.tarantino@sailpoint.com>
*/
package profiles

import (
	"encoding/json"
	"fmt"
	"nerm/cmd/configs"
	"nerm/cmd/utilities"
	"net/url"
	"sort"
	"strconv"

	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
	"golang.org/x/exp/maps"
)

func newProfileDiffCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "diff",
		Short:   "Pulls a count of total Profiles via the Suite and Profile Service",
		Long:    "Pulls a count of total Profiles via the Suite and Profile Service. This is to check for a different in the number of profiles",
		Example: "nerm profiles diff",
		Aliases: []string{"d"},
		RunE: func(cmd *cobra.Command, args []string) error {
			allEnvs, allErr := cmd.Flags().GetBool("all_envs")
			utilities.CheckError(allErr)

			backend := [2]string{"suite", "profile_service"}
			var finalValues [][]string

			currentEnv := configs.GetCurrentEnvironment() // store current env

			if allEnvs {
				allEnvironmentNames := configs.GetAllEnvironmentStrings()

				keys := maps.Keys(allEnvironmentNames)
				sort.Strings(keys)

				for _, e := range keys {
					configs.SetCurrentEnvironment(e)
					// fmt.Println("|| maps values:", maps.Values(allEnvironmentNames))
					fmt.Println("|-+-| Tenant: ", e)

					bar := progressbar.Default(2) // set progress to number of profile types found

					for _, rec := range backend {
						bar.Add(1) // increment progress
						var totalValues []string

						totalValues = append(totalValues, e)

						totalValues = append(totalValues, rec)
						params := url.Values{}
						params.Add("limit", "1")
						params.Add("exclude_attributes", "true")
						params.Add("metadata", "true")
						params.Add("force_backend", rec)

						resp, err := utilities.MakeAPIRequests("get", "profiles", "", params.Encode(), nil)
						utilities.CheckError(err)

						var respMetaData ResponseMetaData
						err = json.Unmarshal(resp, &respMetaData)
						utilities.CheckError(err)

						totalValues = append(totalValues, strconv.Itoa(respMetaData.Metadata.Total))
						finalValues = append(finalValues, totalValues)
					}
				}

				configs.SetCurrentEnvironment(currentEnv) // reset tenant

			} else {

				bar := progressbar.Default(2) // set progress to number of profile types found

				for _, rec := range backend {
					bar.Add(1) // increment progress
					var totalValues []string

					totalValues = append(totalValues, currentEnv)
					totalValues = append(totalValues, rec)
					params := url.Values{}
					params.Add("limit", "1")
					params.Add("metadata", "true")
					params.Add("force_backend", rec)

					resp, err := utilities.MakeAPIRequests("get", "profiles", "", params.Encode(), nil)
					utilities.CheckError(err)

					var respMetaData ResponseMetaData
					err = json.Unmarshal(resp, &respMetaData)
					utilities.CheckError(err)

					totalValues = append(totalValues, strconv.Itoa(respMetaData.Metadata.Total))
					finalValues = append(finalValues, totalValues)
				}
			}

			printDiffTable(finalValues)

			return nil
		},
	}

	cmd.Flags().BoolP("all_envs", "a", false, "If true, Runs the Diff on all configured environments. Else, just the current one")
	return cmd
}
