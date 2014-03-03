package main

import (
	"cmds/entangle/commands"
	"github.com/mitchellh/cli"
	"os"
)

var Commands map[string]cli.CommandFactory

func init() {
	ui := &cli.BasicUi{
		Writer: os.Stdout,
	}

	Commands = map[string]cli.CommandFactory{
		"version": func() (cli.Command, error) {
			return &commands.VersionCommand{
				Ui: ui,
			}, nil
		},
		"validate": func() (cli.Command, error) {
			return &commands.ValidateCommand{
				Ui: ui,
			}, nil
		},
		"generate": func() (cli.Command, error) {
			return &commands.GenerateCommand{
				Ui: ui,
			}, nil
		},
	}
}
