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
	"time"

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

				types_resp, types_err := utilities.MakeAPIRequests("get", "profile_types", "", "", nil) // get all profile types for later use
				utilities.CheckError(types_err)
				var result ProfileTypeResponse
				unmarshalErr := json.Unmarshal(types_resp, &result)
				utilities.CheckError(unmarshalErr)

				riskResp, riskErr := utilities.MakeAPIRequests("get", "risk_levels", "", "", nil) // get all risk Levels for later use
				utilities.CheckError(riskErr)
				var riskLevels RiskLevel
				riskUnmarshalErr := json.Unmarshal(riskResp, &riskLevels)
				utilities.CheckError(riskUnmarshalErr)

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

						json_string = "{\"advanced_search\": {\"condition_rules_attributes\": [{\"type\":\"ProfileTypeRule\",\"comparison_operator\":\"==\",\"value\":\"" + result.ProfileTypes[i-1].ID + "\"}]}}"

					case "2", "Profile Status", "profile status", "status":

						fmt.Println("What Profile Status to add as a Condition Rule? (Active, Inactive, Terminated, On Leave)")
						readType := bufio.NewReader(os.Stdin)
						profStatus, readErr := readType.ReadString('\n')
						utilities.CheckError(readErr)
						statusValue := strings.TrimSpace(profStatus)

						if statusValue != "Active" && statusValue != "Inactive" && statusValue != "Terminated" && statusValue != "On Leave" {
							fmt.Println(statusValue, "is not a valid Status value. Please enter Active, Inactive, Terminated, or On Leave (Capitlized)")
						} else {
							json_string = "{\"advanced_search\": {\"condition_rules_attributes\": [{\"type\":\"ProfileStatusRule\",\"comparison_operator\":\"==\",\"value\":\"" + statusValue + "\"}]}}"
						}

					case "3", "Profile Attribute", "profile attribute", "attribute":
						fmt.Println("What is the ID of the attribute you want to add to the search?: ")
						readID := bufio.NewReader(os.Stdin)
						attributeId, readErr := readID.ReadString('\n')
						utilities.CheckError(readErr)
						attributeId = strings.TrimSpace(attributeId)

						attrResp, attrErr := utilities.MakeAPIRequests("get", "ne_attributes/"+attributeId, "", "", nil)
						utilities.CheckError(attrErr)
						var attribute NeAttribute
						unmarshalErr := json.Unmarshal(attrResp, &attribute)
						utilities.CheckError(unmarshalErr)

						switch attribute.NeAttribute.Type {
						case "TextFieldAttribute":
							fmt.Println("What kind of comparsion do you want to make for the " + attribute.NeAttribute.Label + " value? (==, !=, >, <, start_with?, end_with?, include?): ")
							readCompare := bufio.NewReader(os.Stdin)
							compareOperator, readErr := readCompare.ReadString('\n')
							utilities.CheckError(readErr)
							compareOperator = strings.TrimSpace(compareOperator)

							if compareOperator != "==" && compareOperator != "!=" && compareOperator != ">" && compareOperator != "<" && compareOperator != "start_with?" && compareOperator != "end_with?" && compareOperator != "include?" {
								fmt.Println(compareOperator, "is not a valid Comparison Operator. Please enter ==, !=, >, <, start_with?, end_with?, or include?")
							} else {
								fmt.Println("What value do you want to look for with " + attribute.NeAttribute.Label + "?: ")
								readValue := bufio.NewReader(os.Stdin)
								attributeValue, readErr := readValue.ReadString('\n')
								utilities.CheckError(readErr)
								attributeValue = strings.TrimSpace(attributeValue)

								json_string = "{\"advanced_search\": {\"condition_rules_attributes\": [{\"type\": \"ProfileAttributeRule\",\"condition_object_type\": \"TextFieldAttribute\",\"condition_object_id\": \"" + attributeId + "\",\"comparison_operator\": \"" + compareOperator + "\",\"value\": \"" + attributeValue + "\"}]}}"
							}

						case "DateAttribute":
							var compareOperator string

							for {
								fmt.Println("What kind of comparsion do you want to make for the " + attribute.NeAttribute.Label + " value? (>, <, before, after, ==): ")
								readCompare := bufio.NewReader(os.Stdin)
								compareOperatorRead, readErr := readCompare.ReadString('\n')
								utilities.CheckError(readErr)
								compareOperator = strings.TrimSpace(compareOperatorRead)

								if compareOperator != ">" && compareOperator != "<" && compareOperator != "after" && compareOperator != "before" && compareOperator != "==" {
									fmt.Println(compareOperator, "is not a valid Comparison Operator. Please enter >, <, before, after, ==")
								} else {
									break
								}
							}

							if compareOperator == "before" || compareOperator == "after" {
							compareLoop:
								for {
									fmt.Println("Do you want to compare against 'Today', a date, or another attribute? (enter 'attribute' to compare against another attribute)")
									readValue := bufio.NewReader(os.Stdin)
									compareValue, readErr := readValue.ReadString('\n')
									utilities.CheckError(readErr)
									compareValue = strings.TrimSpace(compareValue)

									// try to parse for a date, used in a later check
									_, timeErr := time.Parse("01/02/2006", compareValue)

									if compareValue == "attribute" {
										for {
											fmt.Println("Enter the Attribute ID for the attribute you want to compare " + attribute.NeAttribute.Label + " with:")
											readAttributeIdValue := bufio.NewReader(os.Stdin)
											attrIDValue, readErr := readAttributeIdValue.ReadString('\n')
											utilities.CheckError(readErr)
											attrIDValue = strings.TrimSpace(attrIDValue)

											attrResp, attrErr := utilities.MakeAPIRequests("get", "ne_attributes/"+attrIDValue, "", "", nil)
											utilities.CheckError(attrErr)

											if attrErr != nil { // if there is an error
												fmt.Println("There was an issue with that attribute. Please enter a valid attribute ID.\n\n", attrResp)
											} else {
												json_string = "{\"advanced_search\": {\"condition_rules_attributes\": [{\"type\": \"ProfileAttributeRule\",\"condition_object_type\": \"DateAttribute\",\"condition_object_id\": \"" + attributeId + "\",\"secondary_attribute_type\": \"DateAttribute\",\"secondary_attribute_id\": \"" + attrIDValue + "\",\"comparison_operator\": \"" + compareOperator + "\"}]}}"
												break compareLoop
											}
										}

									} else if compareValue == "Today" {
										json_string = "{\"advanced_search\": {\"condition_rules_attributes\": [{\"type\": \"ProfileAttributeRule\",\"condition_object_type\": \"DateAttribute\",\"condition_object_id\": \"" + attributeId + "\",\"comparison_operator\": \"" + compareOperator + "\",\"value\": \"Today\"}]}}"

										break
									} else if timeErr != nil {
										fmt.Println("Please enter attribute, 'Today', or a valid date string")
									} else { // must be a date value
										json_string = "{\"advanced_search\": {\"condition_rules_attributes\": [{\"type\": \"ProfileAttributeRule\",\"condition_object_type\": \"DateAttribute\",\"condition_object_id\": \"" + attributeId + "\",\"comparison_operator\": \"" + compareOperator + "\",\"value\": \"" + compareValue + "\"}]}}"
										break
									}

								}

							}

							// fmt.Println("What value do you want to look for with " + attribute.NeAttribute.Label + "?: ")
							// readValue := bufio.NewReader(os.Stdin)
							// attributeValue, readErr := readValue.ReadString('\n')
							// utilities.CheckError(readErr)
							// attributeValue = strings.TrimSpace(attributeValue)

							// fmt.Println(compareOperator, attributeValue)

						case "ProfileSelectAttribute":

						}

					case "4", "Risk", "risk":
						fmt.Println("Risk Rules are based on the 'Risk Level' that a Profile has. Which Risk Level do you want to search for? (enter the number)")
						for index, rec := range riskLevels.RiskLevels {
							fmt.Println(index+1, ". ", rec.Label)
						}

						readType := bufio.NewReader(os.Stdin)
						risk, readErr := readType.ReadString('\n')
						utilities.CheckError(readErr)
						risk = strings.TrimSpace(risk)
						i, err := strconv.Atoi(risk)
						utilities.CheckError(err)

						json_string = "{\"advanced_search\": {\"condition_rules_attributes\": [{\"type\":\"RiskRule\",\"comparison_operator\":\"==\",\"value\":\"" + riskLevels.RiskLevels[i-1].ID + "\"}]}}"

					case "5", "Exit", "exit", "Quit", "quit", "e", "q":
						break outer
					}

					// if not first loop, make: {"advanced_search": {"label": "No operation test"}}

					if firstLoop {
						resp, requestErr := utilities.MakeAPIRequests("post", "advanced_search", "", "", []byte("{\"advanced_search\": {\"label\": \""+advancedSearchLabel+"\"}}"))
						utilities.CheckError(requestErr)

						unmarshalErr = json.Unmarshal(resp, &adv_search)
						utilities.CheckError(unmarshalErr)

						firstLoop = false

					}
					_, requestErr := utilities.MakeAPIRequests("patch", "advanced_search/"+adv_search.AdvancedSearch.ID, "", "", []byte(json_string))
					utilities.CheckError(requestErr)

					formatted, err := json.MarshalIndent(json_string, "", "  ")
					utilities.CheckError(err)
					fmt.Println(string(formatted))

					fmt.Println(adv_search.AdvancedSearch.ID)

					json_string = "" // reset string

				}

			} else {
				fmt.Println("Please provide at least one of the flags.")
			}

			// Available Rules
			// date TODO: > < == uses the 3 value keys
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

			return nil
		},
	}
	cmd.Flags().StringP("file", "f", "", "Use a file to create an Advanced Search. Specify file path here")
	cmd.Flags().StringP("prompt", "p", "", "Use Prompts to create an Advanced Search. Provide a name for the new Advanced Search")
	cmd.MarkFlagsOneRequired("file", "prompt")
	cmd.MarkFlagsMutuallyExclusive("file", "prompt")

	return cmd
}
