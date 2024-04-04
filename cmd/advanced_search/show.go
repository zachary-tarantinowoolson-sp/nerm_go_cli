/*
Copyright © 2024 Zachary Tarantino-Woolson <zachary.tarantino@sailpoint.com>
*/
package advanced_search

import (
	"encoding/json"
	"fmt"
	"nerm/cmd/utilities"
	"net/url"

	"github.com/spf13/cobra"
)

func newAdvancedSearchShowCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "show",
		Short:   "Shows the configuration of a an Advanced Search",
		Long:    "Shows the configuration of a an Advanced Search",
		Example: "nerm advsearch show --id 1234",
		Aliases: []string{"s"},
		RunE: func(cmd *cobra.Command, args []string) error {
			id := cmd.Flags().Lookup("id").Value.String()
			var adv_searches AdvancedSearchConfig

			params := url.Values{}
			params.Add("id", id)

			resp, requestErr := utilities.MakeAPIRequests("get", "advanced_search", "", params.Encode(), nil)
			utilities.CheckError(requestErr)
			err := json.Unmarshal(resp, &adv_searches)
			utilities.CheckError(err)

			formatted, err := json.MarshalIndent(adv_searches.AdvancedSearch, "", "  ")
			utilities.CheckError(err)
			fmt.Println(string(formatted))

			return nil
		},
	}
	cmd.Flags().StringP("id", "i", "", "ID of a specific Advanced Search")
	cmd.MarkFlagRequired("id")

	return cmd
}
