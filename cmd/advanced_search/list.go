/*
Copyright © 2024 Zachary Tarantino-Woolson <zachary.tarantino@sailpoint.com>
*/
package advanced_search

import (
	"encoding/json"
	"math"
	"nerm/cmd/utilities"
	"net/url"
	"strconv"

	"github.com/spf13/cobra"
)

func newAdvancedSearchListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists available saved Advanced Searches",
		Long:    "Lists available saved Advanced Searches from current environment.",
		Example: "nerm advsearch list",
		Aliases: []string{"l"},
		RunE: func(cmd *cobra.Command, args []string) error {
			

			limitInt := 100
			getLimitInt := math.MaxInt32

			var resp []byte
			var requestErr error
			var respMetaData ResponseMetaData
			var finalAdvSearches [][]string

			params := url.Values{}
			params.Add("metadata", "true") //  include metadata for limit/offsets
			params.Add("limit", "1")       // Get 1 to see how many there are total
			params.Add("offset", "0")

			resp, requestErr = utilities.MakeAPIRequests("get", "advanced_search", "", params.Encode(), nil)
			utilities.CheckError(requestErr)
			err := json.Unmarshal(resp, &respMetaData)
			utilities.CheckError(err)

			if getLimitInt > respMetaData.Metadata.Total {
				getLimitInt = respMetaData.Metadata.Total
			}

			params.Set("limit", "100")

			for offset := 0; offset < getLimitInt; offset = offset + limitInt {
				var adv_searches AdvancedSearchConfig

				params.Set("offset", strconv.Itoa(offset))

				resp, requestErr = utilities.MakeAPIRequests("get", "advanced_search", "", params.Encode(), nil)
				utilities.CheckError(requestErr)

				err = json.Unmarshal(resp, &adv_searches)
				utilities.CheckError(err)
				err = json.Unmarshal(resp, &respMetaData)
				utilities.CheckError(err)

				for _, rec := range adv_searches.AdvancedSearch {
					var row []string
					row = append(row, rec.Label)
					row = append(row, rec.ID)
					finalAdvSearches = append(finalAdvSearches, row)
				}

			}
			printAdvSearchListTable(finalAdvSearches)

			return nil
		},
	}

	return cmd
}
