package python2

import (
	"entangle/declarations"
	"entangle/generators"
	"os"
	"path/filepath"
)

// Python 2 generator.
type generator struct {
	options   *Options
}

// New generator.
func NewGenerator(options *Options) (gen generators.Generator, err error) {
	return &generator{
		options: options,
	}, nil
}

// Generate.
func (g *generator) Generate(interfaceDecl *declarations.Interface, outputPath string) (err error) {
	// Build a serialization/deserialization map for helper functions.
	serDesMap := buildSerDesMap(interfaceDecl)

	// Set up the context.
	ctx := &context{
		Interface: interfaceDecl,
		SerDesMap: serDesMap,
	}

	// Generate output files.
	for _, output := range []struct {
		Filename string
		Generator func(*context) (*SourceFile, error)
	} {
		{ "__init__.py", generateInit },
		{ "types.py", generateTypes },
		{ "clients.py", generateClients },
		{ "exceptions.py", generateExceptions },
		{ "deserialization.py", generateDeserialization },
		{ "packing.py", generatePacking },
	} {
		filePath := filepath.Join(outputPath, output.Filename)

		// Get the source file.
		var sourceFile *SourceFile
		if sourceFile, err = output.Generator(ctx); err != nil {
			return
		}

		// Write the output file.
		var outputFile *os.File

		if outputFile, err = os.Create(filePath); err != nil {
			return
		}
		defer outputFile.Close()

		if err = sourceFile.Generate(outputFile); err != nil {
			return
		}
	}

	return
}
