package commands

import (
	"bytes"
	"fmt"
	"entangle"
	"github.com/mitchellh/cli"
)

// Version command.
type VersionCommand struct {
	Ui cli.Ui
}

func (c *VersionCommand) Run(_ []string) int {
	var versionString bytes.Buffer

	fmt.Fprintf(&versionString, "Entangle v%s", entangle.VersionNumber)
	if entangle.VersionSuffix != "" {
		fmt.Fprintf(&versionString, "-%s", entangle.VersionSuffix)

		if entangle.GitCommit != "" {
			fmt.Fprintf(&versionString, " (%s)", entangle.GitCommit)
		}
	}

	c.Ui.Output(versionString.String())

	return 0
}

func (c *VersionCommand) Synopsis() string {
	return "Show Entangle version"
}

func (c *VersionCommand) Help() string {
	return ""
}
