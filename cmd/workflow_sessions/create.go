/*
Copyright © 2024 Zachary Tarantino-Woolson <zachary.tarantino@sailpoint.com>
*/

// CLONED FROM GET - NEEDS TO BE UPDATED STILL
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

func newSessionsCreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "Create",
		Short:   "Create Workflow Sessions",
		Long:    "Uses data from a JSON or CSV file to generate workflow sessions in the current tenant",
		Example: "nerm sessions create -f something.csv ",
		Aliases: []string{"g"},
		RunE: func(cmd *cobra.Command, args []string) error {

			var resp []byte
			var requestErr error
			var days int
			var err error

			params := url.Values{}

			id := cmd.Flags().Lookup("id").Value.String()
			uid := cmd.Flags().Lookup("uid").Value.String()
			profile_id := cmd.Flags().Lookup("profile_id").Value.String()
			status := cmd.Flags().Lookup("status").Value.String()
			workflow_id := cmd.Flags().Lookup("workflow_id").Value.String()
			requester_id := cmd.Flags().Lookup("requester_id").Value.String()
			limit := cmd.Flags().Lookup("limit").Value.String()
			getLimit := cmd.Flags().Lookup("get_limit").Value.String()

			dayString := cmd.Flags().Lookup("days").Value.String()

			if dayString != "" {
				days, err = strconv.Atoi(dayString)
				utilities.CheckError(err)
				params.Add("order", "created_at DESC")
			}

			limitInt, _ := strconv.Atoi(limit)

			var getLimitInt int

			if getLimit != "" {
				getLimitInt, _ = strconv.Atoi(getLimit)
			} else {
				getLimitInt = math.MaxInt
			}

			outputLoc := configs.GetOutputFolder() + configs.GetCurrentEnvironment() + "_Sessions_Export" + strconv.Itoa(int(time.Now().Unix()))

			createSessionsJsonFile(outputLoc + ".json")

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

			err = json.Unmarshal(resp, &sessions_result)
			utilities.CheckError(err)

			err = json.Unmarshal(resp, &respMetaData)
			utilities.CheckError(err)

			if getLimitInt > respMetaData.Metadata.Total {
				// getLimit = strconv.Itoa(respMetaData.Metadata.Total)
				getLimitInt = respMetaData.Metadata.Total
			}

			params.Set("metadata", "false") // metadata makes calls slow, we do not need it for the actual data gets

			bar := progressbar.Default(int64(getLimitInt)) // set progress to number of profile types found

			for offset := 0; offset < getLimitInt; offset = offset + limitInt {
				var sessions SessionResponse      // this round of sessions from Get
				var finalSessions SessionResponse // the sessions that will be put into the file

				params.Add("offset", strconv.Itoa(offset))

				resp, requestErr = utilities.MakeAPIRequests("get", "workflow_sessions", id, params.Encode(), nil)
				utilities.CheckError(requestErr)

				err := json.Unmarshal(resp, &sessions)
				utilities.CheckError(err)

				// err = json.Unmarshal(resp, &respMetaData)
				// utilities.CheckError(err)

				if (offset + limitInt) >= getLimitInt {
					bar.Set(getLimitInt)
				} else {
					bar.Add(limitInt) // increment progress
				}

				for _, rec := range sessions.Sessions {
					t := time.Now().AddDate(0, 0, (days * -1))
					compareDate := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()) // zeros out the day

					createdAtTime, dateErr := time.Parse(time.RFC3339, rec.CreatedAt)
					utilities.CheckError(dateErr)

					if createdAtTime.After(compareDate) {
						// fmt.Println("after")
						finalSessions.Sessions = append(finalSessions.Sessions, rec)
					}
				}

				printJsonToFile(outputLoc+".json", finalSessions)
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
	cmd.Flags().StringP("days", "d", "", "Pull sessions from the last x days")

	return cmd
}
