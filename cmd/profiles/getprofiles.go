/*
Copyright © 2024 Zachary Tarantino-Woolson <zachary.tarantino@sailpoint.com>
*/
package profiles

import (
	"encoding/json"
	"fmt"
	"log"
	"nerm/cmd/configs"
	"nerm/cmd/utilities"
	"net/url"
	"strconv"

	"github.com/spf13/cobra"
)

func newProfileGetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get",
		Short:   "Pulls Profiles from current environment",
		Long:    "Pulls Profiles from current environment based on query parameters",
		Example: "nerm profiles get --profile_type",
		Aliases: []string{"g"},
		RunE: func(cmd *cobra.Command, args []string) error {
			id := cmd.Flags().Lookup("id").Value.String()
			exclude := cmd.Flags().Lookup("exclude").Value.String()
			profile_type := cmd.Flags().Lookup("profile_type").Value.String()
			status := cmd.Flags().Lookup("status").Value.String()
			name := cmd.Flags().Lookup("name").Value.String()
			force_backend := cmd.Flags().Lookup("force_backend").Value.String()
			limit := cmd.Flags().Lookup("limit").Value.String()
			fileLoc := cmd.Flags().Lookup("file").Value.String()

			var resp []byte
			var requestErr error

			createProfilesJsonFile(fileLoc)

			params := url.Values{}
			params.Add("metadata", "true") // always include metadata for limit/offsets

			// add params that have been set
			if exclude != "" {
				params.Add("exclude_attributes", exclude)
			}
			if profile_type != "" {
				params.Add("profile_type_id", profile_type)
			}
			if status != "" {
				params.Add("status", status)
			}
			if name != "" {
				params.Add("name", name)
			}
			if force_backend != "" {
				params.Add("force_backend", force_backend)
			}
			if intLimit, _ := strconv.Atoi(limit); intLimit > 500 {
				fmt.Println("Limit can not be over 500")
				limit = "500"
				params.Add("limit", "500")
			} else {
				params.Add("limit", limit)
			}

			if id != "" {
				resp, requestErr = utilities.MakeAPIRequests("get", "profiles", id, params.Encode(), nil)
			} else {
				resp, requestErr = utilities.MakeAPIRequests("get", "profiles", "", params.Encode(), nil)
			}

			if requestErr != nil {
				log.Fatal(requestErr)
			}

			var profile_result ProfileResponse
			var respMetaData ResponseMetaData
			err := json.Unmarshal(resp, &profile_result)
			if err != nil { // Parse []byte to the go struct pointer
				fmt.Println("Can not unmarshal JSON", err)
			}
			err = json.Unmarshal(resp, &respMetaData)
			if err != nil { // Parse []byte to the go struct pointer
				fmt.Println("Can not unmarshal JSON", err)
			}

			// fmt.Println(string(resp))
			// fmt.Println(profile_result)

			// jsonData, _ := json.MarshalIndent(profile_result, "", "    ")
			// fmt.Println(string(jsonData))

			printToFile(fileLoc, profile_result)

			// var finalValues [][]string

			endProfilesJsonFile(fileLoc)

			return nil
		},
	}
	cmd.Flags().StringP("id", "i", "", "ID of a specific Profile")
	cmd.Flags().StringP("exclude", "e", "", "Exclude attributes from response")
	cmd.Flags().StringP("profile_type", "t", "", "Profile type ID of profiles")
	cmd.Flags().StringP("status", "s", "", "Status of profiles")
	cmd.Flags().StringP("name", "n", "", "Name of the profile(s) to look for")
	cmd.Flags().StringP("force_backend", "b", "", "Force the Profile Service or Identity Suite controllers")
	cmd.Flags().StringP("limit", "l", strconv.Itoa(configs.GetDefaultLimitParam()), "Limit for each GET request")
	cmd.Flags().StringP("file", "f", "", "Set the output location of the Profile Data")

	cmd.MarkFlagRequired("file")

	return cmd
}
