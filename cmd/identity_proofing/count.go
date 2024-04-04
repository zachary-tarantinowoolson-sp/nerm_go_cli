/*
Copyright © 2024 Zachary Tarantino-Woolson <zachary.tarantino@sailpoint.com>
*/
package identity_proofing

import (
	"encoding/json"
	"nerm/cmd/utilities"
	"net/url"
	"strconv"

	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

func newIdentityProofingCountCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "count",
		Short:   "Displays a table of IDP results",
		Long:    "Pulls a count of all IDP results in current environment by Pass/Fail",
		Example: "nerm idproofing count",
		Aliases: []string{"c"},
		RunE: func(cmd *cobra.Command, args []string) error {
			var metadata ResponseMetaData

			var finalValues [2]string

			params := url.Values{}
			params.Add("limit", "1")
			params.Add("metadata", "true")
			params.Add("result", "pass")

			bar := progressbar.Default(2)

			pass_resp, pass_err := utilities.MakeAPIRequests("get", "identity_proofing_results", "", params.Encode(), nil)
			utilities.CheckError(pass_err)

			err := json.Unmarshal(pass_resp, &metadata)
			utilities.CheckError(err)

			finalValues[0] = strconv.Itoa(metadata.Metadata.Total)

			bar.Add(1)

			params.Set("result", "fail")
			fail_resp, fail_err := utilities.MakeAPIRequests("get", "identity_proofing_results", "", params.Encode(), nil)
			utilities.CheckError(fail_err)

			err = json.Unmarshal(fail_resp, &metadata)
			utilities.CheckError(err)

			finalValues[1] = strconv.Itoa(metadata.Metadata.Total)

			bar.Add(1)

			idpResultCountTable(finalValues)

			return nil
		},
	}
}
