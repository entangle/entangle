package commands

import (
	"entangle/declarations"
	"entangle/generators"
	"entangle/generators/golang"
	"flag"
)

// Go target options.
type goTargetOptions struct {
	// Package name override.
	Package string
}

// Go target flag set.
func goTargetFlagSet() (s *flag.FlagSet, options interface {}) {
	o := &goTargetOptions {}
	s = flag.NewFlagSet("go", flag.ExitOnError)
	s.StringVar(&o.Package, "package", "", "Override generated package name. If not provided, the definition name will be used as package name.")
	options = o
	return
}

// Go target generation.
func goTargetGenerate(interfaceDecl *declarations.Interface, outputPath string, o interface{}) (err error) {
	options := o.(*goTargetOptions)

	// Initialize the generator.
	var generator generators.Generator
	if generator, err = golang.NewGenerator(&golang.Options{
		Package: options.Package,
	}); err != nil {
		return
	}

	// Generate.
	return generator.Generate(interfaceDecl, outputPath)
}
