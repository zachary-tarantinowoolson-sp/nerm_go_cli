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

			var resp []byte
			var requestErr error

			params := url.Values{}
			params.Add("metadata", "true")

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
			if limit != "" {
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
			err := json.Unmarshal(resp, &profile_result)
			if err != nil { // Parse []byte to the go struct pointer
				fmt.Println("Can not unmarshal JSON")
			}

			// fmt.Println(string(resp))
			// fmt.Println(profile_result)

			j, _ := json.MarshalIndent(profile_result, "", "    ")

			fmt.Println(string(j))

			// var finalValues [][]string

			// for _, rec := range profile_result.Profiles {
			// 	var rowValues []string

			// 	rowValues = append(rowValues, rec.Attributes)
			// }

			return nil
		},
	}
	cmd.Flags().StringP("id", "i", "", "ID of a specific Profile")
	cmd.Flags().StringP("exclude", "e", "", "Exclude attributes from response")
	cmd.Flags().StringP("profile_type", "t", "", "Profile type ID of profiles")
	cmd.Flags().StringP("status", "s", "", "Status of profiles")
	cmd.Flags().StringP("name", "n", "", "Name of the profile(s) to look for")
	cmd.Flags().StringP("force_backend", "f", "", "Force the Profile Service or Identity Suite controllers")
	cmd.Flags().StringP("limit", "l", "", "Limit for each GET request")

	return cmd
}
