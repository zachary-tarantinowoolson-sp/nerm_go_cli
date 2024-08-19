/*
Copyright © 2024 Zachary Tarantino-Woolson <zachary.tarantino@sailpoint.com>
*/
package utilities

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"nerm/cmd/configs"
	"net/http"
	"os"
	"strings"
)

type Environment struct {
	TenantURL string `mapstructure:"tenanturl"`
	BaseURL   string `mapstructure:"baseurl"`
	APIToken  string `mapstructure:"token"`
}

func init() {
	// rootCmd.AddCommand(utilitiesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// utilitiesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// utilitiesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func FindEnvironments() map[string]interface{} {

	env_names := configs.GetAllEnvironments()

	return env_names

}

func CreateEnvironment(environmentName string) error {
	existing_envs := FindEnvironments()
	var tenant string
	var token string
	var baseurl string

	if existing_envs[environmentName] != nil {
		log.Print(environmentName + " already exists. Please use 'nerm env update' to update that environment")
		return nil
	}

	tenant = environmentName
	baseurl = configs.GetBaseURL()

	var s string
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprint(os.Stderr, "Please enter the Tenant name or press enter if the following is correct (https://{{tenant}}.nonemployee.com) ("+tenant+"):")
		s, _ = r.ReadString('\n')
		if s != "" {
			break
		}
	}
	prompt_input := strings.TrimSpace(s)

	if prompt_input != "" {
		tenant = prompt_input
	}

	var y string
	re := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprint(os.Stderr, "Please enter a new Base URL or press enter if the following is correct (https://tenant.{{baseurl}}) ("+baseurl+"):")
		y, _ = re.ReadString('\n')
		if y != "" {
			break
		}
	}
	prompt_input_3 := strings.TrimSpace(y)

	if prompt_input_3 != "" {
		baseurl = prompt_input_3
	}

	var x string
	a := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprint(os.Stderr, "Please enter a bearer token to use with this tenant (Bearer: {{token}}) (Please only enter the token value, not 'Bearer'):")
		x, _ = a.ReadString('\n')
		if x != "" {
			break
		}
	}
	prompt_input_2 := strings.TrimSpace(x)

	if prompt_input_2 != "" {
		token = prompt_input_2
	}

	configs.SetCurrentEnvironment(environmentName)
	configs.SetTenant(tenant)
	configs.SetBaseURL(baseurl)
	configs.SetAPIToken(token)

	return nil
}

func UpdateEnvironment(environmentName string) error {
	var tenant string
	var token string
	var baseurl string

	tenant = environmentName
	baseurl = configs.GetBaseURL()

	var s string
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprint(os.Stderr, "Please enter a new Tenant name or press enter if the following is correct (https://{{tenant}}.nonemployee.com) ("+tenant+"):")
		s, _ = r.ReadString('\n')
		if s != "" {
			break
		}
	}
	prompt_input := strings.TrimSpace(s)

	if prompt_input != "" {
		tenant = prompt_input
	}

	var y string
	re := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprint(os.Stderr, "Please enter a new Base URL or press enter if the following is correct (https://tenant.{{baseurl}}) ("+baseurl+"):")
		y, _ = re.ReadString('\n')
		if y != "" {
			break
		}
	}
	prompt_input_3 := strings.TrimSpace(y)

	if prompt_input_3 != "" {
		baseurl = prompt_input_3
	}

	var x string
	a := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprint(os.Stderr, "Please enter a new bearer token to use with this tenant (Bearer: {{token}}) (Please only enter the token value, not 'Bearer'):")
		x, _ = a.ReadString('\n')
		if x != "" {
			break
		}
	}
	prompt_input_2 := strings.TrimSpace(x)

	if prompt_input_2 != "" {
		token = prompt_input_2
	}

	configs.SetCurrentEnvironment(environmentName)
	configs.SetTenant(tenant)
	configs.SetBaseURL(baseurl)
	configs.SetAPIToken(token)

	return nil
}

