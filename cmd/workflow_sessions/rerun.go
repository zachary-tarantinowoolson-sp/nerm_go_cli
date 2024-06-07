/*
Copyright © 2024 Zachary Tarantino-Woolson <zachary.tarantino@sailpoint.com>
*/

// CLONED FROM GET - NEEDS TO BE UPDATED STILL
package workflow_sessions

import (
	"encoding/json"
	"nerm/cmd/utilities"
	"net/url"

	"github.com/spf13/cobra"
)

func newSessionsRerunCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "Rerun",
		Short:   "Re-runs a given workflow session, using the same data it already has",
		Long:    "Given a list or file of IDs of sessions, the workflow sessions will be re-run (a new session with the same data is created)",
		Example: "nerm sessions rerun -i 1234abcd-1234-abcd-5678-12345abcd5678 | nerm sessions rerun -f id_list.csv",
		Aliases: []string{"g"},
		RunE: func(cmd *cobra.Command, args []string) error {

			var resp []byte
			var requestErr error
			var err error

			params := url.Values{}

			id := cmd.Flags().Lookup("id").Value.String()
			// fileLoc := cmd.Flags().Lookup("file").Value.String()

			// make first call to get the total number of sessisions to be returned
			resp, requestErr = utilities.MakeAPIRequests("get", "workflow_sessions", id, params.Encode(), nil)
			utilities.CheckError(requestErr)

			var sessions_result SessionResponse

			err = json.Unmarshal(resp, &sessions_result)
			utilities.CheckError(err)

			return nil
		},
	}
	cmd.Flags().StringP("id", "i", "", "ID of a specific Workflow Session")
	cmd.Flags().StringP("file", "f", "", "CSV File that contains workflow sessions IDs")

	return cmd
}
