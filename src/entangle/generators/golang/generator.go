package golang

import (
	"bytes"
	"code.google.com/p/go.tools/imports"
	"entangle/data"
	"entangle/declarations"
	"entangle/generators"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// Go generator.
type generator struct {
	options                    *Options
	exceptionsTmpl             *template.Template
	servicesTmpl               *template.Template
	serviceImplementationsTmpl *template.Template
	enumsTmpl                  *template.Template
	structsTmpl                *template.Template
	deserializationTmpl        *template.Template
	serializationTmpl          *template.Template
	serversTmpl                *template.Template
	clientsTmpl                *template.Template
}

// New generator.
func NewGenerator(options *Options) (gen generators.Generator, err error) {
	g := &generator{
		options: options,
	}

	// Define function mapping.
	funcMap := template.FuncMap{
		"documentation":             documentationHelper,
		"type":                      typeHelper,
		"nonNilableType":            nonNilableTypeHelper,
		"canSkipBeforeField":        canSkipBeforeFieldHelper,
		"deserializationCode":       deserializationCodeHelper,
		"serializationCode":         serializationCodeHelper,
		"structSerializationCode":   structSerializationCodeHelper,
		"typeDeserializationMethod": typeDeserializationMethodHelper,
		"typeSerializationCode":     typeSerializationCodeHelper,
		"fieldIndex": func(fieldDecl *declarations.Field) string {
			return fmt.Sprintf("%d", fieldDecl.Index-1)
		},
		"argIndex": func(arg *declarations.FunctionArgument) string {
			return fmt.Sprintf("%d", arg.Index-1)
		},
		"lowerFirst": func(input string) string {
			if input == "" {
				return input
			}
			return strings.ToLower(string(input[0])) + input[1:]
		},
		"argumentOptional": func(arg *declarations.FunctionArgument, minimumDeserializedLength int) bool {
			return arg.Index > uint(minimumDeserializedLength)
		},
	}

	// Load templates.
	for _, info := range []struct {
		Filename string
		Target   **template.Template
	}{
		{"exceptions.go.tmpl", &g.exceptionsTmpl},
		{"services.go.tmpl", &g.servicesTmpl},
		{"service_implementations.go.tmpl", &g.serviceImplementationsTmpl},
		{"enums.go.tmpl", &g.enumsTmpl},
		{"structs.go.tmpl", &g.structsTmpl},
		{"deserialization.go.tmpl", &g.deserializationTmpl},
		{"serialization.go.tmpl", &g.serializationTmpl},
		{"servers.go.tmpl", &g.serversTmpl},
		{"clients.go.tmpl", &g.clientsTmpl},
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

	// Set up the context.
	ctx := &context{
		Interface:   interfaceDecl,
		SerDesMap:   serDesMap,
		PackageName: interfaceDecl.Name,
	}

	// Generate output files.
	for _, output := range []struct {
		Filename string
		Template *template.Template
	}{
		{"exceptions.go", g.exceptionsTmpl},
		{"services.go", g.servicesTmpl},
		{"service_implementations.go", g.serviceImplementationsTmpl},
		{"enums.go", g.enumsTmpl},
		{"structs.go", g.structsTmpl},
		{"deserialization.go", g.deserializationTmpl},
		{"serialization.go", g.serializationTmpl},
		{"servers.go", g.serversTmpl},
		{"clients.go", g.clientsTmpl},
	} {
		filePath := filepath.Join(outputPath, output.Filename)

		// Generate the output file.
		buffer := new(bytes.Buffer)

		if err = output.Template.Execute(buffer, ctx); err != nil {
			return
		}

		// Clean the output file by running it through imports.
		outputData, cleanErr := imports.Process(filePath, buffer.Bytes(), &imports.Options{
			Fragment:  false,
			AllErrors: true,
			Comments:  true,
			TabIndent: true,
			TabWidth:  4,
		})
		if cleanErr != nil {
			err = fmt.Errorf("error cleaning generated %s: %v", filePath, cleanErr)
			return
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