func MakeAPIRequests(method string, endpoint string, req_id string, params string, jsonStr []byte) ([]byte, error) {
	tenant := configs.GetTenant()
	baseurl := configs.GetBaseURL()

	url := "https://" + tenant + "." + baseurl + "/api/" + endpoint
	if req_id != "" {
		url = url + "/" + req_id + "?" + params
	} else {
		url = url + "?" + params
	}

	switch strings.ToLower(method) {
	case "get":
		// fmt.Println("url: ", url, " |   jsonstr: ", jsonStr)
		resp, resp_err := MakeGetRequest(url, jsonStr)

		return resp, resp_err
	case "post":
		resp, resp_err := MakePostRequest(url, jsonStr)

		return resp, resp_err

	case "patch":
		resp, resp_err := MakePatchRequest(url, jsonStr)

		return resp, resp_err

	case "delete":

	default:
		fmt.Println("No Method given.")
	}

	return nil, nil
}

func MakePatchRequest(url string, jsonStr []byte) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Authorization", configs.GetAPIToken())
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	CheckError(err)

	respBody, err := io.ReadAll(resp.Body)
	CheckError(err)
	defer resp.Body.Close()

	return respBody, nil
}

func MakePostRequest(url string, jsonStr []byte) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Authorization", configs.GetAPIToken())
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	CheckError(err)

	respBody, err := io.ReadAll(resp.Body)
	CheckError(err)
	defer resp.Body.Close()

	return respBody, nil
}

func MakeGetRequest(url string, jsonStr []byte) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Authorization", configs.GetAPIToken())
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	CheckError(err)

	respBody, err := io.ReadAll(resp.Body)
	CheckError(err)
	defer resp.Body.Close()

	return respBody, nil
}

func RunAdvSearchRequest(req_id string, params string) ([]byte, error) {
	tenant := configs.GetTenant()
	baseurl := configs.GetBaseURL()
	url := "https://" + tenant + "." + baseurl + "/api/advanced_search/" + req_id + "/run?" + params

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Authorization", configs.GetAPIToken())
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	CheckError(err)

	respBody, err := io.ReadAll(resp.Body)
	CheckError(err)
	defer resp.Body.Close()

	return respBody, nil
}

func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

/*
func convertJSONToCSV(source string, destination string) error {
	// 2. Read the JSON file into the struct array
	sourceFile, err := os.Open(source)
	utilities.CheckError(err)

	defer sourceFile.Close()

	var profileData []ProfileJsonFileData
	if err := json.NewDecoder(sourceFile).Decode(&profileData); err != nil {
		return err
	}

	// Attributes map[string]string `json:"attributes"`

	var keys []string

	for _, r := range profileData {
		// fmt.Println(r.Attributes, r.Attributes["guests_name_ne_attribute"], r.Attributes["first_name"])
		for k, _ := range r.Attributes {
			keys = append(keys, k)
		}
	}

	// 3. Create a new file to store CSV data
	outputFile, err := os.Create(destination)
	utilities.CheckError(err)

	defer outputFile.Close()

	// 4. Write the header of the CSV file and the successive rows by iterating through the JSON struct array
	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	header := []string{"ID", "UID", "Name", "ProfileTypeID", "Status", "IDProofingStatus", "UpdatedAt", "CreatedAt"}
	for _, k := range keys {
		header = append(header, k)
	}
	err = writer.Write(header)
	utilities.CheckError(err)

	for _, r := range profileData {
		var csvRow []string
		csvRow = append(csvRow, r.ID, r.UID, r.Name, r.ProfileTypeID, r.IDProofingStatus, r.UpdatedAt, r.CreatedAt)

		for j := 7; j < len(header); j++ {
			csvRow = append(csvRow, r.Attributes[header[j]])
		}
		err = writer.Write(csvRow)
		utilities.CheckError(err)
	}
	return nil
}
*/
