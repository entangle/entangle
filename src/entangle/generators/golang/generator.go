package golang

import (
	"bytes"
	"fmt"
	"entangle/data"
	"entangle/declarations"
	"entangle/generators"
	"text/template"
	"path/filepath"
	"os"
	"code.google.com/p/go.tools/imports"
)

// Go generator.
type generator struct {
	exceptionsTmpl *template.Template
	servicesTmpl *template.Template
	enumsTmpl *template.Template
	structsTmpl *template.Template
	deserializationTmpl *template.Template
	serializationTmpl *template.Template
}

// New generator.
func NewGenerator() (gen generators.Generator, err error) {
	g := &generator {}

	// Define function mapping.
	funcMap := template.FuncMap {
		"documentation": documentationHelper,
		"package": packageHelper,
		"type": typeHelper,
		"nonNilableType": nonNilableTypeHelper,
		"canSkipBeforeField": canSkipBeforeFieldHelper,
		"deserializationCode": deserializationCodeHelper,
		"serializationCode": serializationCodeHelper,
		"structSerializationCode": structSerializationCodeHelper,
		"typeDeserializationMethod": typeDeserializationMethodHelper,
		"fieldIndex": func(fieldDecl *declarations.Field) string { return fmt.Sprintf("%d", fieldDecl.Index - 1) },
	}

	// Load templates.
	for _, info := range []struct {
		Filename string
		Target **template.Template
	} {
		{ "exceptions.go.tmpl", &g.exceptionsTmpl },
		{ "services.go.tmpl", &g.servicesTmpl },
		{ "enums.go.tmpl", &g.enumsTmpl },
		{ "structs.go.tmpl", &g.structsTmpl },
		{ "deserialization.go.tmpl", &g.deserializationTmpl },
		{ "serialization.go.tmpl", &g.serializationTmpl },
	} {
		var src []byte
		path := fmt.Sprintf("templates/generators/golang/%s", info.Filename)

		src, err = data.Asset(path)
		if err != nil {
			return
		}

		var tmpl *template.Template
		tmpl, err = template.New(path).Funcs(funcMap).Parse(string(src))
		if err != nil {
			return
		}

		*info.Target = tmpl
	}

	return g, nil
}

// Generate.
func (g *generator) Generate(interfaceDecl *declarations.Interface, outputPath string) (err error) {
	// Build a serialization/deserialization map for helper functions.
	serDesMap := buildSerDesMap(interfaceDecl)
	fmt.Println(serDesMap)

	// Generate output files.
	for _, output := range []struct {
		Filename string
		Template *template.Template
	} {
		{ "exceptions.go", g.exceptionsTmpl },
		{ "services.go", g.servicesTmpl },
		{ "enums.go", g.enumsTmpl },
		{ "structs.go", g.structsTmpl },
		{ "deserialization.go", g.deserializationTmpl },
		{ "serialization.go", g.serializationTmpl },
	} {
		filePath := filepath.Join(outputPath, output.Filename)

		// Generate the output file.
		buffer := new(bytes.Buffer)

		if err = output.Template.Execute(buffer, struct {
			Interface *declarations.Interface
			SerDesMap map[string]declarations.Type
		} {
			Interface: interfaceDecl,
			SerDesMap: serDesMap,
		}); err != nil {
			return
		}

		// Clean the output file by running it through imports.
		outputData, cleanErr := imports.Process(filePath, buffer.Bytes(), &imports.Options {
			Fragment: false,
			AllErrors: true,
			Comments: true,
			TabIndent: true,
			TabWidth: 4,
		})
		if cleanErr != nil {
			err = fmt.Errorf("error cleaning generated %s: %v", filePath, cleanErr)
		}

		// Write the output file.
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

	return
}
