package main

import (
	"fmt"
	"os"
	"github.com/mitchellh/cli"
)

func main() {
	os.Exit(mainImpl())
}

func mainImpl() int {
	args := os.Args[1:]

	// Fast path to version argument as per original uses of CLI.
	for _, arg := range args {
		if arg == "-v" || arg == "--version" {
			newArgs := make([]string, len(args)+1)
			newArgs[0] = "version"
			copy(newArgs[1:], args)
			args = newArgs
			break
		}
	}

	// Set up and execute the CLI.
	cli := &cli.CLI{
		Args:     args,
		Commands: Commands,
		HelpFunc: cli.BasicHelpFunc("entangle"),
	}

	exitCode, err := cli.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing CLI: %s\n", err.Error())
		return 1
	}

	return exitCode
}
