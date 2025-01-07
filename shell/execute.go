package shell

import (
	"fmt"
	"os"
	"strings"
)

// Execute handles the main flow of connecting to a GKE cluster and executing commands
func (a *AccessPods) Execute() {

	// Initialize by reading cluster configurations
	configs, err := readConfigurations(a.FilePath)
	if err != nil {
		fmt.Printf("%sError: %v%s\n", colorRed, err, colorReset)
		os.Exit(1)
	}

	// Display environment options to user
	fmt.Printf("%sAvailable environments:%s\n", colorYellow, colorReset)
	for i, config := range configs {
		fmt.Printf("%d. %s\n", i+1, config.env)
	}

	// Handle environment selection
	var choice int
	fmt.Printf("Select environment (1-%d): ", len(configs))
	fmt.Scan(&choice)
	if choice < 1 || choice > len(configs) {
		fmt.Printf("%sInvalid selection%s\n", colorRed, colorReset)
		os.Exit(1)
	}

	selectedConfig := configs[choice-1]

	// Show selected configuration details
	fmt.Printf("\n%sSelected Configuration:%s\n", colorGreen, colorReset)
	fmt.Printf("Environment: %s\n", selectedConfig.env)
	fmt.Printf("Project: %s\n", selectedConfig.project)
	fmt.Printf("Cluster: %s\n", selectedConfig.cluster)
	fmt.Printf("Zone: %s\n", selectedConfig.zone)
	fmt.Printf("Namespace: %s\n", selectedConfig.namespace)

	// Confirmation before proceeding
	var confirm string
	fmt.Print("Continue? (y/n): ")
	fmt.Scan(&confirm)
	if strings.ToLower(confirm) != "y" {
		fmt.Println("Operation cancelled")
		os.Exit(0)
	}

	// Connect to the selected GKE cluster
	if err := connectToGKE(selectedConfig); err != nil {
		fmt.Printf("%sError: %v%s\n", colorRed, err, colorReset)
		os.Exit(1)
	}

	commandOrder := []CommandType{
		ShowPods,
		ConnectPod,
		ShowLogs,
		DescribePod,
		Exit,
	}

	// Enter the command loop
	for {
		fmt.Printf("\n%sAvailable commands:%s\n", colorYellow, colorReset)
		for i, cmdType := range commandOrder {
			if cmd, ok := a.Commands[cmdType]; ok {
				fmt.Printf("%d. %s\n", i+1, cmd.Description)
			}
		}

		var cmdChoice int
		fmt.Printf("\nSelect command (1-%d): ", len(a.Commands))
		fmt.Scan(&cmdChoice)

		if cmdChoice < 1 || cmdChoice > len(commandOrder) {
			fmt.Printf("%sInvalid command%s\n", colorRed, colorReset)
			continue
		}

		cmdType := commandOrder[cmdChoice-1]
		if cmd, ok := a.Commands[cmdType]; ok {
			if err := cmd.Action(selectedConfig.namespace); err != nil {
				fmt.Printf("%sError: %v%s\n", colorRed, err, colorReset)
			}
			if cmdType == Exit {
				return
			}
		}
	}
}
