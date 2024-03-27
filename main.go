/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"fmt"
	"log"
	"nerm/cmd/configs"
	"nerm/cmd/root"
	"runtime"

	"github.com/spf13/cobra"
)

var rootCmd *cobra.Command

func init() {

	cobra.CheckErr(configs.InitConfig())
	rootCmd = root.NewRootCommand()
}

func main() {
	PrintMemUsage()
	_ = rootCmd.Execute()

	PrintMemUsage()
	if save_error := configs.SaveConfig(); save_error != nil {
		log.Print("Issue saving config file", "error", save_error)
	}
}

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc = %v MiB", m.Alloc/1024/1024)
	fmt.Printf("\tTotalAlloc = %v MiB", m.TotalAlloc/1024/1024)
	fmt.Printf("\tSys = %v MiB", m.Sys/1024/1024)
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}
