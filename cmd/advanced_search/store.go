/*
Copyright © 2024 Zachary Tarantino-Woolson <zachary.tarantino@sailpoint.com>
*/
package advanced_search

import (
	"encoding/json"
	"fmt"
	"nerm/cmd/configs"
	"nerm/cmd/utilities"
	"net/url"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

func newAdvancedSearchDownloadCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "download",
		Short:   "Download the configuration of a an Advanced Search",
		Long:    "Download the configuration of a an Advanced Search",
		Example: "nerm advsearch download --id 1234",
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

			outputLoc := configs.GetOutputFolder() + configs.GetCurrentEnvironment() + "_" + adv_searches.AdvancedSearch[0].Label + "_AdvancedSearch_Config" + strconv.Itoa(int(time.Now().Unix()))
			storeAdvancedSearchJsonFile(outputLoc+".json", adv_searches)

			fmt.Println("\n" + "Advanced Saerch config data stored in " + outputLoc)

			return nil
		},
	}
	cmd.Flags().StringP("id", "i", "", "ID of a specific Advanced Search")
	cmd.MarkFlagRequired("id")

	return cmd
}
