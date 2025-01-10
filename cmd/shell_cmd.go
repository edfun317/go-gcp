package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "Execute shell commands to access related functionalities",
	Long:  "Execute shell commands to access and manage GCP-related functionalities",
	Run: func(cmd *cobra.Command, args []string) {
		// Check if list flag is set
		if list, _ := cmd.Flags().GetBool("list"); list {
			fmt.Println("Available GCP commands:")
			for _, command := range gcpCommands {
				fmt.Printf("  %s\n", command)
			}
			return
		}

		filePath, _ := cmd.Flags().GetString("file")

		// Check if file exists
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			fmt.Printf("Error: file '%s' does not exist\n", filePath)
			return
		}

		// Process the file
		fmt.Printf("Processing file: %s\n", filePath)

		// Handle additional arguments
		if len(args) > 0 {
			fmt.Printf("Additional arguments: %v\n", args)
		}
	},
}
