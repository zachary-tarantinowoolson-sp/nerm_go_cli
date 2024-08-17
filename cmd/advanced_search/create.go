/*
Copyright © 2024 Zachary Tarantino-Woolson <zachary.tarantino@sailpoint.com>
*/
package advanced_search

import (
	"encoding/json"
	"fmt"
	"nerm/cmd/utilities"

	"github.com/spf13/cobra"
)

func newAdvancedSearchCreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Create a new advanced search",
		Long:    "Create and upload an advanced search. Can be done via a JSON file or via prompt (tbd)",
		Example: "nerm advsearch create --file search.json",
		Aliases: []string{"c"},
		RunE: func(cmd *cobra.Command, args []string) error {
			file := cmd.Flags().Lookup("file").Value.String()
			var adv_searches AdvancedSearchConfigForUpload

			// params := url.Values{}
			// params.Add("id", id)

			// resp, requestErr := utilities.MakeAPIRequests("get", "advanced_search", "", params.Encode(), nil)
			// utilities.CheckError(requestErr)
			// err := json.Unmarshal(resp, &adv_searches)
			// utilities.CheckError(err)

			adv_searches = readAdvancedSearchJsonFile(file)

			formatted, err := json.MarshalIndent(adv_searches.AdvancedSearch, "", "  ")
			utilities.CheckError(err)
			fmt.Println(string(formatted))

			return nil
		},
	}
	cmd.Flags().StringP("file", "f", "", "Use a file to create an Advanced Search. Specify file path here")

	return cmd
}
