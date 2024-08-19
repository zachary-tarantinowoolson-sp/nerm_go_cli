/*
Copyright © 2024 Zachary Tarantino-Woolson <zachary.tarantino@sailpoint.com>
*/
package advanced_search

import (
	"bufio"
	"encoding/json"
	"fmt"
	"nerm/cmd/utilities"
	"os"
	"strconv"
	"strings"

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
			prompt := cmd.Flags().Lookup("prompt").Value.String()

			if file != "" {
				adv_searches := readAdvancedSearchJsonFile(file)

				formatted, err := json.Marshal(adv_searches.AdvancedSearch)
				utilities.CheckError(err)

				json_string := "{\"advanced_search\":" + string(formatted) + "}"

				_, requestErr := utilities.MakeAPIRequests("post", "advanced_search", "", "", []byte(json_string))
				utilities.CheckError(requestErr)

				fmt.Println("Advanced Search uploaded. New List:")

				newAdvancedSearchListCommand().Execute()

			} else if prompt != "" {

				var adv_search AdvancedSearchConfigForID

				advancedSearchLabel := prompt

				types_resp, types_err := utilities.MakeAPIRequests("get", "profile_types", "", "", nil)
				utilities.CheckError(types_err)

				var result ProfileTypeResponse
				unmarshalErr := json.Unmarshal(types_resp, &result)
				utilities.CheckError(unmarshalErr)

				r := bufio.NewReader(os.Stdin)

				firstLoop := true

			outer:
				for {
					fmt.Fprint(os.Stderr, "What type of Condition Rule do you want to add to the '"+advancedSearchLabel+"' Advanced Searh?:\n1.Profile Type\n2.Profile Status\n3.Profile Attribute\n4.Risk\n5.Exit\n>")
					conditionRule, readErr := r.ReadString('\n')
					utilities.CheckError(readErr)

					conditionRule = strings.TrimSpace(conditionRule)

					var json_string string

					switch conditionRule {
					case "1", "Profile Type", "profile type", "type":
						fmt.Println("in 1")
						for index, rec := range result.ProfileTypes {
							fmt.Println(index+1, ". ", rec.Name)
						}

						fmt.Println("What Profile Type to add as a Condition Rule? (enter the number)")
						readType := bufio.NewReader(os.Stdin)
						profType, readErr := readType.ReadString('\n')
						utilities.CheckError(readErr)
						profType = strings.TrimSpace(profType)
						i, err := strconv.Atoi(profType)
						utilities.CheckError(err)

						fmt.Println(result.ProfileTypes[i-1])

						json_string = "{\"advanced_search\": {\"condition_rules_attributes\": [{\"type\":\"ProfileTypeRule\",\"comparison_operator\":\"==\",\"value\":\"" + result.ProfileTypes[i-1].ID + "\"}]}}"

					case "2", "Profile Status", "profile status", "status":
						fmt.Println("in 2")

					case "3", "Profile Attribute", "profile attribute", "attribute":
						fmt.Println("in 3")

					case "4", "Risk", "risk":
						fmt.Println("in 4")

					case "5", "Exit", "exit", "Quit", "quit", "e", "q":
						fmt.Println("in 5")
						break outer
					}

					// if not first loop, make: {"advanced_search": {"label": "No operation test"}}

					if firstLoop {
						resp, requestErr := utilities.MakeAPIRequests("post", "advanced_search", "", "", []byte("{\"advanced_search\": {\"label\": \""+advancedSearchLabel+"\"}}"))
						utilities.CheckError(requestErr)

						unmarshalErr = json.Unmarshal(resp, &adv_search)
						utilities.CheckError(unmarshalErr)

					}
					_, requestErr := utilities.MakeAPIRequests("patch", "advanced_search/"+adv_search.AdvancedSearch.ID, "", "", []byte(json_string))
					utilities.CheckError(requestErr)

					json_string = "" // reset string

				}

			} else {
				fmt.Println("Please provide at least one of the flags.")
			}

			// Available Rules:
			// {
			// 	"type": "ProfileTypeRule",
			// 	"comparison_operator": "==",
			// 	"value": "3fa85f64-5717-4562-b3fc-2c963f66afa6"
			//   },
			//   {
			// 	"type": "ProfileStatusRule",
			// 	"comparison_operator": "==",
			// 	"value": "Active"
			//   },
			//   {
			// 	"type": "ProfileAttributeRule",
			// 	"condition_object_type": "TextFieldAttribute",
			// 	"condition_object_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
			// 	"comparison_operator": "==",
			// 	"value": "Some value"
			//   },
			//   {
			// 	"type": "ProfileAttributeRule",
			// 	"condition_object_type": "DateAttribute",
			// 	"condition_object_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
			// 	"secondary_attribute_type": "DateAttribute",
			// 	"secondary_attribute_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
			// 	"comparison_operator": ">",
			// 	"value": "Today",
			// 	"secondary_value": "after",
			// 	"tertiary_value": 30
			//   },
			//   {
			// 	"type": "ProfileAttributeRule",
			// 	"condition_object_type": "ProfileSelectAttribute",
			// 	"condition_object_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
			// 	"comparison_operator": "include?",
			// 	"value": "3fa85f64-5717-4562-b3fc-2c963f66afa6"
			//   },
			//   {
			// 	"type": "RiskRule",
			// 	"comparison_operator": "==",
			// 	"value": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
			// 	"secondary_value": "OverallRisk"
			//   }

			return nil
		},
	}
	cmd.Flags().StringP("file", "f", "", "Use a file to create an Advanced Search. Specify file path here")
	cmd.Flags().StringP("prompt", "p", "", "Use Prompts to create an Advanced Search. Provide a name for the new Advanced Search")
	cmd.MarkFlagsOneRequired("file", "prompt")
	cmd.MarkFlagsMutuallyExclusive("file", "prompt")

	return cmd
}
