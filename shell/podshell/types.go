package podshell

// CommandType 定義可用的命令類型
type CommandType int

const (
	ShowPods CommandType = iota
	ConnectPod
	ShowLogs
	DescribePod
	ShowEnv
	Exit
)

// ClusterConfig holds the configuration for a GKE cluster
// format example:
// env | project | cluster | zone | namespace
// dev|project-dev|cluster-dev|us-central1-a|default
// staging|project-stg|cluster-stg|us-central1-b|staging
// prod|project-prod|cluster-prod|us-central1-c|production
type ClusterConfig struct {
	env       string // Environment name (e.g., dev, staging, prod)
	project   string // GCP project ID
	cluster   string // GKE cluster name
	zone      string // GCP zone where the cluster is located
	namespace string // Kubernetes namespace
}

// ShellCommand represents a single command with its action
type ShellCommand struct {
	Type        CommandType
	Description string
	Action      func(namespace string) error
}

type DBConfig struct {
	Env     string
	Host    string
	Port    string
	Purpose string
	DBName  string
}

// AccessPods is the main structure for handling pod access
type AccessPods struct {
	FilePath string // Path to the configuration file
	Commands map[CommandType]ShellCommand
}

// ANSI color codes for terminal output formatting
const (
	colorRed    = "\033[0;31m"
	colorGreen  = "\033[0;32m"
	colorYellow = "\033[1;33m"
	colorReset  = "\033[0m"
)
