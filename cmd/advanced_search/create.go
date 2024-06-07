/*
Copyright © 2024 Zachary Tarantino-Woolson <zachary.tarantino@sailpoint.com>
*/
package advanced_search

import (
	"github.com/spf13/cobra"
)

func newAdvancedSearchCreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "run",
		Short:   "Pulls Profiles from current environment",
		Long:    "Pulls Profiles from current environment based on query parameters. Stores data in a CSV and JSON file at the defaul output location",
		Example: "nerm profiles run --profile_type",
		Aliases: []string{"g"},
		RunE: func(cmd *cobra.Command, args []string) error {
			// params := url.Values{}

			// var resp []byte
			// var requestErr error

			// fileLoc := cmd.Flags().Lookup("file").Value.String()

			// resp, requestErr = utilities.RunAdvSearchRequest(id, params.Encode())
			// utilities.CheckError(requestErr)

			// var advSearch_result ProfileResponse

			// err := json.Unmarshal(resp, &advSearch_result)
			// utilities.CheckError(err)
			return nil
		},
	}
	cmd.Flags().StringP("file", "f", "", "Use a file to create an Advanced Search. Specify file path here")

	return cmd
}
