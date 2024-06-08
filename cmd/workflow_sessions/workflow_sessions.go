/*
Copyright © 2024 Zachary Tarantino-Woolson <zachary.tarantino@sailpoint.com>
*/
package workflow_sessions

import (
	"encoding/csv"
	"encoding/json"
	"nerm/cmd/utilities"
	"os"
	"slices"
	"strings"

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
	Profile struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"profile"`
}

type ResponseMetaData struct {
	Metadata struct {
		Limit  int    `json:"limit"`
		Offset int    `json:"offset"`
		Total  int    `json:"total"`
		Next   string `json:"next"`
	} `json:"_metadata"`
}

type SessionResponse struct { // full response with header
	Sessions []struct {
		ID            string            `json:"id"`
		UID           string            `json:"uid"`
		WorkflowID    string            `json:"workflow_id"`
		RequesterType string            `json:"requester_type"`
		RequesterID   string            `json:"requester_id"`
		ProfileID     string            `json:"profile_id"`
		Status        string            `json:"status"`
		UpdatedAt     string            `json:"updated_at"`
		CreatedAt     string            `json:"created_at"`
		Attributes    map[string]string `json:"attributes"`
	} `json:"workflow_sessions"`
}

type SessionJsonFileData struct { // individual sessions
	ID            string            `json:"id"`
	UID           string            `json:"uid"`
	WorkflowID    string            `json:"workflow_id"`
	RequesterType string            `json:"requester_type"`
	RequesterID   string            `json:"requester_id"`
	ProfileID     string            `json:"profile_id"`
	Status        string            `json:"status"`
	UpdatedAt     string            `json:"updated_at"`
	CreatedAt     string            `json:"created_at"`
	Attributes    map[string]string `json:"attributes"`
}

type UserResponse struct {
	User struct {
		ID    string `json:"id"`
		Type  string `json:"type"`
		Name  string `json:"name"`
		Email string `json:"email"`
		Login string `json:"login"`
	} `json:"user"`
}

func NewWorkflowSessionsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "sessions",
		Short:   "CRU Workflow Sessions",
		Long:    "Create, read, and update Workflow sessions. This allows an admin to execute commands against a specified NERM tenant",
		Example: "nerm sessions get | nerm sessions create -w 1234abcd-1234-abcd-5678-12345abcd5678 -p 1234abcd-1234-abcd-5678-12345abcd5678",
		Aliases: []string{"s"},
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.AddCommand(
		newSessionsGetCommand(),
		newSessionsCreateCommand(),
		newSessionsRerunCommand(),
	)

	return cmd
}

func createSessionsJsonFile(fileLoc string) {

	file, _ := os.OpenFile(fileLoc, os.O_CREATE|os.O_TRUNC, os.ModePerm)
	defer file.Close()
	// file.WriteString(strings.Trim("{\"profiles\":[", "\""))
	file.WriteString(strings.Trim("[", "\""))
	defer file.Close()
}

func endSessionsJsonFile(fileLoc string) {

	file, _ := os.OpenFile(fileLoc, os.O_APPEND, os.ModePerm)
	defer file.Close()
	file.WriteString(strings.Trim("]", "\""))
	defer file.Close()
}

func printJsonToFile(fileLoc string, jsonData SessionResponse) {

	file, _ := os.OpenFile(fileLoc, os.O_APPEND|os.O_CREATE, os.ModePerm)
	defer file.Close()
	encoder := json.NewEncoder(file)

	for i, rec := range jsonData.Sessions {
		encoder.Encode(rec)

		// fmt.Println(len(jsonData.Sessions), i, i+1)
		if (i + 1) != len(jsonData.Sessions) {
			file.WriteString(strings.Trim(",", "\""))
			// fmt.Println("in if", (i+1) != len(jsonData.Sessions))
		}
	}
	// encoder := json.NewEncoder(file)
	// encoder.Encode(jsonData)
}

func convertJSONToCSV(source string, destination string) error {
	var keys []string

	// Read the JSON file into the struct array
	sourceFile, err := os.Open(source)
	utilities.CheckError(err)

	defer sourceFile.Close()

	var sessionData []SessionJsonFileData
	if err := json.NewDecoder(sourceFile).Decode(&sessionData); err != nil {
		return err
	}

	for _, r := range sessionData { // attribute keys
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

	header := []string{"ID", "UID", "WorkflowID", "RequesterType", "RequesterID", "ProfileID", "Status", "UpdatedAt", "CreatedAt"}
	// for _, k := range keys {
	// 	header = append(header, k)
	// }
	header = append(header, keys...) // attribute keys
	err = writer.Write(header)
	utilities.CheckError(err)

	for _, r := range sessionData {
		var csvRow []string
		csvRow = append(csvRow, r.ID, r.UID, r.WorkflowID, r.RequesterType, r.RequesterID, r.ProfileID, r.Status, r.UpdatedAt, r.CreatedAt)

		for j := 9; j < len(header); j++ {
			csvRow = append(csvRow, r.Attributes[header[j]])
		}
		err = writer.Write(csvRow)
		utilities.CheckError(err)
	}
	return nil
}
