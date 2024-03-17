/*
Copyright © 2024 Zachary Tarantino-Woolson <zachary.tarantino@sailpoint.com>
*/
package cmd

import (
	"log"
	"os"
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

func FindEnvironments() {
	files, err := os.ReadDir("/environment/environment_files/")
	if err != nil {
		log.Fatal(err)
	}
	return files

}

func CreateEnvironment(environmentName string) {
	existing_envs := utilities.FindEnvironments()

}
