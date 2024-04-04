/*
Copyright © 2024 Zachary Tarantino-Woolson <zachary.tarantino@sailpoint.com>
*/
package root

import (
	"nerm/cmd/advanced_search"
	"nerm/cmd/environment"
	"nerm/cmd/health_check"
	"nerm/cmd/identity_proofing"
	"nerm/cmd/profiles"
	"nerm/cmd/workflow_sessions"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "nerm",
	Short: "CLI Tool to make API requests and generate files asd",
	Long:  `CLI Tool to make API requests and generate files asdd`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func NewRootCommand() *cobra.Command {
	root := &cobra.Command{
		Use:          "nerm",
		Long:         "starts the use of the CLi",
		Example:      "nerm",
		SilenceUsage: true,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	root.AddCommand(
		environment.NewEnvironmentCommand(),
		health_check.NewHealthCheckCommand(),
		profiles.NewProfilesCommand(),
		workflow_sessions.NewWorkflowSessionsCommand(),
		identity_proofing.NewIdentityProofingCommand(),
		advanced_search.NewAdvancedSearchCommand(),
	)

	return root
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.nerm.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
