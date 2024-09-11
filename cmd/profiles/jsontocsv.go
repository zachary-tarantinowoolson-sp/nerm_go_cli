/*
Copyright © 2024 Zachary Tarantino-Woolson <zachary.tarantino@sailpoint.com>
*/
package profiles

import (
	"strings"

	"github.com/spf13/cobra"
)

func newProfileJSONtoCSVCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "convert",
		Short:   "Converts a JSON file of Profile data to a CSV",
		Long:    "Converts a JSON file of Profile data to a CSV",
		Example: "nerm profiles convert -f something.json",
		Aliases: []string{"d"},
		RunE: func(cmd *cobra.Command, args []string) error {
			file := cmd.Flags().Lookup("file").Value.String()

			convertJSONToCSV(file, strings.Replace(file, "json", "csv", 1))

			return nil
		},
	}

	cmd.Flags().StringP("file", "f", "", "Source File to convert from JSON to CSV")
	cmd.MarkFlagRequired("file")

	return cmd
}
