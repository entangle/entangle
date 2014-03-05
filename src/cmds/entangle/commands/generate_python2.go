package commands

import (
	"entangle/declarations"
	"entangle/generators"
	"entangle/generators/python2"
	"flag"
)

// Python 2 target options.
type python2TargetOptions struct {}

// Python 2 target flag set.
func python2TargetFlagSet() (s *flag.FlagSet, options interface {}) {
	o := &python2TargetOptions {}
	s = flag.NewFlagSet("go", flag.ExitOnError)
	options = o
	return
}

// Python 2 target generation.
func python2TargetGenerate(interfaceDecl *declarations.Interface, outputPath string, o interface{}) (err error) {
	// Initialize the generator.
	var generator generators.Generator
	if generator, err = python2.NewGenerator(&python2.Options{}); err != nil {
		return
	}

	// Generate.
	return generator.Generate(interfaceDecl, outputPath)
}
