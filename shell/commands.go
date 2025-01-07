package shell

import (
	"os"
	"os/exec"
)

// NewAccessPods creates and initializes a new AccessPods instance
func NewAccessPods(filePath string) *AccessPods {
	a := &AccessPods{
		FilePath: filePath,
		Commands: make(map[CommandType]ShellCommand),
	}

	a.registerCommands()
	return a
}

// registerCommands initializes all available commands
func (a *AccessPods) registerCommands() {

	// set the command and sort
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
			cmdType:     Exit,
			description: "Exit program",
			action:      func(namespace string) error { os.Exit(0); return nil },
		},
	}

	a.Commands = make(map[CommandType]ShellCommand)

	for _, cmd := range commandOrder {
		a.Commands[cmd.cmdType] = ShellCommand{
			Type:        cmd.cmdType,
			Description: cmd.description,
			Action:      cmd.action,
		}
	}
}

// Command implementation functions
func (a *AccessPods) listPods(namespace string) error {
	cmd := exec.Command("kubectl", "get", "pods", "-n", namespace)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

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
