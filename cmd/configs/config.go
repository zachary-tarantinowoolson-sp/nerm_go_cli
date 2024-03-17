/*
Copyright © 2024 Zachary Tarantino-Woolson <zachary.tarantino@sailpoint.com>
*/
package configs

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// where the file is stored and what the file is
const (
	configFolder  = ".nerm"
	configEnvFile = "nerm_config.yaml"
)

type Environment struct {
	Tenant   string `mapstructure:"tenant"`
	APIToken string `mapstructure:"token"`
}

func InitConfig() error {

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	viper.AddConfigPath(filepath.Join(home, ".nerm"))
	viper.SetConfigName("nerm_config")
	viper.SetConfigType("yaml")
	viper.SetEnvPrefix("nerm")

	viper.SetDefault("BASEURL", "nonemployee.com")
	viper.SetDefault("DEFAULT_OUTPUT_LOCATION", "")
	viper.SetDefault("LIMIT", 100)
	viper.SetDefault("CURRENT_ENVIRONMENT", "")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Print("Config file not found")
		} else {
			return err
		}
	}

	return nil
}

func SaveConfig() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	if _, err := os.Stat(filepath.Join(home, configFolder)); os.IsNotExist(err) {
		err = os.Mkdir(filepath.Join(home, configFolder), 0777)
		if err != nil {
			log.Print("Failed to create .nerm folder for config: ", err)
		}
	}

	err = viper.WriteConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			err = viper.SafeWriteConfig()
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}

func SetCurrentEnvironment(current_environment string) {
	viper.Set("CURRENT_ENVIRONMENT", strings.ToLower(current_environment))
}
func SetTenant(tenant string) {
	viper.Set("ALL_ENVIRONMENTS."+GetCurrentEnvironment()+".TENANT", tenant)
}
func SetBaseURL(baseurl string) {
	viper.Set("BASEURL", baseurl)
}
func SetAPIToken(token string) {
	viper.Set("ALL_ENVIRONMENTS."+GetCurrentEnvironment()+".TOKEN", "Bearer "+token)
}
func SetOutputFolder(default_output_location string) {
	viper.Set("DEFAULT_OUTPUT_LOCATION", default_output_location)
}
func SetDefaultLimitParam(limit string) {
	viper.Set("LIMIT", limit)
}

func GetCurrentEnvironment() string {
	return strings.ToLower(viper.GetString("CURRENT_ENVIRONMENT"))
}
func GetTenant() string {
	return viper.GetString("ALL_ENVIRONMENTS." + GetCurrentEnvironment() + ".TENANT")
}
func GetBaseURL() string {
	return viper.GetString("BASEURL")
}
func GetAPIToken() string {
	return viper.GetString("ALL_ENVIRONMENTS." + GetCurrentEnvironment() + ".TOKEN")
}
func GetOutputFolder() string {
	return viper.GetString("DEFAULT_OUTPUT_LOCATION")
}
func GetDefaultLimitParam() int {
	return viper.GetInt("LIMIT")
}

func GetAllEnvironments() map[string]interface{} {
	return viper.GetStringMap("ALL_ENVIRONMENTS")
}
