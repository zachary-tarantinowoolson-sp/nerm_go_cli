/*
Copyright © 2024 Zachary Tarantino-Woolson <zachary.tarantino@sailpoint.com>
*/
package identity_proofing

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"nerm/cmd/utilities"
	"os"
	"slices"
	"strings"

	"github.com/fatih/color"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
)

type IdentityProofingResponse struct {
	IdentityProofingResults []struct {
		ID                       string            `json:"id"`
		IdentityProofingActionID string            `json:"identity_proofing_action_id"`
		WorkflowSessionID        string            `json:"workflow_session_id"`
		ProfileID                string            `json:"profile_id"`
		IdentityProofingWorkflow string            `json:"proofing_workflow"`
		Result                   string            `json:"result"`
		UpdatedAt                string            `json:"updated_at"`
		CreatedAt                string            `json:"created_at"`
		Attributes               map[string]string `json:"proofing_attributes"`
	} `json:"identity_proofing_results"`
}

type IdentityProofingJsonFileData struct {
	ID                       string            `json:"id"`
	IdentityProofingActionID string            `json:"identity_proofing_action_id"`
	WorkflowSessionID        string            `json:"workflow_session_id"`
	ProfileID                string            `json:"profile_id"`
	IdentityProofingWorkflow string            `json:"proofing_workflow"`
	Result                   string            `json:"result"`
	UpdatedAt                string            `json:"updated_at"`
	CreatedAt                string            `json:"created_at"`
	Attributes               map[string]string `json:"proofing_attributes"`
}

type ResponseMetaData struct {
	Metadata struct {
		Limit  int    `json:"limit"`
		Offset int    `json:"offset"`
		Total  int    `json:"total"`
		Next   string `json:"next"`
	} `json:"_metadata"`
}

func NewIdentityProofingCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "idproofing",
		Short:   "Get Identity Proofing results",
		Long:    "Get Identity Proofing results for a specified Profile, Workflow Session, or result",
		Example: "nerm idproofing get | nerm i get --result pass",
		Aliases: []string{"i"},
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.AddCommand(
		newIdentityProofingCountCommand(),
		newIDProofingResultGetCommand(),
	)

	return cmd
}

func idpResultCountTable(data [2]string) {
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New("Pass", "Fail")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
	tbl.AddRow(data[0], data[1])

	tbl.Print()
}

func createIdentityProofingJsonFile(fileLoc string) {

	file, _ := os.OpenFile(fileLoc, os.O_CREATE|os.O_TRUNC, os.ModePerm)
	defer file.Close()
	// file.WriteString(strings.Trim("{\"profiles\":[", "\""))
	file.WriteString(strings.Trim("[", "\""))
	defer file.Close()
}

func endIdentityProofingJsonFile(fileLoc string) {

	file, _ := os.OpenFile(fileLoc, os.O_APPEND, os.ModePerm)
	defer file.Close()
	file.WriteString(strings.Trim("]", "\""))
	defer file.Close()
}

func printJsonToFile(fileLoc string, jsonData IdentityProofingResponse) {

	file, _ := os.OpenFile(fileLoc, os.O_APPEND|os.O_CREATE, os.ModePerm)
	defer file.Close()
	encoder := json.NewEncoder(file)

	for i, rec := range jsonData.IdentityProofingResults {
		fmt.Println(rec)
		encoder.Encode(rec)
		if (i + 1) != len(jsonData.IdentityProofingResults) {
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

	var profileData []IdentityProofingJsonFileData
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

	header := []string{"ID", "IdentityProofingActionID", "WorkflowSessionID", "ProfileID", "IdentityProofingWorkflow", "Result", "UpdatedAt", "CreatedAt"}
	// for _, k := range keys {
	// 	header = append(header, k)
	// }

	header = append(header, keys...)

	err = writer.Write(header)
	utilities.CheckError(err)

	for _, r := range profileData {
		var csvRow []string
		csvRow = append(csvRow, r.ID, r.IdentityProofingActionID, r.WorkflowSessionID, r.ProfileID, r.IdentityProofingWorkflow, r.Result, r.UpdatedAt, r.CreatedAt)

		for j := 7; j < len(header); j++ {
			csvRow = append(csvRow, r.Attributes[header[j]])
		}
		err = writer.Write(csvRow)
		utilities.CheckError(err)
	}
	return nil
}
