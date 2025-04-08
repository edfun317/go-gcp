package podshell

import (
	"fmt"
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
			cmdType:     AdjustCPU,
			description: "Adjust pod CPU resources",
			action:      a.adjustPodCPU,
		},
		{
			cmdType:     AdjustMemory,
			description: "Adjust pod memory resources",
			action:      a.adjustPodMemory,
		},
		{
			cmdType:     ScaleDeployment,
			description: "Scale deployment replicas",
			action:      a.scaleDeployment,
		},
		{
			cmdType:     PortForward,
			description: "Port forward service to localhost",
			action:      a.portForward,
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

// adjustPodCPU modifies the CPU resource requests/limits for a selected pod.
func (a *AccessPods) adjustPodCPU(namespace string) error {
	pods, err := getPods(namespace)
	if err != nil {
		return err
	}
	selectedPod, err := selectPod(pods)
	if err != nil {
		return err
	}

	// Get current resource values
	cmd := exec.Command("kubectl", "get", "pod", selectedPod, "-n", namespace, "-o", "jsonpath={.spec.containers[0].resources}")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	// Prompt for new CPU value
	fmt.Print("\nEnter new CPU value (e.g., '500m' for 500 millicores or '2' for 2 cores): ")
	var cpuValue string
	fmt.Scanln(&cpuValue)

	// Apply the new CPU value
	patchStr := fmt.Sprintf(`{"spec":{"containers":[{"name":"*","resources":{"requests":{"cpu":"%s"},"limits":{"cpu":"%s"}}}]}}`, cpuValue, cpuValue)
	cmd = exec.Command("kubectl", "patch", "pod", selectedPod, "-n", namespace, "--type=strategic", "-p", patchStr)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// adjustPodMemory modifies the memory resource requests/limits for a selected pod.
func (a *AccessPods) adjustPodMemory(namespace string) error {
	pods, err := getPods(namespace)
	if err != nil {
		return err
	}
	selectedPod, err := selectPod(pods)
	if err != nil {
		return err
	}

	// Get current resource values
	cmd := exec.Command("kubectl", "get", "pod", selectedPod, "-n", namespace, "-o", "jsonpath={.spec.containers[0].resources}")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	// Prompt for new memory value
	fmt.Print("\nEnter new memory value (e.g., '512Mi' or '2Gi'): ")
	var memValue string
	fmt.Scanln(&memValue)

	// Apply the new memory value
	patchStr := fmt.Sprintf(`{"spec":{"containers":[{"name":"*","resources":{"requests":{"memory":"%s"},"limits":{"memory":"%s"}}}]}}`, memValue, memValue)
	cmd = exec.Command("kubectl", "patch", "pod", selectedPod, "-n", namespace, "--type=strategic", "-p", patchStr)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// scaleDeployment modifies the number of replicas for a deployment.
func (a *AccessPods) scaleDeployment(namespace string) error {
	// Get list of deployments
	cmd := exec.Command("kubectl", "get", "deployments", "-n", namespace)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	// Get deployment name from user
	fmt.Print("\nEnter deployment name: ")
	var deploymentName string
	fmt.Scanln(&deploymentName)

	// Get current replicas
	cmd = exec.Command("kubectl", "get", "deployment", deploymentName, "-n", namespace, "-o", "jsonpath={.spec.replicas}")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fmt.Printf("\nCurrent replicas: ")
	if err := cmd.Run(); err != nil {
		return err
	}

	// Get new replica count from user
	fmt.Print("\nEnter new number of replicas: ")
	var replicaCount string
	fmt.Scanln(&replicaCount)

	// Scale the deployment
	cmd = exec.Command("kubectl", "scale", "deployment", deploymentName, "-n", namespace, fmt.Sprintf("--replicas=%s", replicaCount))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// portForward forwards a local port to a service in the GKE cluster.
func (a *AccessPods) portForward(namespace string) error {
	// Get list of services
	cmd := exec.Command("kubectl", "get", "services", "-n", namespace)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	// Get service name from user
	fmt.Print("\nEnter service name: ")
	var serviceName string
	fmt.Scanln(&serviceName)

	// Get target port from user
	cmd = exec.Command("kubectl", "get", "service", serviceName, "-n", namespace, "-o", "jsonpath={.spec.ports[*].port}")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fmt.Printf("\nAvailable ports: ")
	if err := cmd.Run(); err != nil {
		return err
	}

	fmt.Print("\nEnter target port: ")
	var targetPort string
	fmt.Scanln(&targetPort)

	// Get local port from user
	fmt.Print("\nEnter local port to forward to: ")
	var localPort string
	fmt.Scanln(&localPort)

	// Start port forwarding
	fmt.Printf("\nStarting port forward from localhost:%s to service %s:%s\n", localPort, serviceName, targetPort)
	cmd = exec.Command("kubectl", "port-forward", fmt.Sprintf("service/%s", serviceName), fmt.Sprintf("%s:%s", localPort, targetPort), "-n", namespace)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
