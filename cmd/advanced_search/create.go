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
			adv_searches := readAdvancedSearchJsonFile(file)

			formatted, err := json.Marshal(adv_searches.AdvancedSearch)
			utilities.CheckError(err)

			json_string := "{\"advanced_search\":" + string(formatted) + "}"
			// fmt.Println(json_string)

			_, requestErr := utilities.MakeAPIRequests("post", "advanced_search", "", "", []byte(json_string))
			utilities.CheckError(requestErr)

			fmt.Println("Advanced Search  uploaded. New List:")

			newAdvancedSearchListCommand().Execute()

			return nil
		},
	}
	cmd.Flags().StringP("file", "f", "", "Use a file to create an Advanced Search. Specify file path here")

	return cmd
}
