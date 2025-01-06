package shell

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// connectToGKE establishes connection to a GKE cluster using gcloud command
func connectToGKE(config ClusterConfig) error {
	cmd := exec.Command("gcloud", "container", "clusters", "get-credentials",
		config.cluster, "--zone", config.zone, "--project", config.project)
	return cmd.Run()
}

// getPods retrieves the list of pods in the specified namespace
func getPods(namespace string) ([]string, error) {
	cmd := exec.Command("kubectl", "get", "pods", "-n", namespace, "--no-headers")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

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
func selectPod(pods []string) (string, error) {
	fmt.Printf("\n%sAvailable pods:%s\n", colorYellow, colorReset)
	for i, pod := range pods {
		fmt.Printf("%d. %s\n", i+1, pod)
	}

	var choice int
	fmt.Printf("\nSelect pod (1-%d): ", len(pods))
	fmt.Scan(&choice)
	if choice < 1 || choice > len(pods) {
		return "", fmt.Errorf("invalid pod selection")
	}
	return pods[choice-1], nil
}

// connectToPod establishes an interactive shell connection to the selected pod
func connectToPod(pod, namespace string) error {
	cmd := exec.Command("kubectl", "exec", "-it", pod, "-n", namespace, "--", "/bin/sh")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// readConfigurations reads and parses the cluster configuration file
func readConfigurations(filePath string) ([]ClusterConfig, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open configuration file: %v", err)
	}
	defer file.Close()

	var configs []ClusterConfig
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Split(line, "|")
		if len(parts) != 5 {
			return nil, fmt.Errorf("invalid configuration line format: %s", line)
		}

		config := ClusterConfig{
			env:       strings.TrimSpace(parts[0]),
			project:   strings.TrimSpace(parts[1]),
			cluster:   strings.TrimSpace(parts[2]),
			zone:      strings.TrimSpace(parts[3]),
			namespace: strings.TrimSpace(parts[4]),
		}

		if config.env == "" || config.project == "" || config.cluster == "" ||
			config.zone == "" || config.namespace == "" {
			return nil, fmt.Errorf("missing required fields in line: %s", line)
		}

		configs = append(configs, config)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading configuration file: %v", err)
	}

	if len(configs) == 0 {
		return nil, fmt.Errorf("no valid configurations found")
	}

	return configs, nil
}
