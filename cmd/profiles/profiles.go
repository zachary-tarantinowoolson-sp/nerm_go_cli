/*
Copyright © 2024 Zachary Tarantino-Woolson <zachary.tarantino@sailpoint.com>
*/
package profiles

import (
	"encoding/csv"
	"encoding/json"
	"nerm/cmd/utilities"
	"os"
	"slices"
	"strings"

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
}

type ProfileJsonFileData struct {
	ID               string            `json:"id"`
	UID              string            `json:"uid"`
	Name             string            `json:"name"`
	ProfileTypeID    string            `json:"profile_type_id"`
	Status           string            `json:"status"`
	IDProofingStatus string            `json:"id_proofing_status"`
	UpdatedAt        string            `json:"updated_at"`
	CreatedAt        string            `json:"created_at"`
	Attributes       map[string]string `json:"attributes"`
}

type ResponseMetaData struct {
	Metadata struct {
		Limit   int    `json:"limit"`
		Offset  int    `json:"offset"`
		Total   int    `json:"total"`
		Next    string `json:"next"`
		AfterID string `json:"after_id"`
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
		newProfileDiffCommand(),
		newProfileJSONtoCSVCommand(),
	)

	return cmd
}

func printCountTable(data [][]string) {
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New("Profile Type", "Active", "Inactive", "On Leave", "Terminated", "Total")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for _, row := range data {
		tbl.AddRow(row[0], row[1], row[2], row[3], row[4], row[5])
	}

	tbl.Print()
}

func printDiffTable(data [][]string) {
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New("Tenant", "Backend", "Total")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for _, row := range data {
		tbl.AddRow(row[0], row[1], row[2])
	}

	tbl.Print()
}

func createProfilesJsonFile(fileLoc string) {

	file, _ := os.OpenFile(fileLoc, os.O_CREATE|os.O_TRUNC, os.ModePerm)
	defer file.Close()
	// file.WriteString(strings.Trim("{\"profiles\":[", "\""))
	file.WriteString(strings.Trim("[", "\""))
	defer file.Close()
}

func endProfilesJsonFile(fileLoc string) {

	file, _ := os.OpenFile(fileLoc, os.O_APPEND, os.ModePerm)
	defer file.Close()
	file.WriteString(strings.Trim("]", "\""))
	defer file.Close()
}

func printJsonToFile(fileLoc string, jsonData ProfileResponse, lastLoop bool) {

	file, _ := os.OpenFile(fileLoc, os.O_APPEND|os.O_CREATE, os.ModePerm)
	defer file.Close()
	encoder := json.NewEncoder(file)

	for i, rec := range jsonData.Profiles {
		encoder.Encode(rec)
		if !lastLoop {
			file.WriteString(strings.Trim(",", "\""))
		} else if (i + 1) != len(jsonData.Profiles) {
			file.WriteString(strings.Trim(",", "\""))
		}
	}
	// encoder := json.NewEncoder(file)
	// encoder.Encode(jsonData)
	defer file.Close()
}

func convertJSONToCSV(source string, destination string) error {
	var keys []string

	// Read the JSON file into the struct array
	sourceFile, err := os.Open(source)
	utilities.CheckError(err)

	defer sourceFile.Close()

	var profileData []ProfileJsonFileData
	if err := json.NewDecoder(sourceFile).Decode(&profileData); err != nil {
		return err
	}

	for _, r := range profileData {
		for k := range r.Attributes {
			keys = append(keys, k)
		}
	}

	slices.Sort(keys)           // sort a-z
	keys = slices.Compact(keys) // remove duplicates

	// Create a new file to store CSV data
	outputFile, err := os.Create(destination)
	utilities.CheckError(err)

	defer outputFile.Close()

	// Write the header of the CSV file and the successive rows by iterating through the JSON struct array
	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	header := []string{"ID", "UID", "Name", "ProfileTypeID", "Status", "IDProofingStatus", "UpdatedAt", "CreatedAt"}
	// for _, k := range keys {
	// 	header = append(header, k)
	// }

	header = append(header, keys...)

	err = writer.Write(header)
	utilities.CheckError(err)

	for _, r := range profileData {
		var csvRow []string
		csvRow = append(csvRow, r.ID, r.UID, r.Name, r.ProfileTypeID, r.Status, r.IDProofingStatus, r.UpdatedAt, r.CreatedAt)

		for j := 7; j < len(header); j++ {
			csvRow = append(csvRow, r.Attributes[header[j]])
		}
		err = writer.Write(csvRow)
		utilities.CheckError(err)
	}
	return nil
}
