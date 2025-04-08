package podshell

import (
	"fmt"
	"os"
	"strings"
)

// Execute handles the main flow of connecting to a GKE cluster and executing commands
func (a *AccessPods) Execute() {
	// Load and select configuration
	selectedConfig, err := a.setupClusterConfig()
	if err != nil {
		a.handleError("Configuration setup failed", err)
		os.Exit(1)
	}

	// Connect to GKE cluster
	if err := connectToGKE(selectedConfig); err != nil {
		a.handleError("GKE connection failed", err)
		os.Exit(1)
	}

	// Start command loop
	a.commandLoop(selectedConfig)
}

// setupClusterConfig handles configuration loading and selection
func (a *AccessPods) setupClusterConfig() (ClusterConfig, error) {
	// Read configurations
	configs, err := readConfigurations(a.FilePath)
	if err != nil {
		return ClusterConfig{}, err
	}

	// Display environments
	fmt.Printf("%sAvailable environments:%s\n", colorYellow, colorReset)
	for i, config := range configs {
		fmt.Printf("%d. %s\n", i+1, config.env)
	}

	// Get user selection
	choice := a.getUserInput(fmt.Sprintf("Select environment (1-%d): ", len(configs)))
	if choice < 1 || choice > len(configs) {
		return ClusterConfig{}, fmt.Errorf("invalid environment selection")
	}

	selectedConfig := configs[choice-1]
	if err := a.confirmConfiguration(selectedConfig); err != nil {
		return ClusterConfig{}, err
	}

	return selectedConfig, nil
}

// confirmConfiguration displays and confirms the selected configuration
func (a *AccessPods) confirmConfiguration(config ClusterConfig) error {
	fmt.Printf("\n%sSelected Configuration:%s\n", colorGreen, colorReset)
	fmt.Printf("Environment: %s\n", config.env)
	fmt.Printf("Project: %s\n", config.project)
	fmt.Printf("Cluster: %s\n", config.cluster)
	fmt.Printf("Zone: %s\n", config.zone)
	fmt.Printf("Namespace: %s\n", config.namespace)

	if !a.getUserConfirmation("Continue? (y/n): ") {
		return fmt.Errorf("operation cancelled by user")
	}
	return nil
}

func (a *AccessPods) commandLoop(config ClusterConfig) {
	if a.Commands == nil {
		a.handleError("Command initialization", fmt.Errorf("commands not initialized"))
		return
	}

	commandOrder := []CommandType{
		ShowPods,
		ConnectPod,
		ShowLogs,
		DescribePod,
		ShowEnv,
		AdjustCPU,
		AdjustMemory,
		ScaleDeployment,
		PortForward,
		Exit,
	}

	for {
		// 顯示命令列表
		fmt.Printf("\n%sAvailable commands:%s\n", colorYellow, colorReset)
		for i, cmdType := range commandOrder {
			cmd, ok := a.Commands[cmdType]
			if !ok {
				a.handleError("Command lookup", fmt.Errorf("command %v not registered", cmdType))
				continue
			}
			fmt.Printf("%d. %s\n", i+1, cmd.Description)
		}

		// 獲取用戶輸入
		choice := a.getUserInput(fmt.Sprintf("\nSelect command (1-%d): ", len(commandOrder)))
		if choice < 1 || choice > len(commandOrder) {
			fmt.Printf("%sInvalid command%s\n", colorRed, colorReset)
			continue
		}

		// 執行命令
		cmdType := commandOrder[choice-1]
		cmd, ok := a.Commands[cmdType]
		if !ok {
			a.handleError("Command execution", fmt.Errorf("command %v not found", cmdType))
			continue
		}

		if cmd.Action == nil {
			a.handleError("Command execution", fmt.Errorf("action not defined for command %v", cmdType))
			continue
		}

		if err := cmd.Action(config.namespace); err != nil {
			a.handleError("Command execution failed", err)
		}

		if cmdType == Exit {
			return
		}
	}
}

// Helper functions

func (a *AccessPods) getUserInput(prompt string) int {

	var choice int
	fmt.Print(prompt)
	fmt.Scan(&choice)
	return choice
}

func (a *AccessPods) getUserConfirmation(prompt string) bool {

	var confirm string
	fmt.Print(prompt)
	fmt.Scan(&confirm)
	return strings.ToLower(confirm) == "y"
}

func (a *AccessPods) handleError(context string, err error) {
	fmt.Printf("%sError: %s: %v%s\n", colorRed, context, err, colorReset)
}
