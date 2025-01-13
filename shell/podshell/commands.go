package podshell

import (
	"os"
	"os/exec"
)

// NewAccessPods creates and initializes a new AccessPods instance to manage pod operations.
// It takes a filePath parameter specifying the location of cluster configuration file.
// Returns a pointer to the initialized AccessPods struct with registered commands.
func NewAccessPods(filePath string) *AccessPods {
	a := &AccessPods{
		FilePath: filePath,
		Commands: make(map[CommandType]ShellCommand),
	}
	a.registerCommands()
	return a
}

// registerCommands initializes all available pod management commands.
// This includes commands for listing, connecting, viewing logs, describing pods,
// showing environment variables, and program exit functionality.
func (a *AccessPods) registerCommands() {
	// Define command order and properties
	commandOrder := []struct {
		cmdType     CommandType
		description string
		action      func(namespace string) error
	}{
		{
			cmdType:     ShowPods,
			description: "List all pods",
			action:      a.listPods,
		},
		{
			cmdType:     ConnectPod,
			description: "Connect to a pod",
			action:      a.connectToPodShell,
		},
		{
			cmdType:     ShowLogs,
			description: "Show pod logs",
			action:      a.showPodLogs,
		},
		{
			cmdType:     DescribePod,
			description: "Describe pod",
			action:      a.describePod,
		},
		{
			cmdType:     ShowEnv,
			description: "Show environment variables",
			action:      a.showPodEnv,
		},
		{
			cmdType:     Exit,
			description: "Exit program",
			action:      func(namespace string) error { os.Exit(0); return nil },
		},
	}

	// Initialize command map
	a.Commands = make(map[CommandType]ShellCommand)

	// Register each command with its properties
	for _, cmd := range commandOrder {
		a.Commands[cmd.cmdType] = ShellCommand{
			Type:        cmd.cmdType,
			Description: cmd.description,
			Action:      cmd.action,
		}
	}
}

// listPods executes kubectl command to list all pods in the specified namespace.
// Displays pod information directly to stdout.
func (a *AccessPods) listPods(namespace string) error {
	cmd := exec.Command("kubectl", "get", "pods", "-n", namespace)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// connectToPodShell establishes an interactive shell connection to a selected pod.
// First retrieves available pods, then lets user select one before connecting.
func (a *AccessPods) connectToPodShell(namespace string) error {
	pods, err := getPods(namespace)
	if err != nil {
		return err
	}
	selectedPod, err := selectPod(pods)
	if err != nil {
		return err
	}
	return connectToPod(selectedPod, namespace)
}

// showPodLogs retrieves and displays logs from a selected pod.
// User can select which pod's logs to view from available pods.
func (a *AccessPods) showPodLogs(namespace string) error {
	pods, err := getPods(namespace)
	if err != nil {
		return err
	}
	selectedPod, err := selectPod(pods)
	if err != nil {
		return err
	}
	cmd := exec.Command("kubectl", "logs", selectedPod, "-n", namespace)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// describePod shows detailed information about a selected pod.
// Executes kubectl describe command on the chosen pod.
func (a *AccessPods) describePod(namespace string) error {
	pods, err := getPods(namespace)
	if err != nil {
		return err
	}
	selectedPod, err := selectPod(pods)
	if err != nil {
		return err
	}
	cmd := exec.Command("kubectl", "describe", "pod", selectedPod, "-n", namespace)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// showPodEnv displays environment variables for a selected pod.
// Executes 'env' command inside the pod to list all environment variables.
func (a *AccessPods) showPodEnv(namespace string) error {
	// Get and select target pod
	pods, err := getPods(namespace)
	if err != nil {
		return err
	}
	selectedPod, err := selectPod(pods)
	if err != nil {
		return err
	}

	// Execute command to retrieve environment variables
	cmd := exec.Command("kubectl", "exec", selectedPod, "-n", namespace, "--", "env")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
