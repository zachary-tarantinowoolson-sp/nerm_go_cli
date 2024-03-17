/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"log"
	"nerm/cmd/configs"
	"nerm/cmd/root"

	"github.com/spf13/cobra"
)

var rootCmd *cobra.Command

func init() {

	cobra.CheckErr(configs.InitConfig())
	rootCmd = root.NewRootCommand()
}

func main() {
	_ = rootCmd.Execute()

	if save_error := configs.SaveConfig(); save_error != nil {
		log.Print("Issue saving config file", "error", save_error)
	}
}
