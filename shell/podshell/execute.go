package podshell

import (
	"fmt"
	"os"
	"strings"
)

// Execute handles the main flow of connecting to a GKE cluster and executing commands
// This is the core function that orchestrates the entire pod access workflow
func (a *AccessPods) Execute() {
	// Step 1: Configuration Loading
	// Read and parse the cluster configurations from the specified file path
	// The configurations contain environment, project, cluster, zone, and namespace information
	configs, err := readConfigurations(a.FilePath)
	if err != nil {
		fmt.Printf("%sError: %v%s\n", colorRed, err, colorReset)
		os.Exit(1)
	}

	// Step 2: Environment Selection
	// Display available environments to the user for selection
	// Each environment represents a different GKE cluster setup (e.g., dev, staging, prod)
	fmt.Printf("%sAvailable environments:%s\n", colorYellow, colorReset)
	for i, config := range configs {
		fmt.Printf("%d. %s\n", i+1, config.env)
	}

	// Step 3: User Input Processing
	// Get and validate user's environment selection
	// Ensure the selection is within valid range
	var choice int
	fmt.Printf("Select environment (1-%d): ", len(configs))
	fmt.Scan(&choice)
	if choice < 1 || choice > len(configs) {
		fmt.Printf("%sInvalid selection%s\n", colorRed, colorReset)
		os.Exit(1)
	}
	selectedConfig := configs[choice-1]

	// Step 4: Configuration Display
	// Show the selected configuration details to the user for verification
	// This includes all relevant GKE cluster connection parameters
	fmt.Printf("\n%sSelected Configuration:%s\n", colorGreen, colorReset)
	fmt.Printf("Environment: %s\n", selectedConfig.env)
	fmt.Printf("Project: %s\n", selectedConfig.project)
	fmt.Printf("Cluster: %s\n", selectedConfig.cluster)
	fmt.Printf("Zone: %s\n", selectedConfig.zone)
	fmt.Printf("Namespace: %s\n", selectedConfig.namespace)

	// Step 5: User Confirmation
	// Ask for user confirmation before proceeding with cluster connection
	// This prevents accidental connections to wrong environments
	var confirm string
	fmt.Print("Continue? (y/n): ")
	fmt.Scan(&confirm)
	if strings.ToLower(confirm) != "y" {
		fmt.Println("Operation cancelled")
		os.Exit(0)
	}

	// Step 6: GKE Cluster Connection
	// Establish connection to the selected GKE cluster using gcloud command
	// This sets up the kubectl context for subsequent commands
	if err := connectToGKE(selectedConfig); err != nil {
		fmt.Printf("%sError: %v%s\n", colorRed, err, colorReset)
		os.Exit(1)
	}

	// Step 7: Command Setup
	// Define the order of available commands in the interactive menu
	// This determines the sequence and options presented to the user
	commandOrder := []CommandType{
		ShowPods,    // Display all pods in the namespace
		ConnectPod,  // Connect to a specific pod
		ShowLogs,    // View pod logs
		DescribePod, // Get detailed pod information
		Exit,        // Exit the program
	}

	// Step 8: Command Loop
	// Enter an infinite loop for command execution until user exits
	// This is the main interactive portion of the program
	for {
		// Display available commands with numbering
		fmt.Printf("\n%sAvailable commands:%s\n", colorYellow, colorReset)
		for i, cmdType := range commandOrder {
			if cmd, ok := a.Commands[cmdType]; ok {
				fmt.Printf("%d. %s\n", i+1, cmd.Description)
			}
		}

		// Get and validate user's command selection
		var cmdChoice int
		fmt.Printf("\nSelect command (1-%d): ", len(a.Commands))
		fmt.Scan(&cmdChoice)
		if cmdChoice < 1 || cmdChoice > len(commandOrder) {
			fmt.Printf("%sInvalid command%s\n", colorRed, colorReset)
			continue
		}

		// Execute the selected command
		cmdType := commandOrder[cmdChoice-1]
		if cmd, ok := a.Commands[cmdType]; ok {
			// Execute the command's action with the selected namespace
			// Handle any errors that occur during execution
			if err := cmd.Action(selectedConfig.namespace); err != nil {
				fmt.Printf("%sError: %v%s\n", colorRed, err, colorReset)
			}
			// Exit the program if the Exit command was selected
			if cmdType == Exit {
				return
			}
		}
	}
}
