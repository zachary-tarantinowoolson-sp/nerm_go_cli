/*
Copyright © 2024 Zachary Tarantino-Woolson <zachary.tarantino@sailpoint.com>
*/
package advanced_search

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

func newAdvancedSearchRunCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "run",
		Short:   "Pulls Profiles from current environment",
		Long:    "Pulls Profiles from current environment based on an advanced Search. Stores data in a CSV and JSON file at the defaul output location",
		Example: "nerm advsearch run --id 123",
		Aliases: []string{"r"},
		RunE: func(cmd *cobra.Command, args []string) error {

			outputLoc := configs.GetOutputFolder() + configs.GetCurrentEnvironment() + "_AdvancedSearch_Export" + strconv.Itoa(int(time.Now().Unix()))

			params := url.Values{}

			id := cmd.Flags().Lookup("id").Value.String()
			limit := cmd.Flags().Lookup("limit").Value.String()
			getLimit := cmd.Flags().Lookup("get_limit").Value.String()

			limitInt, err := strconv.Atoi(limit)
			utilities.CheckError(err)

			if limitInt > 100 {
				fmt.Println("Limit can not be over 100 - Setting it back down to 100")
				params.Add("limit", "100")
			} else {
				params.Add("limit", limit)
			}
			params.Add("offset", "0")

			var getLimitInt int

			if getLimit != "" {
				getLimitInt, err = strconv.Atoi(getLimit)
				utilities.CheckError(err)
			} else {
				getLimitInt = math.MaxInt
			}

			var resp []byte
			var requestErr error

			createAdvancedSearchJsonFile(outputLoc + ".json")

			bar := progressbar.Default(-1, "Getting Profiles...")

			for offset := 0; offset < getLimitInt; offset = offset + limitInt {

				params.Set("offset", strconv.Itoa(offset))

				resp, requestErr = utilities.RunAdvSearchRequest(id, params.Encode())
				utilities.CheckError(requestErr)

				var advSearch_result ProfileResponse

				err := json.Unmarshal(resp, &advSearch_result)
				utilities.CheckError(err)

				// check to compensate for advanced search returning a 200 with empty body
				// when no profiles are found.. If there are profiles, prep file
				if len(advSearch_result.Profiles) == 0 {
					break
				} else if offset != 0 { //not first loop, so dont add comma to json file
					addCommaToFile(outputLoc + ".json")
				}

				printJsonToFile(outputLoc+".json", advSearch_result)

				bar.Add(len(advSearch_result.Profiles))
			}

			endAdvancedSearchJsonFile(outputLoc + ".json")

			convertJSONToCSV(outputLoc+".json", outputLoc+".csv")

			fmt.Println("\n\n\n" + "Profile data stored in " + outputLoc)

			return nil
		},
	}
	cmd.Flags().StringP("id", "i", "", "ID of a specific advanced Search")
	cmd.Flags().StringP("limit", "l", strconv.Itoa(configs.GetDefaultLimitParam()), "Limit for each GET request")
	cmd.Flags().StringP("get_limit", "g", "", "Set a Get limit for how many profiles to pull back (default is All profiles)")

	cmd.MarkFlagRequired("id")

	return cmd
}
