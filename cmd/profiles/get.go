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
		RunE: func(cmd *cobra.Command, args []string) error {
			id := cmd.Flags().Lookup("id").Value.String()
			exclude := cmd.Flags().Lookup("exclude").Value.String()
			profile_type := cmd.Flags().Lookup("profile_type").Value.String()
			status := cmd.Flags().Lookup("status").Value.String()
			name := cmd.Flags().Lookup("name").Value.String()
			force_backend := cmd.Flags().Lookup("force_backend").Value.String()
			limit := cmd.Flags().Lookup("limit").Value.String()
			getLimit := cmd.Flags().Lookup("get_limit").Value.String()

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
			} else {
				params.Add("limit", limit)
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

				params.Add("offset", strconv.Itoa(offset))

				resp, requestErr = utilities.MakeAPIRequests("get", "profiles", id, params.Encode(), nil)

				utilities.CheckError(requestErr)

				var profile_result ProfileResponse
				var respMetaData ResponseMetaData

				err := json.Unmarshal(resp, &profile_result)
				utilities.CheckError(err)

				err = json.Unmarshal(resp, &respMetaData)
				utilities.CheckError(err)

				if (offset + limitInt) >= getLimitInt {
					bar.Set(getLimitInt)
					lastLoop = true
				}

				if lastLoop {
					bar.Set(getLimitInt)
				} else {
					bar.Add(limitInt) // increment progress
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

	return cmd
}
