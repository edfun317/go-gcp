package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "myapp",
	Short: "A command-line tool for executing GCP-related shell commands",
	Long: `A command-line tool for executing GCP-related shell commands 
and managing cloud resources efficiently.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Use --help to see available arguments")
	},
}

// Define available GCP shell commands
var gcpCommands = []string{
	"gke - Manage Google Kubernetes Engine clusters",
	//"gcloud - Manage Google Cloud resources",
	//	"gsutil - Access Google Cloud Storage",
}

// Execute adds all child commands to the root command and sets flags appropriately
func Execute() {
	// Add flags to shell command
	shellCmd.Flags().StringP("file", "f", "", "File path to use")
	shellCmd.Flags().BoolP("list", "l", false, "List available GCP commands")

	// Mark file flag as required
	shellCmd.MarkFlagRequired("file")

	// Add shell command to root command
	rootCmd.AddCommand(shellCmd)

	// Execute root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
