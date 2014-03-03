package commands

import (
	"entangle/declarations"
	"entangle/errors"
	"entangle/parser"
	"entangle/source"
	"flag"
	"fmt"
	"github.com/mitchellh/cli"
	"os"
	"strings"
)

// Target language definition.
type targetLanguageDefinition struct {
	// Synopsis.
	Synopsis string

	// Construct flag set.
	FlagSet func() (*flag.FlagSet, interface {})

	// Generate output.
	Generate func(interfaceDecl *declarations.Interface, outputPath string, options interface{}) error
}

var targetLanguageMapping = map[string]targetLanguageDefinition {
	"go": {
		Synopsis: "Go",
		FlagSet: goTargetFlagSet,
		Generate: goTargetGenerate,
	},
}

// Generate command.
type GenerateCommand struct {
	Ui cli.Ui
}

func (c *GenerateCommand) Help() string {
	targetLanguages := make([]string, 0, len(targetLanguageMapping))
	languageOptions := make([]string, 0, len(targetLanguageMapping))

	for ident, def := range targetLanguageMapping {
		// Add the language to the list of target languages.
		targetLanguages = append(targetLanguages, fmt.Sprintf("  %s%s%s", ident, strings.Repeat(" ", 10 - len(ident)), def.Synopsis))

		// Determine the maximum length of flag names.
		flagSet, _ := def.FlagSet()

		usageListElements := []usageListElement {}

		flagSet.VisitAll(func(f *flag.Flag) {
			var name string
			if len(f.Name) == 1 {
				name = fmt.Sprintf("-%s", f.Name)
			} else {
				name = fmt.Sprintf("--%s", f.Name)
			}

			usageListElements = append(usageListElements, usageListElement {
				Name: name,
				Synopsis: f.Usage,
			})
		})

		if len(usageListElements) == 0 {
			continue
		}

		languageOptions = append(languageOptions, fmt.Sprintf(`Options for %s (%s):

%s`, def.Synopsis, ident, usageList(usageListElements)))
	}

	return fmt.Sprintf(`Usage: entangle generate <language> [options] <definition path> <output path>

  Generate an implementation from an Entangle definition file. Options depend
  on the target language.

Arguments:

  <language>            Target language to generate an implementation for.
                        See the list of available target languages below.
  <definition path>     Path of the definition file.
  <output path>         Output path of the generated code.

Target languages:

%s

%s`, strings.Join(targetLanguages, "\n"), strings.Join(languageOptions, "\n\n"))
}

func (c *GenerateCommand) Run(args []string) int {
	// Parse the target language from the arguments.
	if len(args) < 3 {
		c.Ui.Error("A target language and definition and output paths must be supplied.")
		c.Ui.Error("")
		c.Ui.Error(c.Help())
		return 1
	}

	targetLanguageIdentifier := args[0]

	targetLanguage, targetLanguageValid := targetLanguageMapping[targetLanguageIdentifier]
	if !targetLanguageValid {
		c.Ui.Error(fmt.Sprintf("Invalid target language: %s", targetLanguageIdentifier))
		c.Ui.Error("")
		c.Ui.Error(c.Help())
		return 1
	}

	// Parse the options for the target language.
	flagSet, options := targetLanguage.FlagSet()
	flagSet.Usage = func() {
		c.Ui.Output("")
		c.Ui.Output(c.Help())
	}
	if err := flagSet.Parse(args[1:]); err != nil {
		return 1
	}

	// Parse the path from the arguments.
	paths := flagSet.Args()

	if len(paths) != 2 {
		if len(paths) == 0 {
			c.Ui.Error("A definition file path is required.")
		} else if len(paths) == 1 {
			c.Ui.Error("An output path is required.")
		} else {
			c.Ui.Error("Too many arguments. Only a definition file and output path may be supplied.")
		}
		c.Ui.Error("")
		c.Ui.Error(c.Help())
		return 1
	}

	path := paths[0]
	outputPath := paths[1]

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
	interfaceDecl, err := parser.Parse(src)

	if parseErr, ok := err.(errors.ParseError); ok {
		parser.PrintError(parseErr)
		return 1
	} else if err != nil {
		c.Ui.Error(fmt.Sprintf("Failed to parse %s: %v", path, err))
		return 1
	}

	// Make sure the output directory exists.
	var outputPathStat os.FileInfo
	var statErr error

	if outputPathStat, statErr = os.Stat(outputPath); statErr != nil {
		if !os.IsNotExist(statErr) {
			c.Ui.Error(fmt.Sprintf("Failed to determine status of output directory '%s': %v", outputPath, statErr))
			return 1
		}

		if statErr = os.MkdirAll(outputPath, 0777); statErr != nil {
			c.Ui.Error(fmt.Sprintf("Failed to create output directory '%s': %v", outputPath, statErr))
			return 1
		}
	} else if !outputPathStat.Mode().IsDir() {
		c.Ui.Error(fmt.Sprintf("Output path is not a directory: %s", outputPath))
		return 1
	}

	// Perform the generation.
	if err = targetLanguage.Generate(interfaceDecl, outputPath, options); err != nil {
		c.Ui.Error(fmt.Sprintf("Failed to generate implementation: %v", err))
		return 1
	}

	return 0
}

func (c *GenerateCommand) Synopsis() string {
	return "Generate a target from a definition file."
}
