package python2

import (
	"bytes"
	"entangle/data"
	"entangle/declarations"
	"entangle/generators"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

// Templates.
var templates = []string {

}

// Python 2 generator.
type generator struct {
	options   *Options
	templates []*template.Template
}

// New generator.
func NewGenerator(options *Options) (gen generators.Generator, err error) {
	g := &generator{
		options: options,
		templates: make([]*template.Template, len(templates)),
	}

	// Define function mapping.
	funcMap := template.FuncMap{}

	// Load templates.
	for i, filename := range templates {
		var src []byte
		path := fmt.Sprintf("templates/generators/python2/%s.tmpl", filename)

		src, err = data.Asset(path)
		if err != nil {
			return
		}

		var tmpl *template.Template
		tmpl, err = template.New(path).Funcs(funcMap).Parse(string(src))
		if err != nil {
			return
		}

		g.templates[i] = tmpl
	}

	return g, nil
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
	for i, filename := range templates {
		filePath := filepath.Join(outputPath, filename)

		// Generate the output file.
		buffer := new(bytes.Buffer)

		if err = g.templates[i].Execute(buffer, ctx); err != nil {
			return
		}

		// Write the output file.
		outputData := buffer.Bytes()
		var outputFile *os.File

		if outputFile, err = os.Create(filePath); err != nil {
			return
		}
		defer outputFile.Close()

		var n int
		if n, err = outputFile.Write(outputData); err != nil || n != len(outputData) {
			return
		}
	}

	for _, output := range []struct {
		Filename string
		Generator func(*context) (*SourceFile, error)
	} {
		{ "types.py", generateTypes },
		{ "clients.py", generateClients },
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
