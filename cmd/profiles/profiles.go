/*
Copyright © 2024 Zachary Tarantino-Woolson <zachary.tarantino@sailpoint.com>
*/
package profiles

import (
	"github.com/fatih/color"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
)

type ProfileTypeResponse struct {
	ProfileTypes []struct {
		Name     string `json:"name"`
		ID       string `json:"id"`
		UID      string `json:"uid"`
		Category string `json:"category"`
	} `json:"profile_types"`
}

type ProfileResponse struct {
	Profiles []struct {
		ID               string            `json:"id"`
		UID              string            `json:"uid"`
		Name             string            `json:"name"`
		ProfileTypeID    string            `json:"profile_type_id"`
		Status           string            `json:"status"`
		IDProofingStatus string            `json:"id_proofing_status"`
		UpdatedAt        string            `json:"updated_at"`
		CreatedAt        string            `json:"created_at"`
		Attributes       map[string]string `json:"attributes"`
	} `json:"profiles"`
	Metadata struct {
		Limit  int    `json:"limit"`
		Offset int    `json:"offset"`
		Total  int    `json:"total"`
		Next   string `json:"next"`
	} `json:"_metadata"`
}

func NewProfilesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "profiles",
		Short:   "CRUD Profile data from an eviroment",
		Long:    "Create, read, update, and delete Profiles. This allows an admin to execute commands against a specified NERM tenant",
		Example: "nerm profiles count | nerm profiles get",
		Aliases: []string{"p"},
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.AddCommand(
		newProfileCountCommand(),
		newProfileGetCommand(),
	)

	return cmd
}

func printTable(data [][]string) {
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New("Profile Type", "Active", "Inactive", "On Leave", "Terminated", "Total")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for _, row := range data {
		tbl.AddRow(row[0], row[1], row[2], row[3], row[4], row[5])
	}

	tbl.Print()
}
