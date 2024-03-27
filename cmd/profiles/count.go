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

type ProfileTypeResponse struct {
	ProfileTypes []struct {
		Name     string `json:"name"`
		ID       string `json:"id"`
		UID      string `json:"uid"`
		Category string `json:"category"`
	} `json:"profile_types"`
}

type ProfileResponse struct {
	Profiles []struct {
		ID               string `json:"id"`
		UID              string `json:"uid"`
		Name             string `json:"name"`
		ProfileTypeID    string `json:"profile_type_id"`
		Status           string `json:"status"`
		IDProofingStatus string `json:"id_proofing_status"`
		UpdatedAt        string `json:"updated_at"`
		CreatedAt        string `json:"created_at"`
	} `json:"profiles"`
	Metadata struct {
		Limit  int    `json:"limit"`
		Offset int    `json:"offset"`
		Total  int    `json:"total"`
		Next   string `json:"next"`
	} `json:"_metadata"`
}

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

			// fmt.Println(string(types_resp))

			for _, rec := range result.ProfileTypes {
				fmt.Println(rec.Name)
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
					err = json.Unmarshal(resp, &profile_result)
					if err != nil { // Parse []byte to the go struct pointer
						fmt.Println("Can not unmarshal JSON")
					}

					endTotal = endTotal + profile_result.Metadata.Total
					fmt.Println(status+" Total: ", profile_result.Metadata.Total)

					params.Del(status) // remove status just in case to add the next one
				}

			}

			fmt.Println("Full Total: ", endTotal)

			return nil
		},
	}
}
