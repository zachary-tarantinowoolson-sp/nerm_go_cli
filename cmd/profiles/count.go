/*
Copyright © 2024 Zachary Tarantino-Woolson <zachary.tarantino@sailpoint.com>
*/
package profiles

import (
	"encoding/json"
	"fmt"
	"log"
	"nerm/cmd/utilities"
	"net/url"
	"strconv"

	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

func newProfileCountCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "count",
		Short:   "Pulls a count of all Profiles in current environment",
		Long:    "Pulls a count of all Profiles in current environment by profile Type",
		Example: "nerm profiles count",
		Aliases: []string{"c"},
		RunE: func(cmd *cobra.Command, args []string) error {

			endTotal := 0
			allStatuses := [4]string{"Active", "Inactive", "On Leave", "Terminated"}

			var finalValues [][]string

			type_params := url.Values{}
			type_params.Add("limit", "100")
			type_params.Add("metadata", "true")

			types_resp, types_err := utilities.MakeAPIRequests("get", "profile_types", "", type_params.Encode(), nil)
			if types_err != nil {
				log.Fatal(types_err)
			}
			// dumped_res, dump_err := httputil.DumpResponse(types_resp, true)

			var result ProfileTypeResponse
			err := json.Unmarshal(types_resp, &result)
			if err != nil { // Parse []byte to the go struct pointer
				fmt.Println("Can not unmarshal JSON")
			}

			bar := progressbar.Default(int64(len(result.ProfileTypes))) // set progress to number of profile types found

			for _, rec := range result.ProfileTypes {
				bar.Add(1) // increment progress
				var typeValues []string
				runningTotal := 0

				typeValues = append(typeValues, string(rec.Name))

				// fmt.Println(rec.Name)
				params := url.Values{}
				params.Add("profile_type_id", string(rec.ID))
				params.Add("limit", "1")
				params.Add("metadata", "true")
				params.Add("exclude_attributes", "true")

				for _, status := range allStatuses {

					params.Add("status", string(status))

					resp, err := utilities.MakeAPIRequests("get", "profiles", "", params.Encode(), nil)
					if err != nil {
						log.Fatal(err)
					}
					// fmt.Println(string(resp))

					var profile_result ProfileResponse
					var respMetaData ResponseMetaData
					err = json.Unmarshal(resp, &profile_result)
					if err != nil { // Parse []byte to the go struct pointer
						fmt.Println("Can not unmarshal JSON", err)
					}
					err = json.Unmarshal(resp, &respMetaData)
					if err != nil { // Parse []byte to the go struct pointer
						fmt.Println("Can not unmarshal JSON", err)
					}

					typeValues = append(typeValues, strconv.Itoa(respMetaData.Metadata.Total))
					runningTotal = runningTotal + respMetaData.Metadata.Total

					params.Del(status) // remove status just in case to add the next one
				}

				typeValues = append(typeValues, strconv.Itoa(runningTotal))
				finalValues = append(finalValues, typeValues)
				endTotal = endTotal + runningTotal
			}

			printCountTable(finalValues)
			fmt.Println("Total of all Profiles: ", endTotal)

			return nil
		},
	}
}
