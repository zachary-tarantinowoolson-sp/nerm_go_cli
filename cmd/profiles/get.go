/*
Copyright © 2024 Zachary Tarantino-Woolson <zachary.tarantino@sailpoint.com>
*/
package profiles

import (
	"encoding/json"
	"fmt"
	"math"
	"nerm/cmd/configs"
	"nerm/cmd/utilities"
	"net/url"
	"reflect"
	"strconv"
	"time"

	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

func newProfileGetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get",
		Short:   "Pulls Profiles from current environment",
		Long:    "Pulls Profiles from current environment based on query parameters. Stores data in a CSV and JSON file at the defaul output location",
		Example: "nerm profiles get --profile_type",
		Aliases: []string{"g"},
		PreRun: func(cmd *cobra.Command, args []string) {
			isafterIdSet := cmd.Flags().Lookup("after_id").Changed
			if isafterIdSet {
				cmd.MarkFlagRequired("keep_archived")
			}

		},
		RunE: func(cmd *cobra.Command, args []string) error {
			id := cmd.Flags().Lookup("id").Value.String()
			exclude := cmd.Flags().Lookup("exclude").Value.String()
			profile_type := cmd.Flags().Lookup("profile_type").Value.String()
			status := cmd.Flags().Lookup("status").Value.String()
			name := cmd.Flags().Lookup("name").Value.String()
			force_backend := cmd.Flags().Lookup("force_backend").Value.String()
			limit := cmd.Flags().Lookup("limit").Value.String()
			getLimit := cmd.Flags().Lookup("get_limit").Value.String()
			after_id := cmd.Flags().Lookup("after_id").Value.String()
			isafterIdSet := cmd.Flags().Lookup("after_id").Changed
			keep_archived, _ := cmd.Flags().GetBool("keep_archived")

			limitInt, _ := strconv.Atoi(limit)

			var getLimitInt int

			if getLimit != "" {
				getLimitInt, _ = strconv.Atoi(getLimit)
			} else {
				getLimitInt = math.MaxInt
			}

			outputLoc := configs.GetOutputFolder() + configs.GetCurrentEnvironment() + "_Profile_Export" + strconv.Itoa(int(time.Now().Unix()))

			var resp []byte
			var requestErr error

			createProfilesJsonFile(outputLoc + ".json")

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
			if limitInt > 500 {
				fmt.Println("Limit can not be over 500")
				limit = "500"
				params.Add("limit", "500")
			} else if limitInt > getLimitInt {
				params.Add("limit", getLimit)
			}
			if isafterIdSet {
				params.Add("after_id", after_id)
			}

			// make first call to get the total number of profiles to be returned
			resp, requestErr = utilities.MakeAPIRequests("get", "profiles", id, params.Encode(), nil)

			utilities.CheckError(requestErr)

			var profile_result ProfileResponse
			var respMetaData ResponseMetaData

			err := json.Unmarshal(resp, &profile_result)
			utilities.CheckError(err)

			err = json.Unmarshal(resp, &respMetaData)
			utilities.CheckError(err)

			if getLimitInt > respMetaData.Metadata.Total {
				// getLimit = strconv.Itoa(respMetaData.Metadata.Total)
				getLimitInt = respMetaData.Metadata.Total
			}

			bar := progressbar.Default(int64(getLimitInt)) // set progress to number of profile types found
			lastLoop := false                              // used to determine where to add commas inthe json file, and for progressbar

			for offset := 0; offset < getLimitInt; offset = offset + limitInt {

				if !isafterIdSet {
					params.Add("offset", strconv.Itoa(offset))
				}

				resp, requestErr = utilities.MakeAPIRequests("get", "profiles", id, params.Encode(), nil)

				utilities.CheckError(requestErr)

				var profile_result ProfileResponse
				var profile_filtered_result []struct {
					ID               string            `json:"id"`
					UID              string            `json:"uid"`
					Name             string            `json:"name"`
					ProfileTypeID    string            `json:"profile_type_id"`
					Status           string            `json:"status"`
					IDProofingStatus string            `json:"id_proofing_status"`
					Archived         bool              `json:"archived"`
					UpdatedAt        string            `json:"updated_at"`
					CreatedAt        string            `json:"created_at"`
					Attributes       map[string]string `json:"attributes"`
				}
				// var profileResultZero ProfileResponse
				var respMetaData ResponseMetaData

				err := json.Unmarshal(resp, &profile_result)
				utilities.CheckError(err)

				err = json.Unmarshal(resp, &respMetaData)
				utilities.CheckError(err)

				if reflect.DeepEqual(profile_result, ProfileResponse{}) {
					fmt.Println("\n\nNo more Profiles found!")
					break
				}

				// Use a loop to check for and only store Archived or Non-Archived proifles (depending on the flag set)
				// If both the record and the flag are true or both are false - store the profile
				for _, rec := range profile_result.Profiles {
					if rec.Archived && keep_archived {
						profile_filtered_result = append(profile_filtered_result, rec)
					} else if !rec.Archived && !keep_archived {
						profile_filtered_result = append(profile_filtered_result, rec)
					}
				}
				profile_result.Profiles = profile_filtered_result

				if (offset + limitInt) >= getLimitInt {
					bar.Set(getLimitInt)
					lastLoop = true
				}

				if lastLoop {
					bar.Set(getLimitInt)
				} else {
					bar.Add(limitInt) // increment progress
				}

				if isafterIdSet {
					if respMetaData.Metadata.AfterID == "null" || respMetaData.Metadata.AfterID == "" { // incase the metadata is broken, just use the last id
						e, eErr := strconv.Atoi(limit)
						utilities.CheckError(eErr)
						params.Add("after_id", profile_result.Profiles[e-1].ID)
					} else {
						params.Add("after_id", respMetaData.Metadata.AfterID) // use the metadata value when possible
					}
				}

				printJsonToFile(outputLoc+".json", profile_result, lastLoop)
			}

			// jsonData, _ := json.MarshalIndent(profile_result, "", "    ")
			// fmt.Println(string(jsonData))

			endProfilesJsonFile(outputLoc + ".json")

			convertJSONToCSV(outputLoc+".json", outputLoc+".csv")

			fmt.Println("\n" + "Profile data stored in " + outputLoc)

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
	cmd.Flags().StringP("get_limit", "g", "", "Set a Get limit for how many profiles to pull back (default is All profiles)")
	cmd.Flags().String("after_id", "", "Get all Profiles using the after_id pagination. Leave blank or add a value to start from")
	cmd.Flags().Bool("keep_archived", false, "When using after_id pagination, determin if you want to store records that are archived or not. Requried if using after_id")

	return cmd
}
