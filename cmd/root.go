package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "main",
	Short: "Start a server",
}

func Execute() {
	rootCmd.AddCommand(restApiServiceCmd)

	InitFlags()
	rootCmd.Execute()
}

func InitFlags() {
	restApiServiceCmd.PersistentFlags().Bool("start", false, "Command to start service with default port 8080")

}
