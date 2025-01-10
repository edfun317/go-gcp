package podshell

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// connectToGKE establishes connection to a GKE cluster using gcloud command
// Command: gcloud container clusters get-credentials my-cluster --zone us-central1-a --project my-project
func connectToGKE(config ClusterConfig) error {
	// Constructs and executes gcloud command to get cluster credentials
	cmd := exec.Command("gcloud", "container", "clusters", "get-credentials",
		config.cluster, "--zone", config.zone, "--project", config.project)
	return cmd.Run()
}

// getPods retrieves the list of pods in the specified namespace
// Command: kubectl get pods -n bi-rpa-pd-pnc --no-headers
func getPods(namespace string) ([]string, error) {
	// Execute kubectl command to get pod list
	cmd := exec.Command("kubectl", "get", "pods", "-n", namespace, "--no-headers")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	// Parse the command output to extract pod names
	var pods []string
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) > 0 {
			pods = append(pods, fields[0])
		}
	}
	return pods, nil
}

// selectPod displays available pods and handles pod selection
// Interactive command: User selects pod number from displayed list
func selectPod(pods []string) (string, error) {
	// Display available pods with numbering
	fmt.Printf("\n%sAvailable pods:%s\n", colorYellow, colorReset)
	for i, pod := range pods {
		fmt.Printf("%d. %s\n", i+1, pod)
	}

	// Get user input for pod selection
	var choice int
	fmt.Printf("\nSelect pod (1-%d): ", len(pods))
	fmt.Scan(&choice)

	// Validate user selection
	if choice < 1 || choice > len(pods) {
		return "", fmt.Errorf("invalid pod selection")
	}
	return pods[choice-1], nil
}

// connectToPod establishes an interactive shell connection to the selected pod
// Command: kubectl exec -it pod-name -n namespace -- /bin/sh
func connectToPod(pod, namespace string) error {
	// Setup interactive shell connection to the pod
	cmd := exec.Command("kubectl", "exec", "-it", pod, "-n", namespace, "--", "/bin/sh")

	// Connect standard input/output/error for interactive session
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// readConfigurations reads and parses the cluster configuration file
// File format: env|project|cluster|zone|namespace
// Example line: prod|my-project|my-cluster|us-central1-a|default
func readConfigurations(filePath string) ([]ClusterConfig, error) {
	// Open and defer close of configuration file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open configuration file: %v", err)
	}
	defer file.Close()

	var configs []ClusterConfig
	scanner := bufio.NewScanner(file)

	// Parse configuration file line by line
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Split and validate configuration line
		parts := strings.Split(line, "|")
		if len(parts) != 5 {
			return nil, fmt.Errorf("invalid configuration line format: %s", line)
		}

		// Create and validate config object
		config := ClusterConfig{
			env:       strings.TrimSpace(parts[0]),
			project:   strings.TrimSpace(parts[1]),
			cluster:   strings.TrimSpace(parts[2]),
			zone:      strings.TrimSpace(parts[3]),
			namespace: strings.TrimSpace(parts[4]),
		}

		// Verify all required fields are present
		if config.env == "" || config.project == "" || config.cluster == "" ||
			config.zone == "" || config.namespace == "" {
			return nil, fmt.Errorf("missing required fields in line: %s", line)
		}
		configs = append(configs, config)
	}

	// Check for scanner errors and valid configurations
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading configuration file: %v", err)
	}
	if len(configs) == 0 {
		return nil, fmt.Errorf("no valid configurations found")
	}
	return configs, nil
}
