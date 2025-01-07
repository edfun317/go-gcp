package main

import (
	"fmt"
	"os"
	"swissknife/module_go_log_library/github.com/edfun317/go-gcp/shell"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: program <config-file>")
		os.Exit(1)
	}

	access := shell.NewAccessPods(os.Args[1])
	access.Execute()
}
