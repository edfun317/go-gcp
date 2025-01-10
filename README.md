# go-gcp

A command-line tool for executing GCP-related shell commands and managing cloud resources efficiently.

## Features

- Execute shell commands to access GCP-related functionalities
- Built-in help system for command discovery
- File-based operations support
- List available GCP commands

## Installation

```bash
# Clone the repository
git clone https://github.com/edfun317/go-gcp.git

# Navigate to the project directory
cd go-gcp

# Install dependencies
go mod tidy
```

## Usage

### Basic Commands

1. Show available commands:
```bash
go run .
```

2. Display help information:
```bash
go run . --help
```

### Shell Command

The `shell` command allows you to execute GCP-related shell commands:

```bash
# List available GCP commands
go run . shell -l

# Execute with a file
go run . shell -f /path/to/file
```

Available flags:
- `-f, --file`: Specify the file path (required)
- `-l, --list`: List available GCP commands
- `-h, --help`: Help for shell command

### Available GCP Commands

The following GCP commands are supported:
- `gke`: Manage Google Kubernetes Engine clusters
- `gcloud`: Manage Google Cloud resources
- `gsutil`: Access Google Cloud Storage

## Project Structure

```
go-gcp/
├── cmd/
│   ├── cmd.go       # Command definitions
│   ├── main.go      # Entry point
│   └── shell_cmd.go # Shell command implementation
├── shell/
│   └── podshell/
│       ├── commands.go  # Shell commands
│       ├── execute.go   # Command execution
│       ├── types.go     # Type definitions
│       └── utils.go     # Utility functions
├── .gitignore       # Git ignore file
├── go.mod           # Go module file
├── go.sum           # Go module checksum
└── README.md        # Project documentation
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.