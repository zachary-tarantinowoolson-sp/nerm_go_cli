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
		fmt.Fprint(os.Stderr, "Please enter a new Base URL or press enter if the following is correct (https://tenant.{{baseurl}}.com) ("+baseurl+"):")
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
		fmt.Fprint(os.Stderr, "Please enter a new Base URL or press enter if the following is correct (https://tenant.{{baseurl}}.com) ("+baseurl+"):")
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

	url := "http://" + tenant + "." + baseurl + "/api/" + endpoint
	if req_id != "" {
		url = url + "/" + req_id
	} else {
		url = url + "?" + params
	}

	switch strings.ToLower(method) {
	case "get":
		// fmt.Println("url: ", url, " |   jsonstr: ", jsonStr)
		resp, resp_err := MakeGetRequest(url, jsonStr)

		return resp, resp_err
	case "post":

	case "patch":

	case "delete":

	default:
		fmt.Println("No Method given.")
	}

	return nil, nil
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
	if err != nil {
		log.Fatal(err)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	return respBody, nil
}
