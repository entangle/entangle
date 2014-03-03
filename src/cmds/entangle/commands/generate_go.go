package commands

import (
	"entangle/declarations"
	"entangle/generators"
	"entangle/generators/golang"
	"flag"
)

// Go target options.
type goTargetOptions struct {}

// Go target flag set.
func goTargetFlagSet() (s *flag.FlagSet, options interface {}) {
	o := &goTargetOptions {}
	s = flag.NewFlagSet("go", flag.ExitOnError)
	options = o
	return
}

// Go target generation.
func goTargetGenerate(interfaceDecl *declarations.Interface, outputPath string, o interface{}) (err error) {
	// Initialize the generator.
	var generator generators.Generator
	if generator, err = golang.NewGenerator(&golang.Options{}); err != nil {
		return
	}

	// Generate.
	return generator.Generate(interfaceDecl, outputPath)
}
