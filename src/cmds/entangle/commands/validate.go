package commands

import (
	"entangle/errors"
	"entangle/parser"
	"entangle/source"
	"fmt"
	"github.com/mitchellh/cli"
	"os"
)

// Validate command.
//
// Validates a definition file by parsing it and outputting any errors.
type ValidateCommand struct {
	Ui cli.Ui
}

func (c *ValidateCommand) Help() string {
	return `Usage: entangle validate <path>

  Validate an Entangle definition file.`
}

func (c *ValidateCommand) Run(args []string) int {
	// Parse the path from the arguments.
	if len(args) != 1 {
		if len(args) == 0 {
			c.Ui.Error("A definition file path is required.")
		} else {
			c.Ui.Error("Only one definition file path may be supplied.")
		}
		c.Ui.Error("")
		c.Ui.Error(c.Help())
		return 1
	}

	path := args[0]

	// Open the file for reading.
	f, err := os.Open(path)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Failed to open %s: %v", path, err))
		return 1
	}
	defer f.Close()

	// Read the file into a source.
	src, err := source.FromReader(f, path)

	if err != nil {
		c.Ui.Error(fmt.Sprintf("Failed to read %s: %v", path, err))
		return 1
	}

	// Parse the file.
	_, err = parser.Parse(src)

	if parseErr, ok := err.(errors.ParseError); ok {
		parser.PrintError(parseErr)
		return 1
	} else if err != nil {
		c.Ui.Error(fmt.Sprintf("Failed to parse %s: %v", path, err))
		return 1
	}

	return 0
}

func (c *ValidateCommand) Synopsis() string {
	return "Validate a definition file."
}
