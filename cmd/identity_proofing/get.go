/*
Copyright © 2024 Zachary Tarantino-Woolson <zachary.tarantino@sailpoint.com>
*/
package identity_proofing

import (
	"encoding/json"
	"fmt"
	"nerm/cmd/configs"
	"nerm/cmd/utilities"
	"net/url"
	"strconv"
	"time"

	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

func newIDProofingResultGetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get",
		Short:   "Pulls IDP Results from current environment",
		Long:    "Pulls Identity Proofing Results from current environment based on query parameters. Stores data in a CSV and JSON file at the defaul output location",
		Example: "nerm idproofing get --result fail",
		Aliases: []string{"g"},
		RunE: func(cmd *cobra.Command, args []string) error {
			profile_id := cmd.Flags().Lookup("profile_id").Value.String()
			workflow_session_id := cmd.Flags().Lookup("workflow_session_id").Value.String()
			result := cmd.Flags().Lookup("result").Value.String()
			limit := cmd.Flags().Lookup("limit").Value.String()

			limitInt, _ := strconv.Atoi(limit)

			var getLimitInt int

			outputLoc := configs.GetOutputFolder() + configs.GetCurrentEnvironment() + "_IDP_Export" + strconv.Itoa(int(time.Now().Unix()))

			var resp []byte
			var requestErr error

			createIdentityProofingJsonFile(outputLoc + ".json")

			params := url.Values{}
			params.Add("metadata", "true") // always include metadata for limit/offsets

			// add params that have been set
			if profile_id != "" {
				params.Add("profile_id", profile_id)
			}
			if workflow_session_id != "" {
				params.Add("workflow_session_id", workflow_session_id)
			}
			if result != "" {
				params.Add("result", result)
			}
			if limitInt > 100 {
				fmt.Println("Limit can not be over 500")
				limit = "100"
				params.Add("limit", "100")
			} else {
				params.Add("limit", limit)
			}

			// make first call to get the total number of results to be returned
			resp, requestErr = utilities.MakeAPIRequests("get", "identity_proofing_results", "", params.Encode(), nil)

			utilities.CheckError(requestErr)

			var idp_result IdentityProofingResponse
			var respMetaData ResponseMetaData

			err := json.Unmarshal(resp, &idp_result)
			utilities.CheckError(err)

			err = json.Unmarshal(resp, &respMetaData)
			utilities.CheckError(err)

			if getLimitInt > respMetaData.Metadata.Total {
				// getLimit = strconv.Itoa(respMetaData.Metadata.Total)
				getLimitInt = respMetaData.Metadata.Total
			}

			bar := progressbar.Default(int64(getLimitInt)) // set progress to number of results found

			for offset := 0; offset < getLimitInt; offset = offset + limitInt {

				params.Add("offset", strconv.Itoa(offset))

				resp, requestErr = utilities.MakeAPIRequests("get", "identity_proofing_results", "", params.Encode(), nil)

				utilities.CheckError(requestErr)

				err := json.Unmarshal(resp, &idp_result)
				utilities.CheckError(err)

				err = json.Unmarshal(resp, &respMetaData)
				utilities.CheckError(err)

				if (offset + limitInt) >= getLimitInt {
					bar.Set(getLimitInt)
				} else {
					bar.Add(limitInt) // increment progress
				}

				printJsonToFile(outputLoc+".json", idp_result)
			}

			endIdentityProofingJsonFile(outputLoc + ".json")

			convertJSONToCSV(outputLoc+".json", outputLoc+".csv")

			fmt.Println("\n" + "Identity Proofing data stored in " + outputLoc)

			return nil
		},
	}
	cmd.Flags().StringP("profile_id", "p", "", "ID of a specific Profile")
	cmd.Flags().StringP("workflow_session_id", "w", "", "ID of a specific Workflow Session")
	cmd.Flags().StringP("result", "r", "", "Find IDP results based on Pass/Fail")
	cmd.Flags().StringP("limit", "l", strconv.Itoa(configs.GetDefaultLimitParam()), "Limit for each GET request")

	return cmd
}
