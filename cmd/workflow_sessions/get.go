/*
Copyright © 2024 Zachary Tarantino-Woolson <zachary.tarantino@sailpoint.com>
*/
package workflow_sessions

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

func newSessionsGetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get",
		Short:   "Gets Workflow Sessions from current environment",
		Long:    "Pulls Workflow Sessions from current environment based on query parameters. Stores data in a CSV and JSON file at the defaul output location",
		Example: "nerm profiles get --profile_id 1234abcd-1234-abcd-5678-12345abcd5678 ",
		Aliases: []string{"g"},
		RunE: func(cmd *cobra.Command, args []string) error {
			id := cmd.Flags().Lookup("id").Value.String()
			uid := cmd.Flags().Lookup("uid").Value.String()
			profile_id := cmd.Flags().Lookup("profile_id").Value.String()
			status := cmd.Flags().Lookup("status").Value.String()
			workflow_id := cmd.Flags().Lookup("workflow_id").Value.String()
			requester_id := cmd.Flags().Lookup("requester_id").Value.String()
			limit := cmd.Flags().Lookup("limit").Value.String()
			getLimit := cmd.Flags().Lookup("get_limit").Value.String()

			limitInt, _ := strconv.Atoi(limit)

			var getLimitInt int

			if getLimit != "" {
				getLimitInt, _ = strconv.Atoi(getLimit)
			} else {
				getLimitInt = math.MaxInt
			}

			outputLoc := configs.GetOutputFolder() + configs.GetCurrentEnvironment() + "_Sessions_Export" + strconv.Itoa(int(time.Now().Unix()))

			var resp []byte
			var requestErr error

			createSessionsJsonFile(outputLoc + ".json")

			params := url.Values{}
			params.Add("metadata", "true") //  include metadata for limit/offsets

			// add params that have been set
			if uid != "" {
				params.Add("uid", uid)
			}
			if profile_id != "" {
				params.Add("profile_id", profile_id)
			}
			if status != "" {
				params.Add("status", status)
			}
			if workflow_id != "" {
				params.Add("workflow_id", workflow_id)
			}
			if requester_id != "" {
				params.Add("requester_id", requester_id)
			}

			// make first call to get the total number of sessisions to be returned
			params.Add("limit", "1")
			resp, requestErr = utilities.MakeAPIRequests("get", "workflow_sessions", id, params.Encode(), nil)
			if limitInt > 500 {
				fmt.Println("Limit can not be over 500")
				limit = "500"
				params.Set("limit", "500")
			} else {
				params.Set("limit", limit)
			}
			utilities.CheckError(requestErr)

			var sessions_result SessionResponse
			var respMetaData ResponseMetaData

			err := json.Unmarshal(resp, &sessions_result)
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
				var sessions SessionResponse

				params.Add("offset", strconv.Itoa(offset))

				resp, requestErr = utilities.MakeAPIRequests("get", "workflow_sessions", id, params.Encode(), nil)

				utilities.CheckError(requestErr)

				err := json.Unmarshal(resp, &sessions)
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

				printJsonToFile(outputLoc+".json", sessions, lastLoop)
			}

			endSessionsJsonFile(outputLoc + ".json")

			convertJSONToCSV(outputLoc+".json", outputLoc+".csv")

			fmt.Println("\n" + "Session data stored in " + outputLoc)

			return nil
		},
	}
	cmd.Flags().StringP("id", "i", "", "ID of a specific Workflow Session")
	cmd.Flags().StringP("uid", "u", "", "UID of a specific Workflow Session")
	cmd.Flags().StringP("profile_id", "p", "", "Find all sessions that were run for a Profile")
	cmd.Flags().StringP("status", "s", "", "Status of the Workflow Session")
	cmd.Flags().StringP("workflow_id", "w", "", "ID of a specific Workflow that was run")
	cmd.Flags().StringP("requester_id", "r", "", "Find all sessions that were run by a specific User")
	cmd.Flags().StringP("limit", "l", strconv.Itoa(configs.GetDefaultLimitParam()), "Limit for each GET request")
	cmd.Flags().StringP("get_limit", "g", "", "Set a Get limit for how many sessions to pull back (default is All sessions)")

	return cmd
}
